package userModel

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
