package clients

import (
	"main/dto"
	"main/pkg/database"
	"main/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Statement(c *fiber.Ctx) error {

	id := string(c.Request().Header.Peek("id_client"))

	db := database.GetDB()

	client, err := services.FindClient(c.Context(), db, id)
	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, "Client not found")
	}

	transactions, err := services.Statement(id)
	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, "Error to recover transactions")
	}

	response := dto.StatementResponse{
		LastTransactions: transactions,
	}

	response.Balance.Total = client.Balance
	response.Balance.Limit = client.Limit
	response.Balance.DateStatement = time.Now().UTC().Format(time.RFC3339Nano)

	return c.JSON(response)

}
