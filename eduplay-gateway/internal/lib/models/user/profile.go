package userModel

type Profile struct {
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	UserName string `json:"username"`
	// Password string `json:"password"`
}
