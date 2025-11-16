package pgsql

import (
	"context"
	"log/slog"

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

	GetAnotherReviewerQuery = `SELECT id 
	FROM users
	WHERE id <> $1 -- старый ревьювер
		AND id <> $2 -- автор PR
		AND team_name = $3
		AND is_active = TRUE
		AND id NOT IN ( -- исключаем второго ревьювера если он есть
			SELECT reviewer_id
			FROM pull_requests_reviewers
			WHERE pull_request_id = $4
		);`
)

type UserRepository struct {
	db  *PgDB
	log *slog.Logger
}

func newUserRepo(db *PgDB, log *slog.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (models.User, error) {
	const op = "user_repo.GetUser"
	var user models.User
	if err := r.db.pool.QueryRow(ctx, GetUserQuery, userID).Scan(&user.Id,
		&user.Username,
		&user.TeamName,
		&user.IsActive); err != nil {
		r.log.Error(op, slog.String("error: failed to get user", err.Error()))
		return models.User{}, mapErr(err)
	}

	r.log.Debug(op, slog.String("success", "user got successfully"))
	return user, nil
}

func (r *UserRepository) UpdateUserIsActive(ctx context.Context, userId string, isActive bool) error {
	const op = "user_repo.UpdateUserIsActive"

	pgc, err := r.db.pool.Exec(ctx, UpdateUserIsActiveQuery, userId, isActive)
	if err != nil {
		r.log.Error(op, slog.String("error: failed to update user_is_active", err.Error()))
		return mapErr(err)
	}
	if pgc.RowsAffected() == 0 {
		r.log.Error(op, slog.String("error: failed to update user_is_active", "user not found"))
		return models.ErrNotFound
	}

	r.log.Debug(op, slog.String("success", "user updated successfully"))
	return nil
}

func (r *UserRepository) GetTeamMembers(ctx context.Context, teamName string) ([]models.TeamMember, error) {
	const op = "user_repo.GetTeamMembers"

	rows, err := r.db.pool.Query(ctx, GetTeamMembersQuery, teamName)
	if err != nil {
		r.log.Error(op, slog.String("error: failed to get team members", err.Error()))
		return nil, mapErr(err)
	}
	defer rows.Close()

	memberList := make([]models.TeamMember, 0)

	for rows.Next() {
		mem := models.TeamMember{}

		err := rows.Scan(&mem.Id, &mem.Username, &mem.IsActive)
		if err != nil {
			r.log.Error(op, slog.String("error: failed to get team members", err.Error()))
			return nil, mapErr(err)
		}

		memberList = append(memberList, mem)
	}

	if rows.Err() != nil {
		r.log.Error(op, slog.String("error: failed to get team members", rows.Err().Error()))
		return nil, mapErr(rows.Err())
	}

	r.log.Debug(op, slog.String("success", "team members got successfully"))
	return memberList, nil

}

func (r *UserRepository) GetReviewers(ctx context.Context, userID, teamName string) ([]string, error) {
	const op = "user_repo.GetReviewers"

	rows, err := r.db.pool.Query(ctx, GetActiveUsersQuery, teamName, userID)
	if err != nil {
		r.log.Error(op, slog.String("error: failed to get reviewers", err.Error()))
		return nil, mapErr(err)
	}
	defer rows.Close()

	var reviewers []string

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			r.log.Error(op, slog.String("error: failed to get reviewers", err.Error()))
			return nil, mapErr(err)
		}

		reviewers = append(reviewers, id)
	}

	if rows.Err() != nil {
		r.log.Error(op, slog.String("error: failed to get reviewers", rows.Err().Error()))
		return nil, mapErr(rows.Err())
	}

	r.log.Debug(op, slog.String("success", "reviewers got successfully"))
	return reviewers, nil
}

func (r *UserRepository) GetAnotherReviewers(ctx context.Context, prID, oldReviewerID, authorID string) ([]string, error) {
	const op = "user_repo.GetAnotherReviewers"
	var user models.User

	// Получаем ревьювера, которого хотим заменить, чтобы узнать его текущую команду
	if err := r.db.pool.QueryRow(ctx, GetUserQuery, oldReviewerID).Scan(&user.Id,
		&user.Username,
		&user.TeamName,
		&user.IsActive); err != nil {
		r.log.Error(op, slog.String("error: failed to get reviewers", err.Error()))
		return nil, mapErr(err)
	}

	rows, err := r.db.pool.Query(ctx, GetAnotherReviewerQuery, oldReviewerID, authorID, user.TeamName, prID)
	if err != nil {
		r.log.Error(op, slog.String("error: failed to get reviewers", err.Error()))
		return nil, mapErr(err)
	}
	defer rows.Close()

	var candidates []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			r.log.Error(op, slog.String("error: failed to get reviewers", err.Error()))
			return nil, mapErr(err)
		}

		candidates = append(candidates, id)
	}

	if rows.Err() != nil {
		r.log.Error(op, slog.String("error: failed to get reviewers", rows.Err().Error()))
		return nil, mapErr(rows.Err())
	}

	r.log.Debug(op, slog.String("success", "reviewers got successfully"))
	return candidates, nil
}
