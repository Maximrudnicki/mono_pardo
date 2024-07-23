package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"mono_pardo/cmd/config"
	"mono_pardo/cmd/controller"
	"mono_pardo/cmd/model"
	"mono_pardo/cmd/repository"
	"mono_pardo/cmd/router"
	"mono_pardo/cmd/service"

	"github.com/go-playground/validator"
	"github.com/rs/cors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	validate := validator.New()

	//Database
	db := config.ConnectionDB(&loadConfig)

	db_table_err := db.Table("words").AutoMigrate(&model.Word{})
	if db_table_err != nil {
		log.Fatalf("Database table error: %v\n", db_table_err)
	}

	db_table_err = db.Table("users").AutoMigrate(&model.User{})
	if db_table_err != nil {
		log.Fatalf("Database table error: %v\n", db_table_err)
	}

	mc := fmt.Sprintf(loadConfig.MONGODB_STRING)

	clientOptions := options.Client().ApplyURI(mc)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
	}

	groupCollection := client.Database(loadConfig.MONGODB_DB).Collection("groups")
	statisticsCollection := client.Database(loadConfig.MONGODB_DB).Collection("statistics")

	//Init Repositories
	groupRepository := repository.NewGroupRepositoryImpl(groupCollection)
	statisticsRepository := repository.NewStatisticsRepositoryImpl(statisticsCollection)

	userRepository := repository.NewUsersRepositoryImpl(db)
	wordRepository := repository.NewWordRepositoryImpl(db)

	//Init Services
	authenticationService := service.NewAuthenticationServiceImpl(loadConfig, validate, userRepository)
	vocabService := service.NewVocabServiceImpl(authenticationService, validate, wordRepository)
	groupService := service.NewGroupServiceImpl(
		authenticationService, vocabService, validate, groupRepository, statisticsRepository)

	//Init controllers
	authenticationController := controller.NewAuthenticationController(authenticationService)
	vocabController := controller.NewVocabController(vocabService)
	groupController := controller.NewGroupController(groupService, vocabService)

	r := router.NewRouter(authenticationController, vocabController, groupController)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{loadConfig.ALLOWED_ORIGINS},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", loadConfig.PORT),
		Handler: handler,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
