package httpfiber

import (
	"github.com/gofiber/fiber/v2"
)

type URL struct {
	URL string `json:"url"`
}

func (ha Adapter) redirect(c *fiber.Ctx) error {
	url, err := ha.api.GetURL(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	ha.api.AddMetrics(c.Params("id"))
	return c.Redirect(url, fiber.StatusMovedPermanently)
}

func (ha Adapter) addURL(c *fiber.Ctx) error {

	url := new(URL)
	if err := c.BodyParser(url); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	urlID, err := ha.api.NewURL(url.URL, c.Params("username"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(fiber.Map{"urlID": urlID})
}

func (ha Adapter) deleteURL(c *fiber.Ctx) error {
	err := ha.api.DeleteURL(c.Params("id"), c.Params("username"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return c.JSON(fiber.Map{"status": "success"})
}

func (ha Adapter) getMetrics(c *fiber.Ctx) error {
	id:= c.Params("id")
	metrics, err := ha.api.GetMetrics(c.Params("username"),id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}
	return c.JSON(fiber.Map{"metrics": metrics})
}
