package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/schattenbrot/simple-login-system/utils"
	"golang.org/x/crypto/bcrypt"
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

func (m *Repository) ValidateUpdatePasswordPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var passwordUser struct {
			Token    string `json:"token" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		utils.MiddlewareBodyDecoder(r, &passwordUser)

		// check token validity
		id := chi.URLParam(r, "id")
		user, err := m.DB.GetUserById(id)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		isExpired := time.Now().Unix() > user.ResetPasswordTokenExpire.Unix()
		if user.ResetPasswordToken != passwordUser.Token && isExpired {
			utils.ErrorJSON(w, errors.New("token invalid"), http.StatusUnauthorized)
			return
		}

		// validate payload
		err = m.App.Validator.Struct(passwordUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		// check if password is the same
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordUser.Password))
		if err == nil {
			utils.ErrorJSON(w, errors.New("new password should not match the old password"), http.StatusInternalServerError)
			return
		} else if err != bcrypt.ErrMismatchedHashAndPassword {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return

		}

		next.ServeHTTP(w, r)
	})
}
