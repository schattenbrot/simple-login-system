package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/schattenbrot/simple-login-system/models"
	"github.com/schattenbrot/simple-login-system/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	var authUser AuthUser

	err := json.NewDecoder(r.Body).Decode(&authUser)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), 12)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:     authUser.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// db access
	userId, err := m.DB.CreateUser(user)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	type jsonResp struct {
		ID string `json:"id"`
	}

	response := jsonResp{
		ID: *userId,
	}

	utils.WriteJSON(w, http.StatusCreated, response)
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	// get login data
	var authUser AuthUser

	err := json.NewDecoder(r.Body).Decode(&authUser)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// check if user exists in db
	user, err := m.DB.GetUserByUsername(authUser.Email)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	fmt.Println("user:")
	fmt.Println(user)

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password))
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// create JWT token with bearer and expire date
	expirationDate := time.Now().Add(7 * 24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.ID,
		ExpiresAt: expirationDate,
	})
	tokenString, err := token.SignedString(m.App.Config.JWT)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// send response
	type jsonResp struct {
		ID             string `json:"id"`
		Token          string `json:"token"`
		ExpirationDate int64  `json:"exp"`
	}
	utils.WriteJSON(w, http.StatusOK, jsonResp{
		ID:             user.ID,
		Token:          tokenString,
		ExpirationDate: expirationDate,
	})
}
