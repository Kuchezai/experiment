package usecase

import (
	"fmt"

	"experiment.io/internal/entity"
)

type SegmentRepo interface {
	NewSegment(seg entity.Segment) error
	DeleteSegment(slug string) error
}

type SegmentUsecase struct {
	r SegmentRepo
}

func NewSegmentUsecase(r SegmentRepo) *SegmentUsecase {
	return &SegmentUsecase{r}
}

func (uc *SegmentUsecase) NewSegment(seg entity.Segment) error {
	op := "usecase.segment.New"
	if err := uc.r.NewSegment(seg); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (uc *SegmentUsecase) DeleteSegment(slug string) error {
	op := "usecase.segment.Delete"
	if err := uc.r.DeleteSegment(slug); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
