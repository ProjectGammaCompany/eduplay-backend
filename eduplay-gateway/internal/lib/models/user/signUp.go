package userModel

type SignUpRequest struct {
	Name         string `json:"name" validate:"required"`
	Surname      string `json:"surname" validate:"required"`
	Email        string `json:"email" validate:"required"`
	Organization string `json:"organization" validate:"required"`
	Password     string `json:"password" validate:"required"`
	Phone        string `json:"phone" validate:"required"`
}
