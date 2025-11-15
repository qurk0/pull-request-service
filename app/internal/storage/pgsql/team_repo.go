package pgsql

import "context"

const (
	CheckTeamExists = `SELECT 1 
	FROM teams 
	WHERE team_name = $1;
	`
)

type TeamRepository struct {
	db *PgDB
}

func newTeamRepo(db *PgDB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CheckTeamExists(ctx context.Context, teamName string) error {
	var exists int
	if err := r.db.pool.QueryRow(ctx, CheckTeamExists, teamName).Scan(&exists); err != nil {
		return mapErr(err)
	}

	return nil
}
