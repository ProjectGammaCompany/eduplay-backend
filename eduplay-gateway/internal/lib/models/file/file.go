package models

type File struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Files struct {
	Files []File `json:"files"`
}
