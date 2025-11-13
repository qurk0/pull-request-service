package pgsql

import (
	"context"
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
	User        *UserRepository
	Team        *TeamRepository
	PullRequest *PullRequestRepository
}

func NewStorage(db *PgDB) *Storage {
	return &Storage{
		User:        NewUserRepo(db),
		Team:        NewTeamRepo(db),
		PullRequest: NewPullRequestRepo(db),
	}
}
