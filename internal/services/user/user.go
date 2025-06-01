package user

import pg "github.com/mnocard/go-project/internal/storage"

type UserStorage interface {
	Create(u pg.User) (int, error)
	FindByName(uName string) (pg.User, error)
	Get(id int) (pg.User, error)
	Update(u pg.User) (pg.User, error)
	Delete(id int) (bool, error)
}

type User struct {
	UserName string
	Password string
}

type userService struct {
	uStorage UserStorage
}

func New(s UserStorage) *userService {
	return &userService{uStorage: s}
}

func (uService *userService) Create(uName, password string) (int, error) {
	return uService.uStorage.Create(pg.User{
		UserName: uName,
		Password: password,
	})
}

func (uService *userService) Get(id int) (*User, error) {
	user, err := uService.uStorage.Get(id)
	return &User{
		UserName: user.UserName,
		Password: user.Password,
	}, err
}

func (uService *userService) FindByName(uName string) (*User, error) {
	user, err := uService.uStorage.FindByName(uName)
	return &User{
		UserName: user.UserName,
		Password: user.Password,
	}, err
}

func (uService *userService) Update(uName, password string) (int, error) {
	user, err := uService.uStorage.FindByName(uName)
	if err != nil {
		return 0, err
	}

	user, err = uService.uStorage.Update(pg.User{
		UserName: uName,
		Password: password,
	})

	return user.Id, err
}

func (uService *userService) Delete(uName string) (bool, error) {
	user, err := uService.uStorage.FindByName(uName)
	if err != nil {
		return false, nil
	}

	return uService.uStorage.Delete(user.Id)
}
