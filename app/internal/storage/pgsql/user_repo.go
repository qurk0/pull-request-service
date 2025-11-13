package pgsql

type UserRepository struct {
	db *PgDB
}

func NewUserRepo(db *PgDB) *UserRepository {
	return &UserRepository{db: db}
}
