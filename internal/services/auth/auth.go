package auth

import (
	"context"
	"errors"

	pg "github.com/mnocard/go-project/internal/storage"
)

type UserStorage interface {
	FindByName(context.Context, string) (*pg.User, error)
}

type authService struct {
	uStorage UserStorage
}

func New(s UserStorage) *authService {
	return &authService{uStorage: s}
}

func (uService *authService) Auth(ctx context.Context, uName, password string) (bool, error) {
	user, err := uService.uStorage.FindByName(ctx, uName)
	if err != nil {
		return false, err
	}

	if user == nil || user.Password != password {
		return false, errors.New("user not found or wrong password")
	}

	return user.IsAdmin, nil
}
