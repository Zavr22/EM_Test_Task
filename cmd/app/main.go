package main

import (
	"EMTestTask/cache"
	"EMTestTask/pkg/db"
	"EMTestTask/pkg/kafka"
	"EMTestTask/web/graphql"
	"EMTestTask/web/rest/handler"
	"EMTestTask/web/rest/repository"
	"EMTestTask/web/rest/service"
	"context"
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

	// SWAGGER
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// LOGGERS
	logger := logrus.New()
	logger.Out = os.Stdout

	// REDIS
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

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

	userServ := service.NewUserService(userRepo)

	profileHandler := handler.NewHandler(userServ)
	profileHandler.InitRoutes(e)
	resolver := graphql.NewResolver(userRepo)
	graphqlHandler := graphql.GraphQLHandler(resolver)

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
	go kafka.ListenToKafkaTopic(userRepo, redisClient)

	select {}
}
