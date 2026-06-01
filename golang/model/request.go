package model

type Request struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
