package valueobject

type SignUpVO struct {
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	Roles         []string `json:"roles"`
	ProfileImgURL *string  `json:"profile_img_url"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
}

type SignInVO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
