package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"omni-pixel/apperr"
	"omni-pixel/model"
	"omni-pixel/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// ---- DTOs ----------------------------------------------------------------

type CreateUserRequest struct {
	Name  string `json:"name"  validate:"required,min=2,max=64"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Items    []UserResponse `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

func toUserResponse(u *model.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ---- Handlers ------------------------------------------------------------

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.ErrInvalidParam.WithCause(err)
	}
	if err := validateStruct(&req); err != nil {
		return err
	}

	u, err := h.svc.Create(c.UserContext(), service.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return err
	}
	return Created(c, toUserResponse(u))
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return apperr.ErrInvalidParam.WithDetail("id", "required")
	}
	u, err := h.svc.Get(c.UserContext(), id)
	if err != nil {
		return err
	}
	return Success(c, toUserResponse(u))
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("page_size", 20)

	users, total, err := h.svc.List(c.UserContext(), page, size)
	if err != nil {
		return err
	}
	items := make([]UserResponse, 0, len(users))
	for _, u := range users {
		items = append(items, toUserResponse(u))
	}
	return Success(c, ListUsersResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: size,
	})
}
