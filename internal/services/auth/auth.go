package auth

import (
	"context"
	"errors"
	"log"

	pg "github.com/mnocard/go-project/internal/storage"
)

type UserStorage interface {
	FindUserByName(context.Context, string) (*pg.User, error)
}

type authService struct {
	uStorage UserStorage
}

func New(s UserStorage) *authService {
	return &authService{uStorage: s}
}

func (uService *authService) Auth(ctx context.Context, uName, password string) (bool, error) {
	user, err := uService.uStorage.FindUserByName(ctx, uName)
	if err != nil {
		log.Println("authService) Auth err 1", err)
		return false, err
	}

	if user == nil || user.Password != password {
		log.Println("authService) Auth err 2", user, password)
		return false, errors.New("user not found or wrong password")
	}

	return user.IsAdmin, nil
}
