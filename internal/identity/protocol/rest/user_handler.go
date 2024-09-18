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

func (h *userHandler) RegisterDefaultUser(ctx *fiber.Ctx) error {
	var cmd user.RegisterUserCommand

	err := ctx.BodyParser(&cmd)
	if err != nil {
		return err
	}

	err = cmd.Validate()
	if err != nil {
		return errors.ErrorBadRequest(err)
	}

	err = h.s.RegisterDefaultUser(ctx.Context(), &cmd)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"user data": cmd,
	})
}

func (h *userHandler) LoginUser(ctx *fiber.Ctx) error {
	var cmd user.LoginUserCommand

	err := ctx.BodyParser(&cmd)
	if err != nil {
		return err
	}

	err = cmd.Validate()
	if err != nil {
		return errors.ErrorBadRequest(err)
	}

	result, err := h.s.GetUserByEmail(ctx.Context(), &cmd)
	if err != nil {
		if err == user.ErrUserNotFound || err == user.ErrInvalidPassword {
			return errors.ErrorUnauthorized(err, "Invalid email or password")
		}
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, result)
}

func (h *userHandler) LogoutUser(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	if token == "" {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "No token provided",
		}
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := h.s.InvalidateToken(ctx.Context(), token)
	if err != nil {
		return errors.ErrorInternalServerError(err)
	}

	return response.Ok(ctx, fiber.Map{
		"message": "user logged out successfully!",
	})
}
