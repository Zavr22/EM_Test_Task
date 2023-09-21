package enrich

import (
	"EMTestTask/cache"
	"EMTestTask/pkg/api"
	"EMTestTask/pkg/model"
	"EMTestTask/web/rest/repository"
	"context"
	"fmt"
	"time"
)

func EnrichAndSaveToDB(name, surname, patronymic string, userRepo *repository.UserRepository, cache *cache.RedisClient) error {
	age, err := api.GetAgifyAge(name)
	if err != nil {
		return err
	}
	gender, err := api.GetGenderizeGender(name)
	if err != nil {
		return err
	}
	nationality, err := api.GetNationalizeNationality(name)
	if err != nil {
		return err
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
		return fmt.Errorf("error save user in db, %v", err)
	}
	err = cache.SetData(userID.String(), user, 1*time.Hour)
	if err != nil {
		return fmt.Errorf("error saving in cache, %v", err)
	}
	return nil
}
