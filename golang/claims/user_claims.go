package claims

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Name          string `json:"user_name"`
	Email         string `json:"user_email"`
	TokenType     string `json:"token_type"`
	Role          string `json:"role"`
	SecurityStamp string `json:"security_stamp"`
	jwt.RegisteredClaims
}
