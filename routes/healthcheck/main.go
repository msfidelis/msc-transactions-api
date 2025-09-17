package healthcheck

import fiber "github.com/gofiber/fiber/v2"

func Probe(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
