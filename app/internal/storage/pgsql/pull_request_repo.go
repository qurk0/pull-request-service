package pgsql

type PullRequestRepository struct {
	db *PgDB
}

func NewPullRequestRepo(db *PgDB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}
