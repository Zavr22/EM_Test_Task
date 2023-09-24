package service

import (
	"context"
	"github.com/Zavr22/EMTestTask/pkg/model"
	"github.com/google/uuid"
)

type User interface {
	CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error)
	GetAllUsers(ctx context.Context, offset int) ([]*model.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (model.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input model.EnrichedFIO) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
	SaveUser(ctx context.Context, user *model.FIO) (uuid.UUID, error)
}

type UserService struct {
	userRepo User
}

func NewUserService(userRepo User) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.FIO) (uuid.UUID, error) {
	return s.userRepo.SaveUser(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context, page int) ([]*model.User, error) {
	offset := (page - 1) * 30
	return s.userRepo.GetAllUsers(ctx, offset)
}

func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (model.User, error) {
	return s.userRepo.GetUser(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input model.EnrichedFIO) error {
	return s.userRepo.UpdateProfile(ctx, userID, input)
}

func (s *UserService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}
