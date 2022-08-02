package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/schattenbrot/simple-login-system/controllers"
	"github.com/schattenbrot/simple-login-system/middlewares"
)

func ChiRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		MaxAge:         300,
	}))

	r.Get("/", controllers.Repo.StatusHandler)
	r.Route("/users", userRouter)

	return r
}

func userRouter(r chi.Router) {
	r.With(middlewares.Repo.ValidateAuthUser).Post("/register", controllers.User.Register)
	r.With(middlewares.Repo.ValidateAuthUser).Post("/login", controllers.User.Login)
	r.With().Get("/{id}/resetPassword", controllers.User.RequestResetPassword)
	r.With(middlewares.Repo.ValidateUpdatePasswordPayload).Patch("/{id}/updatePassword", controllers.User.UpdatePassword)

	r.With(middlewares.Repo.IsAuth).Get("/", controllers.User.GetUsers)
	r.With(middlewares.Repo.IsAuth, middlewares.Repo.IsUser, middlewares.Repo.HasEmail).Patch("/{id}/email", controllers.User.UpdateUserEmail)
	r.With(middlewares.Repo.IsAuth, middlewares.Repo.IsUser, middlewares.Repo.HasUsername).Patch("/{id}/name", controllers.User.UpdateUserName)
}
