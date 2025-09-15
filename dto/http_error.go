package dto

import "github.com/gofiber/fiber/v2"

type HttpError struct {
	Message string `json:"message"`
}

func FiberError(c *fiber.Ctx, statusCode int, mensagem string) error {
	return c.Status(statusCode).JSON(HttpError{Message: mensagem})
}
