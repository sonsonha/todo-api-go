package domain

import "time"

type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	IsDone    bool      `json:"is_done"`
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
