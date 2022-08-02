package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/schattenbrot/simple-login-system/models"
	"github.com/schattenbrot/simple-login-system/utils"
	"golang.org/x/crypto/bcrypt"
)

func (m *UserRepository) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.DB.GetAllUsers()
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

func (m *UserRepository) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	// Get ID from uri
	id := chi.URLParam(r, "id")

	// Get email from body
	var emailUser struct {
		Email string `json:"email" bson:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&emailUser)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var user models.User
	user.Email = emailUser.Email

	// update user in db
	err = m.DB.UpdateUser(id, user)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK)
}

func (m *UserRepository) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	// Get ID from uri
	id := chi.URLParam(r, "id")

	// Get name from body
	var nameUser struct {
		Name string `json:"name" bson:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&nameUser)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var user models.User
	user.Name = nameUser.Name

	// update user in db
	err = m.DB.UpdateUser(id, user)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK)
}

func (m *UserRepository) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	// get ID from uri
	id := chi.URLParam(r, "id")

	// Get new password from body
	var passwordUser struct {
		Password string `json:"password" bson:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&passwordUser)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordUser.Password), 12)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var user models.User
	user.Password = string(hashedPassword)

	// update user in db
	err = m.DB.UpdateUser(id, user)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK)
}

func (m *UserRepository) RequestResetPassword(w http.ResponseWriter, r *http.Request) {
	// get ID from uri
	id := chi.URLParam(r, "id")

	// create reset password token
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	token := hex.EncodeToString(b)

	exp := time.Now().Add(10 * time.Minute)

	// set resetpassword and resetpassword expire in db
	var user models.User
	user.ResetPasswordToken = token
	user.ResetPasswordTokenExpire = exp

	err := m.DB.UpdatePasswordResetToken(id, user)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	type respData struct {
		Token string    `json:"token"`
		Exp   time.Time `json:"exp"`
	}

	resp := respData{
		Token: token,
		Exp:   exp,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (m *UserRepository) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// get ID from uri
	id := chi.URLParam(r, "id")

	// get password from body
	var passwordUser struct {
		Password string `bson:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&passwordUser)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordUser.Password), 12)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var user models.User
	user.Password = string(hashedPassword)
	user.ResetPasswordToken = ""

	// update user in db
	err = m.DB.UpdateUser(id, user)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK)
}
