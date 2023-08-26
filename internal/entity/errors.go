package entity

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrSegmentNotFound     = errors.New("segment not found")
	ErrUserAlreadyExist    = errors.New("user already exist")
	ErrSegmentAlreadyExist = errors.New("segment already exist")
)
