package usecase

import (
	"omni-pixel/domain"
)

type HealthUseCase struct {
}

func NewHealthUseCase() *HealthUseCase {
	return &HealthUseCase{}
}

func (u *HealthUseCase) Health() (*domain.SigninResponse, error) {
	return nil, nil
}
