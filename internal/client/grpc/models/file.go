package models

import "time"

// Структура для хранения данных о файле
type File struct {
	ID        string
	Filename  string
	Data      []byte
	Password  string
	UpdatedAt time.Time
}
