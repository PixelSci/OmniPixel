package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"omni-pixel/apperr"
	"omni-pixel/model"
	"omni-pixel/repository"
)

// Service orchestrates identity operations. It composes a set of provider
// plugins: password today, OAuth providers (github, google) will land in the
// oauth map as they're added. Handlers depend only on *Service.
type Service struct {
	users    repository.UserRepository
	tokens   *TokenIssuer
	password PasswordAuthenticator
	oauth    map[string]OAuthAuthenticator // keyed by Provider.Name(); populated as providers are wired
}

func NewService(users repository.UserRepository, tokens *TokenIssuer, password PasswordAuthenticator) *Service {
	return &Service{
		users:    users,
		tokens:   tokens,
		password: password,
		oauth:    make(map[string]OAuthAuthenticator),
	}
}

// RegisterOAuth plugs in a new OAuth provider. Call this from main after
// constructing Service, e.g.  svc.RegisterOAuth(github.New(cfg)).
func (s *Service) RegisterOAuth(p OAuthAuthenticator) {
	s.oauth[p.Name()] = p
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type AuthResult struct {
	User  *model.User
	Token *TokenPair
}

// Register creates a new password-authenticated user and returns an access token.
func (s *Service) Register(ctx context.Context, in RegisterInput) (*AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	hash, err := hashPassword(in.Password)
	if err != nil {
		return nil, apperr.ErrInternal.WithCause(err)
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:           uuid.NewString(),
		Name:         strings.TrimSpace(in.Name),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, apperr.ErrUserExists.WithDetail("email", email).WithCause(err)
		}
		return nil, apperr.ErrInternal.WithCause(err)
	}
	return s.issue(user)
}

// Login authenticates an email/password pair via the password provider.
func (s *Service) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := s.password.Authenticate(ctx, strings.ToLower(strings.TrimSpace(email)), password)
	if err != nil {
		return nil, err
	}
	return s.issue(user)
}

// ---------------------------------------------------------------------------
// Future OAuth entrypoints (github, google). Uncomment and wire once the
// corresponding providers are implemented and registered via RegisterOAuth.
//
// func (s *Service) OAuthURL(providerName, state string) (string, error) {
// 	p, ok := s.oauth[providerName]
// 	if !ok {
// 		return "", apperr.ErrInvalidParam.WithDetail("provider", providerName)
// 	}
// 	return p.AuthCodeURL(state), nil
// }
//
// func (s *Service) OAuthCallback(ctx context.Context, providerName, code string) (*AuthResult, error) {
// 	p, ok := s.oauth[providerName]
// 	if !ok {
// 		return nil, apperr.ErrInvalidParam.WithDetail("provider", providerName)
// 	}
// 	user, err := p.Exchange(ctx, code)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return s.issue(user)
// }
// ---------------------------------------------------------------------------

func (s *Service) issue(user *model.User) (*AuthResult, error) {
	token, err := s.tokens.Issue(user)
	if err != nil {
		return nil, apperr.ErrInternal.WithCause(err)
	}
	return &AuthResult{User: user, Token: token}, nil
}
