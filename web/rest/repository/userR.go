package repository

import (
	"EMTestTask/cache"
	"EMTestTask/internal/enrich"
	"EMTestTask/pkg/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	client *cache.RedisClient
}

func NewUserRepository(db *pgxpool.Pool, client *cache.RedisClient) *UserRepository {
	return &UserRepository{db: db, client: client}
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
	var users []*model.User

	rowsUsers, errUsers := r.db.Query(ctx, "SELECT id, name, surname, patronymic, age, gender, nationality FROM users LIMIT 30 OFFSET $1", offset)
	if errUsers != nil {
		return users, fmt.Errorf("error while getting users, %s", errUsers)
	}
	defer rowsUsers.Close()
	for rowsUsers.Next() {
		var user model.User
		errScan := rowsUsers.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
		if errScan != nil {
			return users, fmt.Errorf("get users scan rows error %w", errScan)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `SELECT id, name, surname, patronymic, age, gender, nationality FROM users WHERE id=$1`, userID).Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
	if err != nil {
		return model.User{}, fmt.Errorf("error get user by id, %v", err)
	}
	return user, nil
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

func (r *UserRepository) SaveUser(ctx context.Context, user *model.FIO) error {
	err := enrich.EnrichAndSaveToDB(user.Name, user.Surname, user.Patronymic, r, r.client)
	if err != nil {
		return fmt.Errorf("error saving user rest, %v", err)
	}
	return nil
}