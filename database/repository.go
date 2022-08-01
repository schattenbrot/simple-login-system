package database

import (
	"context"
	"time"

	"github.com/schattenbrot/simple-login-system/config"
	"github.com/schattenbrot/simple-login-system/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseRepository interface {
	CreateUser(user models.User) (*string, error)
	GetUserById(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
}

type dbRepo struct {
	App *config.Application
	DB  *mongo.Database
}

func NewRepo(app *config.Application, conn *mongo.Database) DatabaseRepository {
	return &dbRepo{
		App: app,
		DB:  conn,
	}
}

// openDB creates a new database connection and returns the Database
func Open(app *config.Application) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(app.Config.DB.DSN))
	if err != nil {
		app.Logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		app.Logger.Fatal(err)
	}
	db := client.Database("simpleLoginSystem")

	return db
}
