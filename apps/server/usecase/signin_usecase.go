package usecase

import (
	"strings"

	"golang.org/x/crypto/bcrypt"

	"omni-pixel/domain"
	"omni-pixel/internal"
)

type SigninUsecase struct {
	userRepository        domain.UserRepository
	accessTokenSecret     string
	accessTokenExpiryHour int
}

func NewSigninUsecase(userRepository domain.UserRepository, accessTokenSecret string, accessTokenExpiryHour int) *SigninUsecase {
	return &SigninUsecase{
		userRepository:        userRepository,
		accessTokenSecret:     accessTokenSecret,
		accessTokenExpiryHour: accessTokenExpiryHour,
	}
}

func (u *SigninUsecase) Signin(request domain.SigninRequest) (*domain.SigninResponse, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	if email == "" || request.Password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := u.userRepository.FindByEmail(email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, domain.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := internal.CreateAccessToken(user, u.accessTokenSecret, u.accessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	if err := u.userRepository.UpdateLastLogin(user.ID); err != nil {
		return nil, err
	}

	user.PasswordHash = ""

	return &domain.SigninResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   u.accessTokenExpiryHour * 60 * 60,
		User:        *user,
	}, nil
}
