package service

import (
	"fmt"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/valueobject"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenService struct {
	userRepo     *repository.UserRepository
	roleRepo     *repository.RoleRepository
	userRoleRepo *repository.UserRoleRepository
}

func NewAuthenService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, userRoleRepo *repository.UserRoleRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo, roleRepo: roleRepo, userRoleRepo: userRoleRepo}
}

func (service *AuthenService) SignUp(signupDTO valueobject.SignUpVO) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(signupDTO.Password), 14)
	if err != nil {
		return "", fmt.Errorf("cannot generate password hash: %w", err)
	}

	signupDTO.Password = string(bytes)

	newUser, err := service.userRepo.CreateUser(signupDTO)
	if err != nil {
		return "", fmt.Errorf("cannot create user: %w", err)
	}

	return newUser.ID, nil
}

func (service *AuthenService) SignIn(signInVO valueobject.SignInVO) (string, error) {

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	getUserVO, err := service.userRepo.GetUserByCredential(signInVO)
	if err != nil {
		return "", fmt.Errorf("cannot get user by credential: %w", err)
	}

	jwtClaims := jwt.MapClaims{
		"sub":   getUserVO.ID,
		"roles": getUserVO.Roles,
		"iat":   jwt.NewNumericDate(time.Now()),
		"exp":   jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("cannot creating access jwt. error : %s", err)
	}

	return accessTokenString, nil
}

func GenerateJWTWithClaims(claims jwt.RegisteredClaims, secretKey string) (*string, error) {

	claimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := claimsToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &signedToken, nil
}

func ClaimsTokenFromJWTString(jwtString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(
		jwtString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("refreshToken"), nil
		},
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}
