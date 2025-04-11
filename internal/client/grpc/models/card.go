package models

import "time"

// Структура для хранения банк данных
type BankDetails struct {
	ID         string
	Info       string
	CardNumber string
	Cvc        string
	ExpDate    string
	CreatedAt  time.Time
}
