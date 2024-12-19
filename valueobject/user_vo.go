package valueobject

import (
	"time"
)

type GetUserVO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ProfileImgURL *string   `json:"profile_img_url"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
	Roles         []string  `json:"roles"`
}

type UpdateUserVO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	Roles         []string `json:"roles"`
	ProfileImgURL *string  `json:"profile_img_url"`
	Password      string   `json:"password"`
}
