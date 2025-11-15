package pgsql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qurk0/pr-service/internal/domain/models"
)

const (
	GetByReviewerQuery = `SELECT id, pull_request_name, author_id, status
	FROM pull_requests pr
	JOIN pull_requests_reviewers prr
		ON pr.id = prr.pull_request_id
	WHERE prr.reviewer_id = $1;
	`

	CreatePRQuery = `INSERT INTO pull_requests (id, pull_request_name, author_id, status)
	VALUES ($1, $2, $3, 'OPEN')
	RETURNING id, pull_request_name, author_id, status, created_at, merged_at;`

	AddReviewersQuery = `INSERT INTO pull_requests_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2);`
)

type PullRequestRepository struct {
	db *PgDB
}

func newPullRequestRepo(db *PgDB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) GetByReviewer(ctx context.Context, userID string) ([]models.PRShort, error) {
	rows, err := r.db.pool.Query(ctx, GetByReviewerQuery, userID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	prList := make([]models.PRShort, 0)
	for rows.Next() {
		pr := models.PRShort{}

		err := rows.Scan(
			&pr.PRID,
			&pr.PRName,
			&pr.AuthorID,
			&pr.Status,
		)

		if err != nil {
			return nil, mapErr(err)
		}

		prList = append(prList, pr)
	}

	if rows.Err() != nil {
		return nil, mapErr(rows.Err())
	}

	return prList, nil
}

func (r *PullRequestRepository) CreatePR(ctx context.Context, prID, prName, authorID string, reviewers []string) (models.PR, error) {
	var pr models.PR

	tx, err := r.db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.PR{}, mapErr(err)
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, CreatePRQuery, prID, prName, authorID).Scan(&pr.PRID,
		&pr.PRName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.PR{}, models.ErrPRExists
			}
		}
		return models.PR{}, mapErr(err)
	}

	pr.AssignedReviewers = make([]string, 0, len(reviewers))
	for _, id := range reviewers {
		if _, err := tx.Exec(ctx, AddReviewersQuery, pr.PRID, id); err != nil {
			return models.PR{}, mapErr(err)
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, id)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.PR{}, mapErr(err)
	}

	return pr, nil
}
