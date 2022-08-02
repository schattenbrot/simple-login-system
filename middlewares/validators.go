package middlewares

import (
	"net/http"

	"github.com/schattenbrot/simple-login-system/utils"
)

func (m *Repository) HasEmail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var emailUser struct {
			Email string `json:"email" validate:"required,email"`
		}

		utils.MiddlewareBodyDecoder(r, &emailUser)

		err := m.App.Validator.Struct(emailUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Repository) HasUsername(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nameUser struct {
			Name string `json:"name" validate:"required"`
		}

		utils.MiddlewareBodyDecoder(r, &nameUser)

		err := m.App.Validator.Struct(nameUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
