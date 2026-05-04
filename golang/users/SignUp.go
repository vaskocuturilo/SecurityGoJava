package users

type SignUp struct {
	UserName string `json:"userName"`
	Password []byte `json:"password"`
}
