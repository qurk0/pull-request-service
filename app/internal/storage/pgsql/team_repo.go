package pgsql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qurk0/pr-service/internal/domain/models"
)

const (
	CheckTeamExists = `SELECT EXISTS(
		SELECT 1 
		FROM teams 
		WHERE team_name = $1
	);
	`

	CreateTeamQuery = `INSERT INTO teams (team_name)
	VALUES ($1);
	`
)

type TeamRepository struct {
	db  *PgDB
	log *slog.Logger
}

func newTeamRepo(db *PgDB, log *slog.Logger) *TeamRepository {
	return &TeamRepository{db: db, log: log}
}

func (r *TeamRepository) CheckTeamExists(ctx context.Context, teamName string) (bool, error) {
	const op = "team_repo.CheckTeamExists"
	var exists bool

	err := r.db.pool.QueryRow(ctx, CheckTeamExists, teamName).Scan(&exists)
	if err != nil {
		r.log.Error(op, slog.String("error: failed to check team", err.Error()))
		return false, mapErr(err)
	}

	r.log.Debug(op, slog.String("success", "team check done successfully"))
	return exists, nil
}

func (r *TeamRepository) CreateTeamWithMembers(ctx context.Context, teamName string, members []models.TeamMember) (models.Team, error) {
	const op = "team_repo.CreateTeamWithMembers"

	tx, err := r.db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.log.Error(op, slog.String("error: failed to create team", err.Error()))
		return models.Team{}, mapErr(err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, CreateTeamQuery, teamName); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				r.log.Warn(op, slog.String("error: failed to create team", "team already exists"))
				return models.Team{}, models.ErrTeamExists
			}
		}
		r.log.Error(op, slog.String("error: failed to create team", err.Error()))
		return models.Team{}, mapErr(err)
	}

	if len(members) > 0 {
		var (
			query   = "INSERT INTO users (id, username, team_name, is_active) VALUES"
			counter = 1
			args    = make([]any, 0, len(members)*4)
		)

		for i, member := range members {
			r.log.Debug("member",
				slog.String("memID", member.Id),
				slog.Any("memActive", member.IsActive),
				slog.String("memUsername", member.Username),
			)
			if i == 0 {
				query += fmt.Sprintf("\n($%d, $%d, $%d, $%d)", counter, counter+1, counter+2, counter+3)
			} else {
				query += fmt.Sprintf(",\n($%d, $%d, $%d, $%d)", counter, counter+1, counter+2, counter+3)
			}
			args = append(args, member.Id, member.Username, teamName, member.IsActive)
			counter += 4
		}
		query += `
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active;`

		if _, err = tx.Exec(ctx, query, args...); err != nil {
			r.log.Error(op, slog.String("error: failed to create team", err.Error()))
			return models.Team{}, mapErr(err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		r.log.Error(op, slog.String("error: failed to create team", err.Error()))
		return models.Team{}, mapErr(err)
	}

	r.log.Debug(op, slog.String("success", "team created"))
	return models.Team{
		TeamName:    teamName,
		TeamMembers: members,
	}, nil
}
