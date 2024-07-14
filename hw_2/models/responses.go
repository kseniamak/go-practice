package models

type GetAccountRequest struct {
	Name string `json:"name"`
}

type GetAccountResponse struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}
