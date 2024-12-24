package service_test

import (
	"errors"
	"fmt"
	"golang-todo-api-tdd-ddd/domain"
	"golang-todo-api-tdd-ddd/helper"
	"golang-todo-api-tdd-ddd/service"
	"golang-todo-api-tdd-ddd/valueobject"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccessSignUp1(t *testing.T) {

	fmt.Println("testing TestSuccessSignUp1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signUpVO := valueobject.SignUpVO{
		Name:     "test",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
		Username: "username",
		Password: "password",
	}

	mockRepo.On("CreateUser", mock.Anything).Return(domain.User{ID: "xxxx-xxxx-xxxx-xxxx"}, nil)

	userID, err := authenService.SignUp(signUpVO)

	assert.NoError(t, err)
	assert.Equal(t, "xxxx-xxxx-xxxx-xxxx", userID)
}

func TestFailSignUp1(t *testing.T) {

	fmt.Println("testing TestFailSignUp1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signUpVO := valueobject.SignUpVO{
		Name:     "test",
		Email:    "test@example.com",
		Roles:    []string{"admin"},
		Username: "username",
		Password: "password",
	}

	mockRepo.On("CreateUser", mock.Anything).Return(domain.User{ID: "xxxx-xxxx-xxxx-xxxx"}, errors.New("SignUp error!"))

	userID, err := authenService.SignUp(signUpVO)

	assert.Error(t, err)
	assert.Equal(t, "", userID)
	assert.EqualError(t, errors.New("SignUp error!"), err.Error())

	mockRepo.AssertExpectations(t)
}

func TestSuccessSignIn1(t *testing.T) {

	fmt.Println("testing TestSuccessSignIn1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signInVO := valueobject.SignInVO{
		Username: "username",
		Password: "password",
	}

	secretKey := "golang-todo-api-tdd-ddd"

	mockRepo.On("GetUserByCredential", mock.Anything).Return(valueobject.GetUserVO{ID: "xxxx-xxxx-xxxx-xxxx", Name: "name", Email: "test@example.com", Roles: []string{"guest"}}, nil)

	accessTokenString, err := authenService.SignIn(signInVO, secretKey)

	tokenClaims, _ := helper.TokenClaimsFromAccessTokenString(accessTokenString, secretKey)

	sub, _ := tokenClaims.Claims.GetSubject()

	assert.NoError(t, err)
	assert.Equal(t, "xxxx-xxxx-xxxx-xxxx", sub)

	mockRepo.AssertExpectations(t)
}

func TestFailSignIn1(t *testing.T) {

	fmt.Println("testing TestFailSignIn1()....")

	mockRepo := new(MockUserRepository)

	authenService := service.NewAuthenService(mockRepo)

	signInVO := valueobject.SignInVO{
		Username: "username",
		Password: "password",
	}

	secretKey := "golang-todo-api-tdd-ddd"

	mockRepo.On("GetUserByCredential", mock.Anything).Return(valueobject.GetUserVO{}, errors.New("SignIn error!"))

	accessTokenString, err := authenService.SignIn(signInVO, secretKey)

	assert.Error(t, err)
	assert.Equal(t, "", accessTokenString)
	assert.EqualError(t, errors.New("SignIn error!"), err.Error())

	mockRepo.AssertExpectations(t)
}
