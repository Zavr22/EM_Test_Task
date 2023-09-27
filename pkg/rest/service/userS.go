package service

import (
	"context"
	"fmt"
	"github.com/Zavr22/EMTestTask/pkg/api"
	"github.com/Zavr22/EMTestTask/pkg/cache"
	"github.com/Zavr22/EMTestTask/pkg/models"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -source=userS.go -destination=mocks/user_mock.go

type User interface {
	CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetAllUsers(ctx context.Context, offset int) ([]*models.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, input models.EnrichedFIO) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
}

type UserService struct {
	userRepo User
	redis    *cache.RedisClient
}

func NewUserService(userRepo User, redis *cache.RedisClient) *UserService {
	return &UserService{userRepo: userRepo, redis: redis}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.FIO) (uuid.UUID, error) {
	return s.EnrichAndSaveToDB(ctx, user.Name, user.Surname, user.Patronymic)
}

func (s *UserService) GetAllUsers(ctx context.Context, page int) ([]*models.User, error) {
	offset := (page - 1) * 30
	return s.userRepo.GetAllUsers(ctx, offset)
}

func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	return s.userRepo.GetUser(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input models.EnrichedFIO) error {
	return s.userRepo.UpdateProfile(ctx, userID, input)
}

func (s *UserService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}

func (s *UserService) EnrichAndSaveToDB(ctx context.Context, name, surname, patronymic string) (uuid.UUID, error) {
	age, err := api.GetAgifyAge(name)
	if err != nil {
		return uuid.Nil, err
	}
	gender, err := api.GetGenderizeGender(name)
	if err != nil {
		return uuid.Nil, err
	}
	nationality, err := api.GetNationalizeNationality(name)
	if err != nil {
		return uuid.Nil, err
	}
	user := &models.User{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}
	userID, err := s.userRepo.CreateUser(context.Background(), user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error save user in db, %v", err)
	}
	err = s.redis.SetData(userID.String(), user, 1*time.Hour)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error saving in cache, %v", err)
	}
	return userID, nil
}
