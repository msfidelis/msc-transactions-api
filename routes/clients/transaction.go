package clients

import (
	"fmt"
	"main/dto"
	"main/entities"
	"main/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

var functionName = "NewTransaction"

var request dto.TransactionRequest

var balance int64
var limit int64
var semLimit bool
var err error

func NewTransaction(c *fiber.Ctx) error {

	id := string(c.Request().Header.Peek("id_client"))

	// Parser do Request
	if err := c.BodyParser(&request); err != nil {
		return dto.FiberError(c, fiber.StatusBadRequest, err.Error())
	}

	transacao := &entities.Transaction{
		IDClient:    id,
		Type:        request.Type,
		Amount:      request.Amount,
		Description: request.Description,
		Date:        time.Now().UTC().Format(time.RFC3339Nano),
	}

	balance, limit, semLimit, err = services.Process(*transacao)

	if semLimit {
		return dto.FiberError(c, fiber.StatusUnprocessableEntity, "No limit available")
	}

	if err != nil {
		return dto.FiberError(c, fiber.StatusInternalServerError, err.Error())
	}

	fmt.Printf("[%v] Cliente: %v\n", c.Context().ID(), id)
	fmt.Printf("[%v] Type: %v\n", c.Context().ID(), request.Type)
	fmt.Printf("[%v] Amount: %v\n", c.Context().ID(), request.Amount)
	fmt.Printf("[%v] Description: %v\n", c.Context().ID(), request.Description)
	fmt.Printf("[%v] Balance: %v\n", c.Context().ID(), balance)
	fmt.Printf("[%v] Limit: %v\n", c.Context().ID(), limit)

	response := dto.TransactionResponse{
		Limit:   limit,
		Balance: balance,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
