package model

type UserInfo struct {
	Email string `json:"email"`
	Hash  []byte `json:"-"`
}
