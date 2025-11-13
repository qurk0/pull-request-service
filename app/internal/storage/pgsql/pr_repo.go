package pgsql

type PullRequestRepository struct {
	db *PgDB
}

func newPullRequestRepo(db *PgDB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}
