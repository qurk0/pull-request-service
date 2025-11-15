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

	GetTeamMembersQuery = `SELECT id, username, is_active
	FROM users
	WHERE team_name = $1;
	`

	GetActiveUsersQuery = `SELECT id 
	FROM users
	WHERE team_name = $1
		AND id <> $2
		AND is_active = TRUE;`
)

type UserRepository struct {
	db *PgDB
}

func newUserRepo(db *PgDB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	if err := r.db.pool.QueryRow(ctx, GetUserQuery, userID).Scan(&user.Id,
		&user.Username,
		&user.TeamName,
		&user.IsActive); err != nil {
		return models.User{}, mapErr(err)
	}

	return user, nil
}

func (r *UserRepository) UpdateUserIsActive(ctx context.Context, userId string, isActive bool) error {
	pgc, err := r.db.pool.Exec(ctx, UpdateUserIsActiveQuery, userId, isActive)
	if err != nil {
		return mapErr(err)
	}
	if pgc.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *UserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]models.TeamMember, error) {
	rows, err := r.db.pool.Query(ctx, GetTeamMembersQuery, teamName)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	memberList := make([]models.TeamMember, 0)

	for rows.Next() {
		mem := models.TeamMember{}

		err := rows.Scan(&mem.Id, &mem.Username, &mem.IsActive)
		if err != nil {
			return nil, mapErr(err)
		}

		memberList = append(memberList, mem)
	}

	if rows.Err() != nil {
		return nil, mapErr(rows.Err())
	}

	return memberList, nil

}

func (r *UserRepository) GetReviewers(ctx context.Context, userID, teamName string) ([]string, error) {
	rows, err := r.db.pool.Query(ctx, GetActiveUsersQuery, teamName, userID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	var reviewers []string

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, mapErr(err)
		}

		reviewers = append(reviewers, id)
	}

	if rows.Err() != nil {
		return nil, mapErr(rows.Err())
	}

	return reviewers, nil
}
