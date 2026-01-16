package userModel

type Role string

type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// Role         Role   `json:"role"`
	// AccessLevel  int64  `json:"access_level"`
}
