package pgsql

import (
	"context"

	"github.com/qurk0/pr-service/internal/domain/models"
)

const (
	GetByReviewerQuery = `SELECT id, pull_request_name, author_id, status
	FROM pull_requests pr
	JOIN pull_requests_reviewers prr
		ON pr.id = prr.pull_request_id
	WHERE prr.reviewer_id = $1;
	`
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
