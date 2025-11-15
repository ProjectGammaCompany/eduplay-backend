package userModel

type UserAccess struct {
	UserId      string `json:"userId"`
	Role        string `json:"role"`
	AccessLevel int64  `json:"accessLevel"`
}
