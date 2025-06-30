package handlers

import (
	"context"

	tService "github.com/mnocard/go-project/internal/services/task"
	uService "github.com/mnocard/go-project/internal/services/user"
)

type Handler struct {
	uService userService
	tService taskService
}

type userService interface {
	Create(context.Context, *uService.User) (int, error)
	Update(context.Context, string, string) (int, error)
	GetRating(context.Context, string) (int, error)
	FindByName(context.Context, string) (*uService.User, error)
	UpdateRating(ctx context.Context, id, rating int) error
	// Get(context.Context, int) (*uService.User, error)
	// Delete(context.Context, string) (bool, error)
}

type taskService interface {
	Create(context.Context, *tService.Task) (int, error)
	Update(context.Context, *tService.Task) (int, error)
	CompleteTask(context.Context, int, int) (int, error)
	// Get(context.Context, int) (*tService.Task, error)
	// FindTaskByUserId(context.Context, int) (*tService.Task, error)
	// Delete(context.Context, int) (bool, error)
}

func New(uService userService, tService taskService) *Handler {
	return &Handler{
		uService: uService,
		tService: tService,
	}
}

type NewPassword struct {
	Value string `json:"value" binding:"required"`
}
