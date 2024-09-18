package server

import (
	"amg/internal/api/response"
	"amg/internal/db"
	"amg/internal/identity/protocol/rest"
	"amg/internal/identity/user/userimpl"
	"amg/internal/middleware"
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	requireCreateUser = middleware.RequirePermission("create")
	requireReadUser   = middleware.RequirePermission("read")
	requireUpdateUser = middleware.RequirePermission("update")
	requireDeleteUser = middleware.RequirePermission("delete")

	reqOnlyByAdmin      = middleware.RequireRole("admin")
	reqBothUserAndAdmin = middleware.RequireRole("user", "admin")
)

func healthCheck(db db.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var result int64
		err := db.Get(ctx.Context(), &result, "SELECT 1")
		if err != nil {
			return errors.New("database unavailable")
		}

		return response.Ok(ctx, fiber.Map{
			"database": "available",
		})
	}
}

func (s *Server) SetupRoutes() {
	api := s.app.Group("/api")
	api.Get("/health", healthCheck(s.db))

	// User Routes

	user := userimpl.NewService(s.db, s.cfg)
	userHttp := rest.NewUserHandler(user)

	api.Post("/users/register", userHttp.RegisterDefaultUser)
	api.Post("/users/login", userHttp.LoginUser)

	api.Use(middleware.JWTProtected(s.jwtSecret, user))
	api.Post("/users", reqOnlyByAdmin, requireCreateUser, userHttp.CreateUser)
	api.Get("/users", reqBothUserAndAdmin, requireReadUser, userHttp.SearchUser)
	api.Get("/users/:id", reqBothUserAndAdmin, requireReadUser, userHttp.GetByUserID)
	api.Put("/users/:id", reqOnlyByAdmin, requireUpdateUser, userHttp.UpdateUser)
	api.Delete("/users/:id", reqOnlyByAdmin, requireDeleteUser, userHttp.DeleteUser)

	// Logout
	api.Post("/users/logout", userHttp.LogoutUser)

}
