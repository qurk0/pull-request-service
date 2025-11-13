package pgsql

type TeamRepository struct {
	db *PgDB
}

func NewTeamRepo(db *PgDB) *TeamRepository {
	return &TeamRepository{db: db}
}
