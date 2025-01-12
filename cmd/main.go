package main

import (
	"fmt"
	"log"
	"net/http"

	"mono_pardo/internal/api"
	"mono_pardo/internal/api/controller"
	usersDomain "mono_pardo/internal/domain/users"
	wordsDomain "mono_pardo/internal/domain/words"
	usersInfra "mono_pardo/internal/infrastructure/users"
	wordsInfra "mono_pardo/internal/infrastructure/words"
	"mono_pardo/pkg/config"

	"github.com/go-playground/validator"
	"github.com/rs/cors"
)

func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	validate := validator.New()

	//Database
	db := config.ConnectionDB(&loadConfig)

	if err = db.Table("words").AutoMigrate(&wordsDomain.Word{}); err != nil {
		log.Fatalf("Database table error: %v\n", err)
	}

	if err = db.Table("users").AutoMigrate(&usersDomain.User{}); err != nil {
		log.Fatalf("Database table error: %v\n", err)
	}

	//Init Repositories
	userRepository := usersInfra.NewPostgresRepositoryImpl(db)
	wordRepository := wordsInfra.NewPostgresRepositoryImpl(db)

	//Init Services
	authenticationService := usersDomain.NewServiceImpl(loadConfig, validate, userRepository)
	vocabService := wordsDomain.NewServiceImpl(validate, wordRepository)

	//Init controllers
	authenticationController := controller.NewAuthenticationController(authenticationService)
	vocabController := controller.NewVocabController(vocabService)

	router := api.NewRouter(authenticationController, vocabController)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{loadConfig.ALLOWED_ORIGINS},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", loadConfig.PORT),
		Handler: handler,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
