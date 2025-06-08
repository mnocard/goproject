package user

import (
	"context"

	pg "github.com/mnocard/go-project/internal/storage"
)

type UserStorage interface {
	Create(context.Context, *pg.User) (int, error)
	FindByName(context.Context, string) (*pg.User, error)
	FindById(context.Context, int) (*pg.User, error)
	Update(context.Context, pg.User) (*pg.User, error)
	Delete(context.Context, int) (bool, error)
}

type User struct {
	UserName string
	Password string
	IsAdmin  bool
}

type userService struct {
	uStorage UserStorage
}

func New(s UserStorage) *userService {
	return &userService{uStorage: s}
}

func (uService *userService) Create(ctx context.Context, uName, password string, isAdmin bool) (int, error) {
	return uService.uStorage.Create(ctx, &pg.User{
		UserName: uName,
		Password: password,
		IsAdmin:  isAdmin,
	})
}

func (uService *userService) Get(ctx context.Context, id int) (*User, error) {
	user, err := uService.uStorage.FindById(ctx, id)
	return &User{
		UserName: user.UserName,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}, err
}

func (uService *userService) FindByName(ctx context.Context, uName string) (*User, error) {
	user, err := uService.uStorage.FindByName(ctx, uName)
	return &User{
		UserName: user.UserName,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}, err
}

func (uService *userService) Update(ctx context.Context, uName, password string) (int, error) {
	user, err := uService.uStorage.FindByName(ctx, uName)
	if err != nil {
		return 0, err
	}

	user, err = uService.uStorage.Update(ctx, pg.User{
		UserName: uName,
		Password: password,
		IsAdmin:  user.IsAdmin,
	})

	return user.Id, err
}

func (uService *userService) Delete(ctx context.Context, uName string) (bool, error) {
	user, err := uService.uStorage.FindByName(ctx, uName)
	if err != nil {
		return false, nil
	}

	return uService.uStorage.Delete(ctx, user.Id)
}
