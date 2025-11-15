package userModel

type ChangeUserData struct {
	Name                   string `json:"name"`
	Surname                string `json:"surname"`
	JobTitle               string `json:"jobTitle"`
	Organisation           string `json:"organisation"`
	Phone                  string `json:"phone"`
	Email                  string `json:"email"`
	City                   string `json:"city"`
	ShortOrganisationTitle string `json:"shortOrganisationTitle"`
	INN                    string `json:"INN"`
	OrganisationType       string `json:"organisationType"`
}
