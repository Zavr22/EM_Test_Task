package handler

import (
	"context"
	"github.com/Zavr22/EMTestTask/pkg/graph"
	"github.com/Zavr22/EMTestTask/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// User interface consists of user service methods
type User interface {
	CreateUser(ctx context.Context, user *models.FIO) (uuid.UUID, error)
	GetAllUsers(ctx context.Context, page int) ([]*models.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.EnrichedFIO) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
	EnrichAndSaveToDB(ctx context.Context, name, surname, patronymic string) (uuid.UUID, error)
}

type Handler struct {
	userS User
}

func NewHandler(userS User) *Handler {
	return &Handler{userS: userS}
}

// InitRoutes is used to init routes for web service
func (h *Handler) InitRoutes(router *echo.Echo, graphqlHandler echo.HandlerFunc) *echo.Echo {

	api := router.Group("/api")
	api.POST("/users", h.CreateUser)
	api.GET("/users", h.GetUsers)
	api.GET("/users/:id", h.GetUserByID)
	api.PUT("/users/:id", h.UpdateUser)
	api.DELETE("/users/:id", h.DeleteUser)
	router.POST("/query", graphqlHandler)
	router.GET("/", graph.PlaygroundHandler())

	router.Logger.Fatal(router.Start(":9000"))
	return router
}
