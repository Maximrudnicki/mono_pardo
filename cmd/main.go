package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mono_pardo/internal/api"
	"mono_pardo/internal/api/controller"
	setsDomain "mono_pardo/internal/domain/sets"
	usersDomain "mono_pardo/internal/domain/users"
	wordsDomain "mono_pardo/internal/domain/words"
	setsInfra "mono_pardo/internal/infrastructure/sets"
	usersInfra "mono_pardo/internal/infrastructure/users"
	wordsInfra "mono_pardo/internal/infrastructure/words"
	"mono_pardo/pkg/config"
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

	mc := fmt.Sprintf(loadConfig.MONGODB_STRING)

	clientOptions := options.Client().ApplyURI(mc)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
	}

	setsCollection := client.Database(loadConfig.MONGODB_DB).Collection("sets")

	//Init Repositories
	userRepository := usersInfra.NewPostgresRepositoryImpl(db)
	wordRepository := wordsInfra.NewPostgresRepositoryImpl(db)
	setsRepository := setsInfra.NewMongoRepositoryImpl(setsCollection)

	//Init Services
	authenticationService := usersDomain.NewServiceImpl(loadConfig, validate, userRepository)
	vocabService := wordsDomain.NewServiceImpl(validate, wordRepository)
	setsService := setsDomain.NewServiceImpl(validate, setsRepository)

	//Init controllers
	authenticationController := controller.NewAuthenticationController(authenticationService)
	vocabController := controller.NewVocabController(vocabService)
	setsController := controller.NewSetsController(setsService)

	router := api.NewRouter(authenticationController, vocabController, setsController)

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
