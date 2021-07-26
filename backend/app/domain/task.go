package domain

import "time"

type Task struct {
	ID        int
	UserID    int
	Title     string
	Note      string
	Completed int
	CreatedAt time.Time
	UpdatedAt time.Time
}
