package usecase

import (
	"omni-pixel/domain"
	"omni-pixel/internal"
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
	user, err := u.userRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if !internal.CheckPassword(request.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := internal.CreateAccessToken(user, u.accessTokenSecret, u.accessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	if err := u.userRepository.UpdateLastLogin(user.ID); err != nil {
		return nil, err
	}

	return &domain.SigninResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   u.accessTokenExpiryHour * 3600,
		User:        *user,
	}, nil
}

func (u *AccountUseCase) Signup(request domain.SignupRequest) (*domain.SignupResponse, error) {
	if err := u.userRepository.ExistsByEmail(request.Email); err != nil {
		return nil, err
	}

	passwordHash, err := internal.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: passwordHash,
		Status:       "active",
	}

	if err := u.userRepository.Create(user); err != nil {
		return nil, err
	}

	accessToken, err := internal.CreateAccessToken(user, u.accessTokenSecret, u.accessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	return &domain.SignupResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   u.accessTokenExpiryHour * 3600,
		User:        *user,
	}, nil
}
