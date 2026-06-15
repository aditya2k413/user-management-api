package handler

import (
	customErrors "UserAgeAPI/internal/errors"
	"UserAgeAPI/internal/models"
	"UserAgeAPI/internal/service"
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service   *service.UserService
	validator *validator.Validate
}

func NewUserHandler(service *service.UserService) *UserHandler {
	validate := validator.New()
	return &UserHandler{
		service:   service,
		validator: validate,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "validation failed: invalid input format",
		})
	}

	user, err := h.service.CreateUser(
		c.UserContext(),
		req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}
	user, err := h.service.GetUser(
		c.UserContext(),
		int32(id),
	)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(user)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.ListUsers(c.UserContext())
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(users)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid user id",
			},
		)
	}

	var req models.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid request body",
			},
		)
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "validation failed: invalid input format",
		})
	}

	user, err := h.service.UpdateUser(
		c.UserContext(),
		int32(id),
		req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(user)
}
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid user id",
			},
		)
	}

	err = h.service.DeleteUser(
		c.UserContext(),
		int32(id),
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, customErrors.ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"error": "user not found",
			},
		)

	case errors.Is(err, customErrors.ErrInvalidDate):
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid date format, use YYYY-MM-DD",
			},
		)

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "internal server error",
			},
		)
	}
}
