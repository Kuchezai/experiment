package pg

import (
	"context"
	"errors"
	"fmt"

	"experiment.io/internal/entity"
	"experiment.io/pkg/storage/pg"
	"github.com/jackc/pgx/v5/pgconn"
)

type SegmentRepository struct {
	db *pg.Postgres
}

func NewSegmentRepository(db *pg.Postgres) *SegmentRepository {
	return &SegmentRepository{db}
}

func (r *SegmentRepository) NewSegment(seg entity.Segment) error {
	op := "repo.pg.segment.New"

	query := `
	INSERT INTO segments
	(slug) 
	VALUES($1)
	`

	if _, err := r.db.Exec(context.TODO(), query, seg.Slug); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == DuplicatePKErrCode {
			return fmt.Errorf("%s: %w", op, entity.ErrSegmentAlreadyExist)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Creates a segment and returns the user IDs assigned to it
func (r *SegmentRepository) NewSegmentWithAutoAssign(seg entity.Segment, percentAssigned int) ([]int, error) {
	op := "repo.pg.segment.NewWithAutoAssign"

	query := `
	SELECT * FROM create_segment_and_add_users($1, $2);
	`
	rows, err := r.db.Query(context.TODO(), query, seg.Slug, percentAssigned)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == DuplicatePKErrCode {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrSegmentAlreadyExist)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var isCreated bool
		var id int
		if err := rows.Scan(&id, &isCreated); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !isCreated {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrSegmentAlreadyExist)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *SegmentRepository) DeleteSegment(slug string) error {
	op := "repo.pg.segment.Delete"

	query := `
	DELETE FROM segments
	WHERE slug = $1
	`
	res, err := r.db.Exec(context.TODO(), query, slug)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrSegmentNotFound)
	}

	return nil
}
