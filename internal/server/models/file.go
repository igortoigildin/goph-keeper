package model

type FileInfo struct {
	Login string `db:"login"`   // owner login
	Id    string `db:"data_id"` // file id
}
