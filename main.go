package main

import (
	"context"
	"fmt"
	handlerG "github.com/99designs/gqlgen/graphql/handler"
	"github.com/Zavr22/EMTestTask/pkg/cache"
	database "github.com/Zavr22/EMTestTask/pkg/db"
	"github.com/Zavr22/EMTestTask/pkg/graph"
	kafkaServ "github.com/Zavr22/EMTestTask/pkg/kafka"
	"github.com/Zavr22/EMTestTask/pkg/rest/handler"
	"github.com/Zavr22/EMTestTask/pkg/rest/repository"
	"github.com/Zavr22/EMTestTask/pkg/rest/service"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
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
		Addr:     fmt.Sprintf("redis:%s", os.Getenv("REDIS_PORT")),
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

	resolver := graph.NewResolver(userServ)

	graphqlH := handlerG.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)

	network := os.Getenv("KAFKA_NETWORK")
	address := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	topic := os.Getenv("KAFKA_TOPIC")
	partition := 0

	conn, errKafka := kafka.DialLeader(context.Background(), network, address, topic, partition)
	if errKafka != nil {
		logrus.WithFields(logrus.Fields{
			"Error connection to kafka": errKafka,
		}).Fatal("kafka error connection")
	}
	defer conn.Close()

	kafkaConsumer := kafkaServ.NewKafkaConsumer(userServ, conn)
	profileHandler := handler.NewHandler(userServ)

	go kafkaConsumer.ListenToKafkaTopic()
	kafkaConsumer.ProduceMessage()

	profileHandler.InitRoutes(e, graph.NewGraphqlHandler(graphqlH))
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
	if err := conn.Close(); err != nil {
		logrus.Errorf("error occured on kafka connection close: %s", err.Error())
	}
	select {}
}
