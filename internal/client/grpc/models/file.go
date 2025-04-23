package models

import "time"

// Структура для хранения данных о файле
type File struct {
	ID        string
	Filename  string
	Data      []byte
	UpdatedAt time.Time
	Info      string
}
