package model

type UserInfo struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Hash  []byte `json:"-"`
}
