package model

type UserInfo struct {
	Email string `db:"email"`
	Hash  []byte `db:"password_hash"`
}
