package model

type FileInfo struct {
	Login 	string	`db:"login"`	// owner login
	Id		string	`db:"id"`		// file id
}