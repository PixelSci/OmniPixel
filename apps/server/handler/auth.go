package handler

import (
	"github.com/gofiber/fiber/v2"

	"omni-pixel/apperr"
	"omni-pixel/auth"
)

type AuthHandler struct {
	svc *auth.Service
}

func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// ---- DTOs ----------------------------------------------------------------

type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=64"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token TokenDTO     `json:"token"`
}

func toAuthResponse(r *auth.AuthResult) AuthResponse {
	return AuthResponse{
		User: toUserResponse(r.User),
		Token: TokenDTO{
			AccessToken: r.Token.AccessToken,
			TokenType:   r.Token.TokenType,
			ExpiresIn:   r.Token.ExpiresIn,
		},
	}
}

// ---- Handlers ------------------------------------------------------------

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.ErrInvalidParam.WithCause(err)
	}
	if err := validateStruct(&req); err != nil {
		return err
	}
	res, err := h.svc.Register(c.UserContext(), auth.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return Created(c, toAuthResponse(res))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.ErrInvalidParam.WithCause(err)
	}
	if err := validateStruct(&req); err != nil {
		return err
	}
	res, err := h.svc.Login(c.UserContext(), req.Email, req.Password)
	if err != nil {
		return err
	}
	return Success(c, toAuthResponse(res))
}

// ---- Future OAuth endpoints (github / google) ----------------------------
//
// func (h *AuthHandler) OAuthRedirect(c *fiber.Ctx) error {
// 	provider := c.Params("provider")
// 	state := c.Query("state")
// 	url, err := h.svc.OAuthURL(provider, state)
// 	if err != nil {
// 		return err
// 	}
// 	return c.Redirect(url, fiber.StatusTemporaryRedirect)
// }
//
// func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
// 	provider := c.Params("provider")
// 	code := c.Query("code")
// 	res, err := h.svc.OAuthCallback(c.UserContext(), provider, code)
// 	if err != nil {
// 		return err
// 	}
// 	return Success(c, toAuthResponse(res))
// }
