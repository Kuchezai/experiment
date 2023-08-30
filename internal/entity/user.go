package entity

import "time"

type User struct {
	Name     string
	Password string
}

type UserSegmentsHistory struct {
	OperationID int
	UserID      int
	SegmentSlug string
	IsAdded     bool
	Date        time.Time
}
