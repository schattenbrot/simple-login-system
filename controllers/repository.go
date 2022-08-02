package controllers

import (
	"github.com/schattenbrot/simple-login-system/config"
	"github.com/schattenbrot/simple-login-system/database"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	App *config.Application
	DB  database.DatabaseRepository
}

type UserRepository Repository

var Repo *Repository
var User *UserRepository

func NewRepo(app *config.Application, db *mongo.Database) *Repository {
	return &Repository{
		App: app,
		DB:  database.NewRepo(app, db),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
	user := UserRepository(*r)
	User = &user
}
