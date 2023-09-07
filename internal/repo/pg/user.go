package pg

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"experiment.io/internal/entity"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

var MaxTime = time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)

type UserPGX interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type UserRepository struct {
	db            UserPGX
	dirToStoreCSV string
}

func NewUserRepository(db UserPGX, dirToStoreCSV string) (*UserRepository, error) {
	if _, err := os.Stat(dirToStoreCSV); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dirToStoreCSV, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &UserRepository{db, dirToStoreCSV}, nil
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

func (r *UserRepository) Password(username string) (string, error) {
	op := "repo.pg.user.Password"

	query := `
	SELECT encrypted_pwd FROM users
	WHERE name = $1
	`
	var password string
	err := r.db.QueryRow(context.TODO(), query, username).Scan(&password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, entity.ErrInvalidNameOrPass)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return password, nil
}

// TODO : Remove the loop and enter everything in one big request

// Adds expire time only if ttl > 0, otherwise make it infinity
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
			return r.checkUserToSegmentError(op, err)
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
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

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

func (r *UserRepository) UsersHistoryInByDate(year int, month int) ([]entity.UserSegmentsHistory, error) {
	op := "repo.pg.user.UsersHistoryInByDate"

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0)

	query := `
	SELECT operation_id, user_id, segment_slug, isAdded, operation_date
	FROM segment_user_operations
	WHERE operation_date >= $1 AND operation_date <= $2
	`

	rows, err := r.db.Query(context.TODO(), query, firstDay, lastDay)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var history []entity.UserSegmentsHistory
	for rows.Next() {
		var hist entity.UserSegmentsHistory

		if err := rows.Scan(
			&hist.OperationID,
			&hist.UserID,
			&hist.SegmentSlug,
			&hist.IsAdded,
			&hist.Date,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		history = append(history, hist)
	}

	return history, nil
}

func (r *UserRepository) WriteHistoryToCSV(history []entity.UserSegmentsHistory, year int, month int) (string, error) {
	op := "usecase.user.WriteHistoryToCSV"
	filePath := fmt.Sprintf("%s/user_segments_history-%d-%d.csv", r.dirToStoreCSV, year, month)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"operation_id", "user_id", "segment_slug", "is_added", "date"}
	if err = writer.Write(header); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var rows [][]string
	for _, h := range history {
		row := []string{strconv.Itoa(h.OperationID), strconv.Itoa(h.UserID), h.SegmentSlug, strconv.FormatBool(h.IsAdded), h.Date.Format(time.RFC3339)}
		rows = append(rows, row)
	}
	fmt.Println(rows)
	if err = writer.WriteAll(rows); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return filePath, nil
}

func (r *UserRepository) checkUserToSegmentError(op string, err error) error {
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
