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

const (
	DuplicatePKErrCode = "23505"
)

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
