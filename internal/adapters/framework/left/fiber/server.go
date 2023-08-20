package httpfiber

/*
	NewURL(url,username string) (string, error)
	GetURL(id string) (string, error)
	AddMetrics(id string) error
	DeleteURL(id,username string) error
	GetMetrics(username string) (map[string]int, error)
*/

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/sinakeshmiri/shortner/internal/ports"
)

// Adapter implements the http Port interface
type Adapter struct {
	api ports.APIPort
}

func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (ha Adapter) Run() {
	app := fiber.New()
	app.Get("/r/:id", ha.redirect)
	app.Post("/api/:username", ha.addURL)
	app.Delete("/api/:username/:id", ha.deleteURL)
	app.Get("/api/:username", ha.getMetrics)
	log.Fatal(app.Listen(":3000"))
}
