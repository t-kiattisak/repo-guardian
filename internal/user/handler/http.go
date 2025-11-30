package handler

import (
	"net/http"
	"strconv"

	"repo-guardian/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(f *fiber.App, us domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: us,
	}

	f.Post("/users", handler.Register)
	f.Get("/users/:id", handler.GetUser)
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx := c.Context()
	if err := h.UserUsecase.Register(ctx, &user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx := c.Context()
	user, err := h.UserUsecase.GetUser(ctx, id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}
