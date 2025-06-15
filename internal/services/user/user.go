package user

import (
	"context"

	pg "github.com/mnocard/go-project/internal/storage"
)

type UserStorage interface {
	CreateUser(context.Context, *pg.User) (int, error)
	FindUserByName(context.Context, string) (*pg.User, error)
	FindUserById(context.Context, int) (*pg.User, error)
	UpdateUser(context.Context, pg.User) (*pg.User, error)
	DeleteUser(context.Context, int) (bool, error)
}

type User struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Rating   int    `json:"rating"`
	IsAdmin  bool   `json:"is_admin" `
}

type userService struct {
	uStorage UserStorage
}

func New(s UserStorage) *userService {
	return &userService{uStorage: s}
}

func (uService *userService) Create(ctx context.Context, u *User) (int, error) {
	return uService.uStorage.CreateUser(ctx, &pg.User{
		UserName: u.UserName,
		Password: u.Password,
		Rating:   u.Rating,
		IsAdmin:  u.IsAdmin,
	})
}

func (uService *userService) Get(ctx context.Context, id int) (*User, error) {
	user, err := uService.uStorage.FindUserById(ctx, id)
	return &User{
		UserName: user.UserName,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}, err
}

func (uService *userService) FindByName(ctx context.Context, uName string) (*User, error) {
	user, err := uService.uStorage.FindUserByName(ctx, uName)
	return &User{
		UserName: user.UserName,
		Password: user.Password,
		IsAdmin:  user.IsAdmin,
	}, err
}

func (uService *userService) Update(ctx context.Context, uName, password string) (int, error) {
	user, err := uService.uStorage.FindUserByName(ctx, uName)
	if err != nil {
		return 0, err
	}

	user, err = uService.uStorage.UpdateUser(ctx, pg.User{
		UserName: uName,
		Password: password,
		IsAdmin:  user.IsAdmin,
	})

	return user.Id, err
}

func (uService *userService) Delete(ctx context.Context, uName string) (bool, error) {
	user, err := uService.uStorage.FindUserByName(ctx, uName)
	if err != nil {
		return false, nil
	}

	return uService.uStorage.DeleteUser(ctx, user.Id)
}
