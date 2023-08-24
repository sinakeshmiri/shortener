package httpfiber

import (
	"log"

	"github.com/keratin/authn-go/authn"
	"github.com/gofiber/fiber/v2"

	"github.com/sinakeshmiri/shortner/internal/ports"


	"github.com/gofiber/contrib/otelfiber"


)

// Adapter implements the http Port interface
type Adapter struct {
	authc *authn.Client
	api ports.APIPort

}

func NewAdapter(authc authn.Config, api ports.APIPort) (*Adapter,error) {
	c,err := authn.NewClient(authc)
	if err != nil {
		return nil,err
	}
	return &Adapter{api: api,authc: c},nil
}

func (ha Adapter) Run() {
	
	app := fiber.New()
	app.Use(otelfiber.Middleware())
	app.Get("/:id",ha.redirect)
	app.Post("/api/:username",ha.authnMiddleware, ha.addURL)
	app.Delete("/api/:username/:id",ha.authnMiddleware, ha.deleteURL)
	app.Get("/api/:username",ha.authnMiddleware, ha.getMetrics)
	log.Fatal(app.Listen(":3000"))
}
