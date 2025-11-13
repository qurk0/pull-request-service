package pgsql

type TeamRepository struct {
	db *PgDB
}

func newTeamRepo(db *PgDB) *TeamRepository {
	return &TeamRepository{db: db}
}
