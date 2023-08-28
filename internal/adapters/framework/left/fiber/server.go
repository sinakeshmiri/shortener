package httpfiber

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/keratin/authn-go/authn"

	"github.com/sinakeshmiri/shortener/internal/ports"

	"github.com/gofiber/contrib/otelfiber"
)

// Adapter implements the http Port interface
type Adapter struct {
	authc *authn.Client
	api   ports.APIPort
}

func NewAdapter(authc authn.Config, api ports.APIPort) (*Adapter, error) {
	c, err := authn.NewClient(authc)
	if err != nil {
		return nil, err
	}
	return &Adapter{api: api, authc: c}, nil
}

func (ha Adapter) Run() {

	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Get("/:id", ha.redirect)
	app.Post("/api/v1/url", ha.authnMiddleware, ha.addURL)
	app.Delete("/api/v1/url/:id", ha.authnMiddleware, ha.deleteURL)
	app.Get("/api/v1/url/:id?", ha.authnMiddleware, ha.getMetrics)
	log.Fatal(app.Listen(":3000"))
}
