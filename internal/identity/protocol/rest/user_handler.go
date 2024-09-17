package rest

import (
	"amg/internal/api/errors"
	"amg/internal/api/response"
	"amg/internal/identity/user"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	s user.Service
}

func NewUserHandler(s user.Service) *userHandler {
	return &userHandler{
		s: s,
	}
}

func (h *userHandler) CreateUser(ctx *fiber.Ctx) error {
	var cmd user.CreateUserCommand

	err := ctx.BodyParser(&cmd)
	if err != nil {
		return err
	}

	err = cmd.Validate()
	if err != nil {
		return errors.ErrorBadRequest(err)
	}

	err = h.s.CreateUser(ctx.Context(), &cmd)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Created(ctx, fiber.Map{
		"user data": cmd,
	})
}

func (h *userHandler) GetByUserID(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	userID := int64(id)

	result, err := h.s.GetByUserID(ctx.Context(), userID)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"user data": result,
	})
}

func (h *userHandler) UpdateUser(ctx *fiber.Ctx) error {
	var cmd user.UpdateUserCommand

	err := ctx.BodyParser(&cmd)
	if err != nil {
		return err
	}

	err = cmd.Validate()
	if err != nil {
		return errors.ErrorBadRequest(err)
	}

	err = h.s.UpdateUser(ctx.Context(), &cmd)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"user data": cmd,
	})
}

func (h *userHandler) SearchUser(ctx *fiber.Ctx) error {
	var query user.SearchUserQuery

	err := ctx.QueryParser(&query)
	if err != nil {
		return err
	}

	result, err := h.s.SearchUser(ctx.Context(), &query)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, result)
}

func (h *userHandler) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("id")

	userID := int64(id)

	err := h.s.DeleteUser(ctx.Context(), userID)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"message": "user deleted successfully!",
	})
}
