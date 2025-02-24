package delivery

import (
	"errors"
	"strconv"

	"news-rest-api/internal/dto"
	e "news-rest-api/internal/pkg/errors"
	"news-rest-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) EditNews(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var news dto.News
	if err = c.Bind().JSON(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if news.ID != 0 && news.ID != uint64(id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input ID"})
	}

	if err = h.service.UpdateNews(news, uint64(id)); err != nil {
		switch {
		case errors.Is(err, e.ErrInvalidNewsId):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid news id"})
		case errors.Is(err, e.ErrInvalidCategory):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid category id"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": true})
}

func (h *Handler) GetNewsList(c fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "5"))
	if err != nil || pageSize < 1 {
		pageSize = 5
	}

	news, err := h.service.GetNewsList(page, pageSize)
	if err != nil {
		if errors.Is(err, e.ErrNotFoundNews) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "news not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch news"})
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponseNewsList{
		Success: true,
		News:    news,
	})
}
