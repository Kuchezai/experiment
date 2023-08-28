package entity

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInternalServer        = errors.New("internal server error")
	ErrSegmentNotFound       = errors.New("segment not found")
	ErrUserAlreadyExist      = errors.New("user already exist")
	ErrSegmentAlreadyExist   = errors.New("segment already exist")
	ErrSegmentsIntersect     = errors.New("added and removed segments intersect")
	ErrUserAlreadyAssigned   = errors.New("the user is already assigned this segment")
	ErrUserToSegmentNotFound = errors.New("the user is not assigned this segment")
)
