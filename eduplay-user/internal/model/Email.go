package model

type EmailRequest struct {
	To        string   `json:"to"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Cc        []string `json:"cc,omitempty"`
	Broadcast *bool    `json:"broadcast,omitempty"`
}
