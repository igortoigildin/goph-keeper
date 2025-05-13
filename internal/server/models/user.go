package model

type UserInfo struct {
	Login string `db:"login"`
	Hash  []byte `db:"password_hash"`
}
