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

	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	return page, limit
}

func ParseIdParam(c *fiber.Ctx) (uuid.UUID, error) {
	id := c.Params("id")
	if err := uuid.Validate(id); err != nil {
		return uuid.Nil, apperror.BadRequestError(err, "Invalid ID format")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "Can not parse UUID")
	}

	return parsedID, nil
}
