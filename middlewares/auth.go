package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/schattenbrot/simple-login-system/utils"
)

type AuthUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (m *Repository) ValidateAuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var authUser AuthUser

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&authUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		err = m.App.Validator.Struct(authUser)
		if err != nil {
			utils.ErrorJSON(w, err)
			m.App.Logger.Println("validateauthuser 2")
			return
		}
		err = utils.PasswordValidator(authUser.Password)
		if err != nil {
			utils.ErrorJSON(w, err)
			m.App.Logger.Println("validateauthuser 3")
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
		// check if the token is expired
		expirationTime := claims.ExpiresAt
		currTime := time.Now().Unix()

		if currTime > expirationTime {
			utils.ErrorJSON(w, errors.New("token expired"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
