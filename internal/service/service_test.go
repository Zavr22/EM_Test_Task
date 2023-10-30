package service_test

import (
	"context"
	"errors"
	"github.com/Zavr22/EMTestTask/internal/service"
	"github.com/Zavr22/EMTestTask/internal/service/mocks"
	"github.com/Zavr22/EMTestTask/pkg/models"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)

	userService := service.NewUserService(ctrl)

	ctx := context.Background()
	name := "John"
	surname := "Doe"
	patronymic := "Smith"

	expectedUserID := uuid.Nil

	mockUserRepo.EXPECT().CreateUser(ctx, &models.User{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Gender:      "male",
		Age:         72,
		Nationality: "IE",
	}).Return(expectedUserID, nil)

	userID, err := userService.CreateUser(ctx, &models.FIO{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if userID != expectedUserID {
		t.Errorf("expected userID: %s, got: %s", expectedUserID, userID)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)
	mockRedis := cache.NewRedisClient(rdb)

	userService := service.NewUserService(mockUserRepo, mockRedis)

	ctx := context.Background()
	page := 2
	limit := 2

	expectedUsers := []*models.User{
		{ID: uuid.New(), Name: "John", Surname: "Doe"},
		{ID: uuid.New(), Name: "Jane", Surname: "Smith"},
	}

	mockUserRepo.EXPECT().GetAllUsers(ctx, 30).Return(expectedUsers, nil)

	users, err := userService.GetAllUsers(ctx, page, limit)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(users) != len(expectedUsers) {
		t.Errorf("expected %d users, got: %d", len(expectedUsers), len(users))
	}
	for i := range users {
		if users[i].ID != expectedUsers[i].ID {
			t.Errorf("expected user ID: %s, got: %s", expectedUsers[i].ID, users[i].ID)
		}
		if users[i].Name != expectedUsers[i].Name {
			t.Errorf("expected user Name: %s, got: %s", expectedUsers[i].Name, users[i].Name)
		}
		if users[i].Surname != expectedUsers[i].Surname {
			t.Errorf("expected user Surname: %s, got: %s", expectedUsers[i].Surname, users[i].Surname)
		}
	}
}

func TestUserService_GetUser(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)
	mockRedis := cache.NewRedisClient(rdb)

	userService := service.NewUserService(mockUserRepo, mockRedis)

	ctx := context.Background()
	userID := uuid.New()

	expectedUser := models.User{
		ID:      userID,
		Name:    "John",
		Surname: "Doe",
	}

	mockUserRepo.EXPECT().GetUser(ctx, userID).Return(expectedUser, nil)

	user, err := userService.GetUser(ctx, userID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID: %s, got: %s", expectedUser.ID, user.ID)
	}
	if user.Name != expectedUser.Name {
		t.Errorf("expected user Name: %s, got: %s", expectedUser.Name, user.Name)
	}
	if user.Surname != expectedUser.Surname {
		t.Errorf("expected user Surname: %s, got: %s", expectedUser.Surname, user.Surname)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)
	mockRedis := cache.NewRedisClient(rdb)

	userService := service.NewUserService(mockUserRepo, mockRedis)

	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.EXPECT().DeleteProfile(ctx, userID).Return(nil)

	err := userService.DeleteProfile(ctx, userID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)
	mockRedis := cache.NewRedisClient(rdb)

	userService := service.NewUserService(mockUserRepo, mockRedis)

	ctx := context.Background()
	userID := uuid.New()
	name := "John"
	surname := "Doe"
	patronymic := "Smith"
	fio := models.FIO{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
	}
	mockUserRepo.EXPECT().UpdateProfile(ctx, userID, models.EnrichedFIO{
		FIO: fio,
	}).Return(nil)

	err := userService.UpdateProfile(ctx, userID, models.EnrichedFIO{
		FIO: fio,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserService_GetUserWithError(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)
	mockRedis := cache.NewRedisClient(rdb)

	userService := service.NewUserService(mockUserRepo, mockRedis)

	ctx := context.Background()
	userID := uuid.New()

	expectedError := errors.New("user not found")
	mockUserRepo.EXPECT().GetUser(ctx, userID).Return(models.User{}, expectedError)
	_, err := userService.GetUser(ctx, userID)
	if err != expectedError {
		t.Errorf("expected error: %v, got: %v", expectedError, err)
	}
}

func TestUserService_GetAllUsersWithError(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_service.NewMockUser(ctrl)
	mockRedis := cache.NewRedisClient(rdb)

	userService := service.NewUserService(mockUserRepo, mockRedis)

	ctx := context.Background()
	page := 2

	expectedError := errors.New("failed to get users")
	mockUserRepo.EXPECT().GetAllUsers(ctx, 30).Return(nil, expectedError)
	_, err := userService.GetAllUsers(ctx, page)
	if err != expectedError {
		t.Errorf("expected error: %v, got: %v", expectedError, err)
	}
}
