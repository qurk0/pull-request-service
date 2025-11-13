package pgsql

type UserRepository struct {
	db *PgDB
}

func newUserRepo(db *PgDB) *UserRepository {
	return &UserRepository{db: db}
}
