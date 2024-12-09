package models

type Currency struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
