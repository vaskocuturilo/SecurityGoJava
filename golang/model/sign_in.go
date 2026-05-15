package model

type SignIn struct {
	UserName string `json:"userName"`
	Password []byte `json:"password"`
}
