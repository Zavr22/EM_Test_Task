package graphql

import (
	"EMTestTask/pkg/model"
	"EMTestTask/web/rest/repository"
	"context"
	"github.com/google/uuid"
)

type Resolver struct {
	userRepo *repository.UserRepository
}

func (r *Resolver) Mutation() MutationResolver {
	return nil
}

func (r *Resolver) Query() QueryResolver {
	return nil
}

func NewResolver(userRepo *repository.UserRepository) *Resolver {
	return &Resolver{userRepo: userRepo}
}

func (r *Resolver) Query_getUsers(ctx context.Context, page int) ([]*model.User, error) {
	users, err := r.userRepo.GetAllUsers(ctx, page)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Resolver) Query_getUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	user, err := r.userRepo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *Resolver) Mutation_createUser(ctx context.Context, name string, surname string, patronymic string, age int, gender string, nationality string) (uuid.UUID, error) {
	user := &model.User{
		Name:        name,
		Surname:     surname,
		Patronymic:  patronymic,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}
	id, err := r.userRepo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *Resolver) Mutation_updateUser(ctx context.Context, id uuid.UUID, name *string, surname *string, patronymic *string, age *int, gender *string, nationality *string) (uuid.UUID, error) {
	fio := &model.FIO{
		Name:       *name,
		Surname:    *surname,
		Patronymic: *patronymic,
	}
	input := model.EnrichedFIO{
		FIO:         *fio,
		Age:         *age,
		Gender:      *gender,
		Nationality: *nationality,
	}
	if err := r.userRepo.UpdateProfile(ctx, id, input); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *Resolver) Mutation_deleteUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	if err := r.userRepo.DeleteProfile(ctx, id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
