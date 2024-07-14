package models

type Account struct {
	Name   string `json: "name"`
	Amount int    `json: "amount"`
}
