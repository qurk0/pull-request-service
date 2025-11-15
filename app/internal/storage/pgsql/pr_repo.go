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

	GetPRByIDQuery = `SELECT id, pull_request_name, author_id, status, created_at, merged_at
	FROM pull_requests
	WHERE id = $1;`

	ReassignPRReviewerQuery = `UPDATE pull_requests_reviewers
	SET reviewer_id = $3
	WHERE pull_request_id = $1
		AND reviewer_id = $2;
	`

	GetReviewersQuery = `SELECT reviewer_id
	FROM pull_requests_reviewers
	WHERE pull_request_id = $1;`

	MergePRQuery = `UPDATE pull_requests
	SET 
		status = 'MERGED',
		merged_at = COALESCE(merged_at, NOW())
	WHERE id = $1
	RETURNING id, pull_request_name, author_id, status, created_at, merged_at;`
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

func (r *PullRequestRepository) GetPRByID(ctx context.Context, prID string) (models.PR, error) {
	var pr models.PR
	err := r.db.pool.QueryRow(ctx, GetPRByIDQuery, prID).Scan(&pr.PRID,
		&pr.PRName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		return models.PR{}, mapErr(err)
	}

	return pr, nil
}

func (r *PullRequestRepository) ReassignPRReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (models.PR, []string, error) {
	tx, err := r.db.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return models.PR{}, nil, mapErr(err)
	}
	defer tx.Rollback(ctx)

	var pr models.PR
	const GetPRForUpdate string = `
	SELECT id, pull_request_name, author_id, status, created_at, merged_at
	FROM pull_requests
	WHERE id = $1
	FOR UPDATE;
	`

	if err := tx.QueryRow(ctx, GetPRForUpdate, prID).Scan(&pr.PRID,
		&pr.PRName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt); err != nil {
		return models.PR{}, nil, mapErr(err)
	}

	if pr.Status != models.OpenStatus {
		return models.PR{}, nil, models.ErrPRMerged
	}

	pgc, err := tx.Exec(ctx, ReassignPRReviewerQuery, prID, oldReviewerID, newReviewerID)
	if err != nil {
		return models.PR{}, nil, mapErr(err)
	}
	if pgc.RowsAffected() == 0 {
		return models.PR{}, nil, models.ErrNotAssigned
	}

	rows, err := tx.Query(ctx, GetReviewersQuery, pr.PRID)
	if err != nil {
		return models.PR{}, nil, mapErr(err)
	}
	defer rows.Close()
	var reviewers []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return models.PR{}, nil, mapErr(err)
		}

		reviewers = append(reviewers, id)
	}

	if rows.Err() != nil {
		return models.PR{}, nil, mapErr(rows.Err())
	}

	if err := tx.Commit(ctx); err != nil {
		return models.PR{}, nil, mapErr(err)
	}

	return pr, reviewers, nil
}

func (r *PullRequestRepository) GetPRReviewers(ctx context.Context, prID string) ([]string, error) {
	rows, err := r.db.pool.Query(ctx, GetReviewersQuery, prID)
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
func (r *PullRequestRepository) MergePR(ctx context.Context, prID string) (models.PR, error) {
	var pr models.PR

	err := r.db.pool.QueryRow(ctx, MergePRQuery, prID).Scan(&pr.PRID,
		&pr.PRName,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt)
	if err != nil {
		return models.PR{}, mapErr(err)
	}

	return pr, nil
}
