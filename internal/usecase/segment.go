package usecase

import "experiment.io/internal/entity"

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
	return uc.r.NewSegment(seg)
}

func (uc *SegmentUsecase) DeleteSegment(slug string) error {
	return uc.r.DeleteSegment(slug)
}
