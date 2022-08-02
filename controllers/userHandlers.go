package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/schattenbrot/simple-login-system/models"
	"github.com/schattenbrot/simple-login-system/utils"
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
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK)
}
