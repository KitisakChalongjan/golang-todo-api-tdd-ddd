package valueobject

import "time"

type GetTodoVO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due"`
	UserID      string     `json:"user_id"`
}

type CreateTodoVO struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due"`
	UserID      string     `json:"user_id"`
}

type UpdateTodoVO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	Priority    string     `json:"priority"`
	Due         *time.Time `json:"due"`
}
