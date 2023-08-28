package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"experiment.io/internal/entity"
	"experiment.io/pkg/storage/pg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

var MaxTime = time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)

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

// TODO : Remove the loop and enter everything in one big request

// Adds expire time only if ttl > 0, otherwise make it indefinite
func (r *UserRepository) AddUserSegments(userID int, added []entity.SlugWithExpiredDate) error {
	op := "repo.pg.user.AddUserSegments"

	tx, err := r.db.Begin(context.TODO())
	defer tx.Rollback(context.TODO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `
	INSERT INTO segments_to_users
	(segment_slug, user_id, expiration_date)
	VALUES ($1, $2, CASE WHEN $3 > NOW() THEN $3 ELSE 'infinity' END)
	`
	for _, segmentToAdd := range added {
		if _, err := tx.Exec(context.TODO(), query, segmentToAdd.Slug, userID, segmentToAdd.ExpiredDate); err != nil {
			fmt.Println(err)
			return r.checkUserToSegmentErr(op, err)
		}
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *UserRepository) RemoveUserSegments(userID int, removed []string) error {
	op := "repo.pg.user.RemoveUserSegments"

	tx, err := r.db.Begin(context.TODO())
	defer tx.Rollback(context.TODO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `
	DELETE FROM segments_to_users
	WHERE user_id = $1 AND segment_slug = $2
	`
	for _, segmentToRemove := range removed {

		res, err := tx.Exec(context.TODO(), query, userID, segmentToRemove)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if res.RowsAffected() == 0 {
			return fmt.Errorf("%s: %w", op, entity.ErrUserToSegmentNotFound)
		}
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (r *UserRepository) UserSegments(userID int) ([]entity.SlugWithExpiredDate, error) {
	op := "repo.pg.user.UserSegments"

	query := `
	SELECT segment_slug, expiration_date FROM segments_to_users
	WHERE user_id = $1 AND expiration_date > NOW()
	`

	rows, err := r.db.Query(context.TODO(), query, userID)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segments []entity.SlugWithExpiredDate
	for rows.Next() {
		var seg entity.SlugWithExpiredDate
		var expirationDate pq.NullTime // needed in order to scan infinity time
		if err := rows.Scan(
			&seg.Slug,
			&expirationDate,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if expirationDate.Valid {
			seg.ExpiredDate = expirationDate.Time
		} else {
			seg.ExpiredDate = MaxTime // set max allowed time if expired time is infinite
		}

		segments = append(segments, seg)
	}

	return segments, nil

}

func (r *UserRepository) checkUserToSegmentErr(op string, err error) error {
	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	switch {
	case pgErr.Code == NonExistentFKErrCode && pgErr.ConstraintName == InvalidSegmentFK:
		err = entity.ErrSegmentNotFound
	case pgErr.Code == NonExistentFKErrCode && pgErr.ConstraintName == InvalidUserFK:
		err = entity.ErrUserNotFound
	case pgErr.Code == DuplicatePKErrCode:
		err = entity.ErrUserAlreadyAssigned
	}
	return fmt.Errorf("%s: %w", op, err)
}
