package handler

import (
	"EMTestTask/pkg/model"
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// User interface consists of user service methods
type User interface {
	CreateUser(ctx context.Context, user *model.FIO) error
	GetAllUsers(ctx context.Context, page int) ([]*model.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (model.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input model.EnrichedFIO) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type Handler struct {
	userS User
}

func NewHandler(userS User) *Handler {
	return &Handler{userS: userS}
}

// InitRoutes is used to init routes for web service
func (h *Handler) InitRoutes(router *echo.Echo) *echo.Echo {

	api := router.Group("api")
	api.POST("/users", h.CreateUser)
	api.GET("/users", h.GetUsers)
	api.GET("/users/:id", h.GetUserByID)
	api.PUT("/users/:id", h.UpdateUser)
	api.DELETE("/users/:id", h.DeleteUser)

	router.Logger.Fatal(router.Start(":9000"))
	return router
}