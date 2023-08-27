package pg

import (
	"context"
	"errors"
	"fmt"

	"experiment.io/internal/entity"
	"experiment.io/pkg/storage/pg"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db *pg.Postgres
}

func NewUserRepository(db *pg.Postgres) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) NewUser(u entity.User) (int, error) {
	op := "repo.pg.user.New"

	query := `
	INSERT INTO users
	(name, encrypted_pwd) 
	VALUES($1, $2)
	RETURNING id
	`
	var id int
	err := r.db.QueryRow(context.TODO(), query, u.Name, u.Password).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == DuplicatePKErrCode {
			return 0, fmt.Errorf("%s: %w", op, entity.ErrUserAlreadyExist)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
