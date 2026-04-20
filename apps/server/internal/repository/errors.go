package repository

import "errors"

// Repository sentinels are storage-layer signals. The service layer is
// responsible for translating them into apperr values — repositories must
// never depend on HTTP or apperr concerns.
var (
	ErrNotFound  = errors.New("repository: not found")
	ErrDuplicate = errors.New("repository: duplicate")
)
