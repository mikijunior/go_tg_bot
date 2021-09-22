package models

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}