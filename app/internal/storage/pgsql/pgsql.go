package pgsql

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgDB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connString string) (*PgDB, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return &PgDB{pool: pool}, nil
}

func (db *PgDB) Close() {
	db.pool.Close()
}

// Как единая структура, которую создаём в main.go для удобства
// Потом в соответствующий сервис передаём соответствующий репозиторий для удобства разработки
type Storage struct {
	Db          *PgDB
	User        *UserRepository
	Team        *TeamRepository
	PullRequest *PullRequestRepository
}

func NewStorage(db *PgDB, log *slog.Logger) *Storage {
	return &Storage{
		User:        newUserRepo(db, log),
		Team:        newTeamRepo(db, log),
		PullRequest: newPullRequestRepo(db, log),
	}
}
