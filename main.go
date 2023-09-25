package main

import (
	"context"
	"fmt"
	"github.com/Zavr22/EMTestTask/cache"
	database "github.com/Zavr22/EMTestTask/pkg/db"
	"github.com/Zavr22/EMTestTask/pkg/kafka"
	"github.com/Zavr22/EMTestTask/web/graphql"
	"github.com/Zavr22/EMTestTask/web/rest/handler"
	"github.com/Zavr22/EMTestTask/web/rest/repository"
	"github.com/Zavr22/EMTestTask/web/rest/service"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	e := echo.New()
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalf("Error getting current working directory: %v", err)
	}
	fmt.Printf("Current Working Directory: %s\n", cwd)
	// SWAGGER
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// LOGGERS
	logger := logrus.New()
	logger.Out = os.Stdout

	// REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//if err := godotenv.Load(); err != nil {
	//	logrus.Fatalf("error loading .env: %v", err)
	//}
	// POSTGRES
	db, err := database.NewPostgresDB()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error connection to database rep.NewPostgresDB()": err,
		}).Fatal("DB ERROR CONNECTION")
	}

	defer database.ClosePool(db)

	// init utils, services, repos
	redisClient := cache.NewRedisClient(rdb)
	userRepo := repository.NewUserRepository(db, redisClient)

	userServ := service.NewUserService(userRepo, redisClient)

	profileHandler := handler.NewHandler(userServ)
	profileHandler.InitRoutes(e)
	resolver := graphql.NewResolver(userRepo)
	graphqlHandler := graphql.GraphQLHandler(resolver)
	kafkaConsumer := kafka.NewKafkaConsumer(userServ)

	e.POST("/graphql", graphqlHandler)

	// Graceful shutdown
	logrus.Print("App Started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("App Shutting Down")

	if err := e.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := rdb.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
	go kafkaConsumer.ListenToKafkaTopic()

	select {}
}
