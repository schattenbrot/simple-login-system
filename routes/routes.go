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
	r.With(middlewares.Repo.ValidateAuthUser, middlewares.Repo.IsAuth).Post("/register", controllers.Repo.Register)
	r.With(middlewares.Repo.ValidateAuthUser).Post("/login", controllers.Repo.Login)
}
