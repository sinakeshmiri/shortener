package main

import (
	"log"
	"os"

	// application
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
	NODE_ID = "node1"
	MONGO_URI = "mongodb://localhost:27017"
	///
	dbAdapter, err := db.NewAdapter(MONGO_URI )
	if err != nil {
		log.Fatalf("failed to initiate dbase connection: %v", err)
	}
	defer dbAdapter.CloseDbConnection()

	
	// core
	core,err:= shortner.New(NODE_ID)
	if err != nil {
		log.Fatalf("failed to initiate core: %v", err)
	}
	// NOTE: The application's right side port for driven
	// adapters, in this case, a db adapter.
	// Therefore the type for the dbAdapter parameter
	// that is to be injected into the NewApplication will
	// be of type DbPort
	applicationAPI := api.NewApplication(dbAdapter, core)

	// NOTE: We use dependency injection to give the grpc
	// adapter access to the application, therefore
	// the location of the port is inverted. That is
	// the grpc adapter accesses the hexagon's driving port at the
	// application boundary via dependency injection,
	// therefore the type for the applicaitonAPI parameter
	// that is to be injected into the gRPC adapter will
	// be of type APIPort which is our hexagons left side
	// port for driving adapters
	hFiberAdapter := hFiber.NewAdapter(applicationAPI)
	hFiberAdapter.Run()
}
