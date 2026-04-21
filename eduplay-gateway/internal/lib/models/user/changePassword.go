package userModel

type ChangePasswordRequest struct {
	Code       string `json:"code" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RepeatPass string `json:"repeatPassword" validate:"required"`
}

type Email struct {
	Email string `json:"email" validate:"required"`
}

type Code struct {
	Code string `json:"code" validate:"required"`
}
