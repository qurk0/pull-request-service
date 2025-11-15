package pgsql

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

const (
	GetUserQuery = `SELECT id, username, team_name, is_active
	FROM users
	WHERE id = $1;
	`

	UpdateUserIsActiveQuery = `UPDATE users
	SET is_active = $2
	WHERE id = $1;
	`
)

type UserRepository struct {
	db *PgDB
}

func newUserRepo(db *PgDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	if err := r.db.pool.QueryRow(ctx, GetUserQuery, userID).Scan(&user.Id,
		&user.Username,
		&user.TeamName,
		&user.IsActive); err != nil {
		return nil, mapErr(err)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserIsActive(ctx context.Context, user *models.User) error {
	pgc, err := r.db.pool.Exec(ctx, user.Id, user.IsActive)
	if err != nil {
		return mapErr(err)
	}
	if pgc.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}
