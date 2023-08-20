package main

import (
	"log"
	"os"

	// application
	"github.com/keratin/authn-go/authn"
	"github.com/sinakeshmiri/shortner/internal/application/api"
	"github.com/sinakeshmiri/shortner/internal/application/core/shortner"

	// adapters
	hFiber "github.com/sinakeshmiri/shortner/internal/adapters/framework/left/fiber"
	"github.com/sinakeshmiri/shortner/internal/adapters/framework/right/db"
)

func main() {
	var err error

	NODE_ID := os.Getenv("NODE_ID")
	MONGO_URI := os.Getenv("MONGO_URI")
	AUTHN_URL := os.Getenv("AUTHN_URL")
	AUTHN_PASSWORD := os.Getenv("AUTHN_PASSWORD")
	AUTHN_USERNAME := os.Getenv("AUTHN_USERNAME")
	AUTHN_ISSUER := os.Getenv("AUTHN_ISSUER")
	AUTHN_AUDIENCE := os.Getenv("AUTHN_AUDIENCE")

	dbAdapter, err := db.NewAdapter(MONGO_URI)
	if err != nil {
		log.Fatalf("failed to initiate dbase connection: %v", err)
	}
	defer dbAdapter.CloseDbConnection()

	// core
	core, err := shortner.New(NODE_ID)
	if err != nil {
		log.Fatalf("failed to initiate core: %v", err)
	}

	applicationAPI := api.NewApplication(dbAdapter, core)


	ac:=authn.Config{
		Issuer: AUTHN_ISSUER,
		Audience: AUTHN_AUDIENCE,
		Username: AUTHN_USERNAME,
		Password: AUTHN_PASSWORD,
		PrivateBaseURL: AUTHN_URL,
	}

	hFiberAdapter,err := hFiber.NewAdapter(applicationAPI,ac)
	if err != nil {
		log.Fatalf("failed to initiate http adapter: %v", err)
	}
	hFiberAdapter.Run()
}
