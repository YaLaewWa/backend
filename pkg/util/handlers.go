package util

import (
	"socket/pkg/apperror"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PaginationQuery(c *fiber.Ctx) (int, int) {
	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}

	limit := min(c.QueryInt("limit", 10), 50)

	return page, limit
}

func ParseIdParam(id string) (uuid.UUID, error) {
	if err := uuid.Validate(id); err != nil {
		return uuid.Nil, apperror.BadRequestError(err, "Invalid ID format")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "Can not parse UUID")
	}

	return parsedID, nil
}
