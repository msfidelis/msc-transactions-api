package balance

import (
	"main/dto"
	"main/services"

	"github.com/gofiber/fiber/v2"
)

func GetStatement(c *fiber.Ctx) error {
	id := string(c.Request().Header.Peek("id_client"))
	balance, err := services.GetBalance(id)
	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, "Error fetching balance")
	}
	return c.JSON(fiber.Map{"balance": balance})
}
