package service

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/repository"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenService struct {
	userRepo   *repository.UserRepository
	authenRepo *repository.AuthenRepository
}

func NewAuthenService(userRepo *repository.UserRepository, authenRepo *repository.AuthenRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo, authenRepo: authenRepo}
}

func (service *AuthenService) SignUpUser(signupDTO domain.SignUpDTO) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(signupDTO.Password), 14)
	if err != nil {
		return "", fmt.Errorf("cannot generate password hash: %w", err)
	}

	signupDTO.Password = string(bytes)

	newUser, err := service.authenRepo.CreateUser(signupDTO)
	if err != nil {
		return "", fmt.Errorf("cannot create user: %w", err)
	}

	return newUser.ID, nil
}

func (service *AuthenService) SignIn(loginDTO domain.SignInDTO) (string, error) {

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	userDTO, err := service.authenRepo.GetUserByCredential(loginDTO)
	if err != nil {
		return "", fmt.Errorf("cannot get user by credential: %w", err)
	}

	jwtClaims := jwt.RegisteredClaims{
		Subject:   userDTO.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("cannot creating access jwt. error : %s", err)
	}

	return accessTokenString, nil
}

func (service *AuthenService) Logout(accessToken *jwt.Token) error {

	var getUserDTO domain.GetUserDTO

	userID, err := accessToken.Claims.GetSubject()
	if err != nil {
		return fmt.Errorf("cannot get userID from jwt(accessToken). error : %s", err)
	}

	if err := service.userRepo.GetUserById(&getUserDTO, userID); err != nil {
		return err
	}

	issuedAt := jwt.NewNumericDate(time.Now())

	updateRefreshTokenDTO := domain.UpdateRefreshTokenDTO{
		UserID:       getUserDTO.ID,
		RefreshToken: "",
		IsRevoked:    true,
		DeviceInfo:   "",
		IpAddress:    "",
		IssuedAt:     issuedAt.Time,
		ExpiresAt:    issuedAt.Time,
	}

	if err := service.authenRepo.UpdateRefreshToken(updateRefreshTokenDTO); err != nil {
		return err
	}

	return nil
}

func (service *AuthenService) ReAccessToken(refreshTokenCookie *http.Cookie, newAccessToken *string, newRefreshToken *string, reAccessTokenDTO domain.ReAccessDTO) error {

	var refreshToken jwt.Token
	var userDTO domain.GetUserDTO

	err := ClaimsFromJWTString(&refreshToken, refreshTokenCookie.Value)
	if err != nil {
		return fmt.Errorf("cannot get claims from jwt string. error : %s", err.Error())
	}

	expires, err := refreshToken.Claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("cannot get expiresTime from jwt(refreshToken). error : %s", err.Error())
	}

	isExpired := expires.Before(time.Now())
	if isExpired {
		return fmt.Errorf("refresh token expired")
	}

	userID, err := refreshToken.Claims.GetSubject()
	if err != nil {
		return fmt.Errorf("cannot get userID from jwt(refreshToken). error : %s", err.Error())
	}

	if err := service.userRepo.GetUserById(&userDTO, userID); err != nil {
		return fmt.Errorf("cannot get user by id. error : %s", err.Error())
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("JWT_SECRET environment variable not set")
	}

	AccessTokenIssuedAt := jwt.NewNumericDate(time.Now())
	AccessTokenExpiresAt := jwt.NewNumericDate(time.Now().Add(time.Minute * 15))

	accessTokenClaims := jwt.RegisteredClaims{
		Subject:   userDTO.ID,
		IssuedAt:  AccessTokenIssuedAt,
		ExpiresAt: AccessTokenExpiresAt,
	}

	generatedAccessToken, err := GenerateJWTWithClaims(accessTokenClaims, secretKey)
	if err != nil {
		return fmt.Errorf("cannot create jwt accessToken. error : %s", err.Error())
	}

	*newAccessToken = *generatedAccessToken

	RefreshTokenIssuedAt := jwt.NewNumericDate(time.Now())
	RefreshTokenExpiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30))

	refreshTokenClaims := jwt.RegisteredClaims{
		Subject:   userDTO.ID,
		IssuedAt:  RefreshTokenIssuedAt,
		ExpiresAt: RefreshTokenExpiresAt,
	}

	generatedRefreshToken, err := GenerateJWTWithClaims(refreshTokenClaims, secretKey)
	if err != nil {
		return fmt.Errorf("cannot create jwt refreshToken. error : %s", err.Error())
	}

	*newRefreshToken = *generatedRefreshToken

	updateRefreshTokenDTO := domain.UpdateRefreshTokenDTO{
		UserID:       userDTO.ID,
		RefreshToken: *newRefreshToken,
		IsRevoked:    false,
		DeviceInfo:   reAccessTokenDTO.DeviceInfo,
		IpAddress:    reAccessTokenDTO.IpAddress,
		IssuedAt:     RefreshTokenIssuedAt.Time,
		ExpiresAt:    RefreshTokenExpiresAt.Time,
	}

	if err := service.authenRepo.UpdateRefreshToken(updateRefreshTokenDTO); err != nil {
		return err
	}

	return nil
}

func GenerateJWTWithClaims(claims jwt.RegisteredClaims, secretKey string) (*string, error) {

	claimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := claimsToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func ClaimsFromJWTString(refreshToken *jwt.Token, jwtString string) error {

	token, err := jwt.ParseWithClaims(
		jwtString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("refreshToken"), nil
		},
	)
	if err != nil {
		return err
	}

	*refreshToken = *token

	return nil
}
