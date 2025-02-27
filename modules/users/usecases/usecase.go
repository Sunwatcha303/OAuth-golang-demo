package usecases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Sunwatcha303/OAuth-golang-demo/configs"
	"github.com/Sunwatcha303/OAuth-golang-demo/modules/entities"
	"github.com/Sunwatcha303/OAuth-golang-demo/modules/users/repositories"
	"github.com/Sunwatcha303/OAuth-golang-demo/pkg/utils"
	"github.com/golang-jwt/jwt"
)

type UserUsecase struct {
	usersRepository *repositories.UserRepository
	cfg             *configs.Configs
}

func NewUsersUsecase(userRepository *repositories.UserRepository, cfg *configs.Configs) *UserUsecase {
	return &UserUsecase{
		usersRepository: userRepository,
		cfg:             cfg,
	}
}

func (u *UserUsecase) GetUrlOAuth() string {
	scope := "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=offline&prompt=consent", u.cfg.OAuth.ClientID, u.cfg.OAuth.RedirectUri, scope)
	return url
}

func (u *UserUsecase) GetToken(token *entities.TokenClaims) (newToken string, err error) {
	accessToken, err := u.usersRepository.GetAccessTokenByUserId(token.Id)
	if err != nil || time.Now().After(accessToken.Expires_at) {
		var refreshToken string
		if refreshToken, err = u.usersRepository.GetRefreshTokenByUserId(token.Id); err != nil {
			return "", fmt.Errorf("refresh token not found: %v", err)
		}
		data := url.Values{}
		data.Set("client_id", u.cfg.OAuth.ClientID)
		data.Set("client_secret", u.cfg.OAuth.ClientSecret)
		data.Set("refresh_token", refreshToken)
		data.Set("grant_type", "refresh_token")

		resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var tokenData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
			return "", err
		}

		accessToken, ok := tokenData["access_token"].(string)
		if !ok {
			return "", fmt.Errorf("access token not found")
		}

		expireIn, ok := tokenData["expires_in"].(float64)
		if !ok {
			return "", fmt.Errorf("expire not found")
		}

		expireTime := time.Now().Add(time.Duration(int(expireIn)))
		if err = u.usersRepository.SaveAccessToken(token.Id, accessToken, expireTime); err != nil {
			return "", err
		}
	}
	user, err := u.usersRepository.GetUserById(token.Id)
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &entities.TokenClaims{
		Id:         user.Id,
		Email:      user.Email,
		Name:       user.Name,
		Picture:    user.Picture,
		Expires_at: user.Expires_at,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	newToken, err = utils.GenerateToken(claims, u.cfg)
	if err != nil {
		return "", err
	}
	return
}

func (u *UserUsecase) GetNewToken(authCode string) (string, error) {
	data := url.Values{}
	data.Set("code", authCode)
	data.Set("client_id", u.cfg.OAuth.ClientID)
	data.Set("client_secret", u.cfg.OAuth.ClientSecret)
	data.Set("redirect_uri", u.cfg.OAuth.RedirectUri)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		return "", err
	}

	tokenId, ok := tokenData["id_token"].(string)
	if !ok {
		return "", fmt.Errorf("id_token missing or invalid")
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token missing or invalid")
	}

	expiresIn, ok := tokenData["expires_in"].(float64)
	if !ok {
		return "", fmt.Errorf("expires_in missing or invalid")
	}

	refreshToken, ok := tokenData["refresh_token"].(string)
	if !ok {
		return "", fmt.Errorf("fresh_token missing or invalid")
	}

	scope, ok := tokenData["scope"].(string)
	if !ok {
		return "", fmt.Errorf("scope missing or invalid")
	}

	expiresTime := time.Now().Add(time.Duration(expiresIn) * time.Second)

	token, _, err := new(jwt.Parser).ParseUnverified(tokenId, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var user entities.User

		if sub, ok := claims["sub"].(string); ok {
			user.Sub = sub
		}
		if email, ok := claims["email"].(string); ok {
			user.Email = email
		}
		if name, ok := claims["name"].(string); ok {
			user.Name = name
		}
		if picture, ok := claims["picture"].(string); ok {
			user.Picture = picture
		}

		user.Refresh_token = refreshToken
		user.Scope = scope
		user.Expires_at = expiresTime

		existUser, _ := u.usersRepository.GetUserBySup(user.Sub)
		var newId int
		if existUser != nil {
			user.Id = existUser.Id
			if err := u.usersRepository.UpdateUser(&user); err != nil {
				return "", fmt.Errorf("failed to update user: %v", err)
			}
		} else {
			if newId, err = u.usersRepository.CreateUser(&user); err != nil {
				return "", fmt.Errorf("failed to create user: %v", err)
			}
			user.Id = newId
		}

		if err := u.usersRepository.SaveAccessToken(user.Id, accessToken, expiresTime); err != nil {
			return "", fmt.Errorf("failed to save access token: %v", err)
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &entities.TokenClaims{
			Id:         user.Id,
			Email:      user.Email,
			Name:       user.Name,
			Picture:    user.Picture,
			Expires_at: user.Expires_at,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		newToken, err := utils.GenerateToken(claims, u.cfg)
		if err != nil {
			return "", fmt.Errorf("failed to sign the token: %v", err)
		}

		return newToken, nil
	}

	return "", fmt.Errorf("failed to extract claims from JWT token")
}

func (u *UserUsecase) VerifyAndExtractToken(tokenString string) (*entities.TokenClaims, error) {
	return utils.VerifyAndExtractToken(tokenString, u.cfg)
}
