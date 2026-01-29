package userModel

type Role string

type Credentials struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	// Role         Role   `json:"role"`
	// AccessLevel  int64  `json:"access_level"`
}
