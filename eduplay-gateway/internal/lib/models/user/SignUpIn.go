package userModel

type SignUpIn struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
}
