package user

import (
	"log"

	"github.com/mnocard/go-project/internal/services/config"
)

var users = []User{
	{
		UserName: "test",
		Password: "test1",
	},
	{
		UserName: "admin",
		Password: "admin",
	},
	{
		UserName: "test2",
		Password: "test2",
	},
}

type notFoundError struct {
}

func (e *notFoundError) Error() string {
	return "not found"
}

type User struct {
	Id       int
	UserName string
	Password string
}

type storage struct {
	connectionString string
}

func NewStorage() (*storage, error) {
	conn, err := config.GetConnectionString()
	if err != nil {
		return nil, err
	}
	return &storage{connectionString: conn}, nil
}

func (s *storage) Create(u User) (int, error) {
	user, err := s.FindByName(u.UserName)
	if err != nil {
		users = append(users, u)
		return len(users), nil
	}

	return user.Id, nil
}

func (s *storage) FindByName(uName string) (User, error) {
	for _, v := range users {
		if v.UserName == uName {
			return v, nil
		}
	}
	return User{}, &notFoundError{}
}

func (s *storage) Get(id int) (User, error) {
	if id >= len(users) {
		return User{}, &notFoundError{}
	}
	return users[id], nil
}

func (s *storage) Update(u User) (User, error) {
	log.Println("Update, ", u)
	for i := range users {
		if i == u.Id {
			users[i].UserName = u.UserName
			users[i].Password = u.Password
			return users[i], nil
		}
	}

	return User{}, &notFoundError{}
}

func (s *storage) Delete(id int) (bool, error) {
	user, err := s.Get(id)
	if err != nil {
		return false, &notFoundError{}
	}

	users = append(users[:user.Id], users[user.Id+1:]...)
	return true, nil
}
