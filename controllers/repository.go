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

var Repo *Repository

func NewRepo(app *config.Application, db *mongo.Database) *Repository {
	return &Repository{
		App: app,
		DB:  database.NewRepo(app, db),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}
