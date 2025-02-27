package entities

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type OAuthRequest struct {
	Code string `json:"code"`
}

type User struct {
	Id            int       `json:"id" db:"id"`
	Sub           string    `json:"sub" db:"sub"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	Picture       string    `json:"picture" db:"picture"`
	Refresh_token string    `json:"refresh_token" db:"refresh_token"`
	Expires_at    time.Time `json:"expires_at" db:"expires_at"`
	Scope         string    `json:"scope" db:"scope"`
	Created_at    time.Time `json:"created_at" db:"created_at"`
	Updated_at    time.Time `json:"updated_at" db:"updated_at"`
}

type TokenClaims struct {
	Id         int       `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Picture    string    `json:"picture"`
	Expires_at time.Time `json:"expires_at"`
	jwt.StandardClaims
}

type AccessToken struct {
	Id         int       `json:"id"`
	Token      string    `json:"token"`
	Expires_at time.Time `json:"expires_at"`
}
