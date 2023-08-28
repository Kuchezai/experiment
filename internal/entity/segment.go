package entity

import "time"

type Segment struct {
	Slug string
}

type SlugWithExpiredDate struct {
	Slug string
	ExpiredDate  time.Time
}
