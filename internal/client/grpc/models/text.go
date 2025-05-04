package models

import "time"

// Структура для хранения учетных данных
type Text struct {
	ID        string
	Info      string
	Text      string
	CreatedAt time.Time
	Etag      string
}
