package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/schattenbrot/simple-login-system/config"
	"github.com/schattenbrot/simple-login-system/controllers"
	"github.com/schattenbrot/simple-login-system/database"
	"github.com/schattenbrot/simple-login-system/middlewares"
	"github.com/schattenbrot/simple-login-system/routes"
)

func main() {
	var cfg config.Config
	var stringJWT string
	flag.IntVar(&cfg.Port, "port", 4000, "Server port to listen on.")
	flag.StringVar(&cfg.Env, "env", "dev", "Application environment (dev | prod)")
	flag.StringVar(&cfg.DB.DSN, "db", "mongodb://localhost:27017", "Mongo DB connection string")
	flag.StringVar(&stringJWT, "jwt", "3710df151cbf8006c799e99037b26deb998a2df2f067f0f750172eebcdc8439a", "jwt secret")
	flag.Parse()
	cfg.JWT = []byte(stringJWT)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &config.Application{
		Version:   "1.0.0",
		Config:    cfg,
		Logger:    logger,
		StartTime: time.Now(),
		Validator: validator.New(),
	}

	db := database.Open(app)
	controllerRepo := controllers.NewRepo(app, db)
	controllers.NewHandlers(controllerRepo)
	middlewaresRepo := middlewares.NewRepo(app, db)
	middlewares.NewMiddlewares(middlewaresRepo)

	app.Logger.Println("Starting server on port", cfg.Port)

	serve := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routes.ChiRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := serve.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
