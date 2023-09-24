package enrich

import (
	"context"
	"fmt"
	"github.com/Zavr22/EMTestTask/cache"
	"github.com/Zavr22/EMTestTask/pkg/api"
	"github.com/Zavr22/EMTestTask/pkg/model"
	"github.com/Zavr22/EMTestTask/web/rest/repository"
	"github.com/google/uuid"
	"time"
)

func EnrichAndSaveToDB(name, surname, patronymic string, userRepo *repository.UserRepository, cache *cache.RedisClient) (uuid.UUID, error) {
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
	user := &model.User{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}
	userID, err := userRepo.CreateUser(context.Background(), user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error save user in db, %v", err)
	}
	err = cache.SetData(userID.String(), user, 1*time.Hour)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error saving in cache, %v", err)
	}
	return userID, nil
}
