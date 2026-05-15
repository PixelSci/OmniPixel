package usecase

import (
	"omni-pixel/domain"
)

type AccountUseCase struct {
	userRepository        domain.UserRepository
	accessTokenSecret     string
	accessTokenExpiryHour int
}

func NewAccountUseCase(userRepository domain.UserRepository, accessTokenSecret string, accessTokenExpiryHour int) *AccountUseCase {
	return &AccountUseCase{
		userRepository:        userRepository,
		accessTokenSecret:     accessTokenSecret,
		accessTokenExpiryHour: accessTokenExpiryHour,
	}
}

func (u *AccountUseCase) Signin(request domain.SigninRequest) (*domain.SigninResponse, error) {
	return nil, nil
}

func (u *AccountUseCase) Signup(request domain.SignupRequest) (*domain.SignupResponse, error) {
	return nil, nil
}
