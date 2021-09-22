package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	ChatId   int    `json:"chat_id"`
}
