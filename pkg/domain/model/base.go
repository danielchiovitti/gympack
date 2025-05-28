package model

import (
	"time"
)

type BaseModel struct {
	CreatedAt   time.Time
	CreatedById string
	UpdatedAt   time.Time
	UpdatedById string
	Deleted     bool
	DeletedAt   time.Time
	DeletedById string
}
