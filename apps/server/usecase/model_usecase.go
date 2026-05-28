package usecase

import "omni-pixel/domain"

type ModelUseCase struct {
	modelRepository domain.ModelRepository
}

func NewModelUseCase(modelRepository domain.ModelRepository) *ModelUseCase {
	return &ModelUseCase{modelRepository}
}

func (u *ModelUseCase) ListModels(filter domain.ModelFilter) ([]domain.ModelResponse, error) {
	models, err := u.modelRepository.FindAll(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]domain.ModelResponse, 0, len(models))
	for _, m := range models {
		providerName := ""
		if m.Provider != nil {
			providerName = m.Provider.Name
		}

		responses = append(responses, domain.ModelResponse{
			ID:           m.ID,
			ProviderID:   m.ProviderID,
			ProviderName: providerName,
			ModelName:    m.ModelName,
			IsEnabled:    m.IsEnabled,
			ExpireTime:   m.ExpireTime,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		})
	}

	return responses, nil
}
