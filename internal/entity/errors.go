package entity

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInternalServer        = errors.New("internal server error")
	ErrSegmentNotFound       = errors.New("segment not found")
	ErrInvalidPassString     = errors.New("invalid password string") // if it is not possible to hash the password
	ErrInvalidNameOrPass     = errors.New("invalid username or password")
	ErrInvalidToken          = errors.New("invalid or unspecified token")
	ErrUserAlreadyExist      = errors.New("user already exist")
	ErrInvalidAddedSegment   = errors.New("add_segments: ttl must be less or equal 366, greater or equal 0. slug must be provided")
	ErrSegmentAlreadyExist   = errors.New("segment already exist")
	ErrSegmentsIntersect     = errors.New("added and removed segments intersect")
	ErrUserAlreadyAssigned   = errors.New("the user is already assigned this segment")
	ErrUserToSegmentNotFound = errors.New("the user is not assigned this segment")
)
