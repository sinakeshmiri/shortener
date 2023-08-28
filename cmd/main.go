package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel"

	"net/http"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	// application
	"github.com/keratin/authn-go/authn"
	"github.com/sinakeshmiri/shortener/internal/application/api"
	"github.com/sinakeshmiri/shortener/internal/application/core/shortener"

	// adapters
	hFiber "github.com/sinakeshmiri/shortener/internal/adapters/framework/left/fiber"
	"github.com/sinakeshmiri/shortener/internal/adapters/framework/right/db"
)

func main() {
	var err error
	OTEL_EXPORTER_JAEGER_ENDPOINT := os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT")
	NODE_ID := os.Getenv("NODE_ID")
	MONGO_URI := os.Getenv("MONGO_URI")
	AUTHN_URL := os.Getenv("AUTHN_URL")
	AUTHN_PASSWORD := os.Getenv("AUTHN_PASSWORD")
	AUTHN_USERNAME := os.Getenv("AUTHN_USERNAME")
	AUTHN_ISSUER := os.Getenv("AUTHN_ISSUER")
	AUTHN_AUDIENCE := os.Getenv("AUTHN_AUDIENCE")
	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(":2112", nil)

	tp := initTracer(OTEL_EXPORTER_JAEGER_ENDPOINT)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	dbAdapter, err := db.NewAdapter(MONGO_URI)
	if err != nil {
		log.Fatalf("failed to initiate dbase connection: %v", err)
	}
	defer dbAdapter.CloseDbConnection()

	// core
	urlschan := make(chan string, 1000)
	core, err := shortener.New(NODE_ID, urlschan)
	if err != nil {
		log.Fatalf("failed to initiate core: %v", err)
	}

	applicationAPI := api.NewApplication(dbAdapter, core, urlschan)

	ac := authn.Config{
		Issuer:         AUTHN_ISSUER,
		Audience:       AUTHN_AUDIENCE,
		Username:       AUTHN_USERNAME,
		Password:       AUTHN_PASSWORD,
		PrivateBaseURL: AUTHN_URL,
	}

	hFiberAdapter, err := hFiber.NewAdapter(ac, applicationAPI)
	if err != nil {
		log.Fatalf("failed to initiate http adapter: %v", err)
	}
	hFiberAdapter.Run()
}
func initTracer(OTEL_EXPORTER_JAEGER_ENDPOINT string) *sdktrace.TracerProvider {
	//exporter, err := stdout.New(stdout.WithPrettyPrint())
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(OTEL_EXPORTER_JAEGER_ENDPOINT)))
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("shorter"),
			)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
