package model

type Session struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// Role         Role   `json:"role"`
	Role        string `json:"role"`
	AccessLevel int    `json:"access_level"`
	IsActive    bool   `json:"is_active"`
}
