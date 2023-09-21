package repository

import (
	"EMTestTask/pkg/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.db.Exec(ctx, `INSERT INTO users VALUES ($1,$2, $3, $4, $5, $6, $7) RETURNING users.id`, id, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error saving user, %v", err)
	}
	return id, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, offset int) ([]*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) UpdateProfile(ctx context.Context, userID uuid.UUID, input model.EnrichedFIO) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, userID)
	if err != nil {
		return fmt.Errorf("error while delete profile in user repository: %v", err)
	}
	return nil
}
