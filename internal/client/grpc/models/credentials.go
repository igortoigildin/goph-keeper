package models

import "time"

// Структура для хранения учетных данных
type Credential struct {
	ID        string
	Service   string
	Username  string
	Password  string
	CreatedAt time.Time
	Etag      string
}
