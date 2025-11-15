package userModel

type ChangePasswordRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
