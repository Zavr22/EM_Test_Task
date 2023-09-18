package rest

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

func (r *UserRepository) SaveUser(ctx context.Context, user *model.User) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.db.Exec(ctx, `INSERT INTO users VALUES ($1,$2, $3, $4, $5, $6, $7) RETURNING users.id`, id, user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error saving user, %v", err)
	}
	return id, nil
}
