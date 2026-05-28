package usecase

import "omni-pixel/domain"

type ProviderUseCase struct {
	providerRepository domain.ProviderRepository
}

func NewProviderUseCase(providerRepository domain.ProviderRepository) *ProviderUseCase {
	return &ProviderUseCase{providerRepository}
}

func (u *ProviderUseCase) ListProviders() ([]domain.ProviderResponse, error) {
	providers, err := u.providerRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]domain.ProviderResponse, 0, len(providers))
	for _, p := range providers {
		responses = append(responses, domain.ProviderResponse{
			ID:        p.ID,
			Name:      p.Name,
			BaseURL:   p.BaseURL,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return responses, nil
}
