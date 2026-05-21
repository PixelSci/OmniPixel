package response

import "github.com/gofiber/fiber/v3"

func Write(c fiber.Ctx, apiErr *APIError) error {
	return c.Status(HTTPStatus(apiErr.Code)).JSON(apiErr)
}

func DomainError(c fiber.Ctx, err error) error {
	if apiErr, ok := domainMappings[err]; ok {
		return Write(c, apiErr)
	}
	return Write(c, ErrInternalServer)
}
