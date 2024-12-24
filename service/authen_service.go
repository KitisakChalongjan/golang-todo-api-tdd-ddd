package service

import (
	"fmt"
	"golang-todo-api-tdd-ddd/helper"
	"golang-todo-api-tdd-ddd/repository"
	"golang-todo-api-tdd-ddd/valueobject"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenService struct {
	userRepo repository.IUserRepository
}

func NewAuthenService(userRepo repository.IUserRepository) *AuthenService {
	return &AuthenService{userRepo: userRepo}
}

func (service *AuthenService) SignUp(signUpVO valueobject.SignUpVO) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(signUpVO.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to generate password hash: %w", err)
	}

	hashedPasswordSignUpVO := signUpVO
	hashedPasswordSignUpVO.Password = string(bytes)

	user, err := service.userRepo.CreateUser(hashedPasswordSignUpVO)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (service *AuthenService) SignIn(signInVO valueobject.SignInVO, secretKey string) (string, error) {

	getUserVO, err := service.userRepo.GetUserByCredential(signInVO)
	if err != nil {
		return "", err
	}

	jwtClaims := jwt.MapClaims{
		"sub":   getUserVO.ID,
		"roles": getUserVO.Roles,
		"iat":   jwt.NewNumericDate(time.Now()),
		"exp":   jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
	}

	accessTokenString, err := helper.GenerateAccessTokenWithClaims(jwtClaims, secretKey)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
