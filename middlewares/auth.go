package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/schattenbrot/simple-login-system/utils"
)

type AuthUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type contextKey int

const authenticatedUserKey contextKey = 0

func (m *Repository) ValidateAuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authUser AuthUser

		utils.MiddlewareBodyDecoder(r, &authUser)

		err := m.App.Validator.Struct(authUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
		err = utils.PasswordValidator(authUser.Password)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Repository) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			utils.ErrorJSON(w, errors.New("authorization header not found"), http.StatusUnauthorized)
			return
		}
		token := strings.Split(authHeader, " ")[1] // Bearer token
		claims := jwt.StandardClaims{}
		_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
			return m.App.Config.JWT, nil
		})
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		// check if bearer exists
		issuer := claims.Issuer
		_, err = m.DB.GetUserById(issuer)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		// set user in context
		ctxWithUser := context.WithValue(r.Context(), authenticatedUserKey, issuer)

		// check if the token is expired
		expirationTime := claims.ExpiresAt
		currTime := time.Now().Unix()

		if currTime > expirationTime {
			utils.ErrorJSON(w, errors.New("token expired"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

func (m *Repository) IsUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issuerId := r.Context().Value(authenticatedUserKey).(string)
		targetId := chi.URLParam(r, "id")

		if issuerId != targetId {
			utils.ErrorJSON(w, errors.New("wrong user"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
