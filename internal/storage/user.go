package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Rating   int    `json:"rating"`
	IsAdmin  bool   `json:"is_admin"`
}

func (s *storage) CreateUser(ctx context.Context, u *User) (int, error) {
	user, err := s.FindUserByName(ctx, u.UserName)
	if err == nil {
		return user.Id, nil
	}

	sql := `
        INSERT INTO users (username, password, rating, is_admin)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
	`

	var id int
	err = s.pool.QueryRow(ctx, sql, u.UserName, u.Password, u.Rating, u.IsAdmin).Scan(&id)
	if err != nil {
		s.logError(err, "storage) CreateUser")
		return 0, err
	}

	return id, nil
}

func (s *storage) FindUserByName(ctx context.Context, uName string) (*User, error) {
	var user User

	sql := `
    	SELECT id, username, password, rating, is_admin FROM users
    	WHERE username = $1;
	`

	err := s.pool.QueryRow(ctx, sql, uName).Scan(&user.Id, &user.UserName, &user.Password, &user.Rating, &user.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		s.logError(err, "storage) FindUserByName error")
		return nil, err
	}

	log.Println("storage) FindUserByName user:", user)
	return &user, nil
}

func (s *storage) FindUserById(ctx context.Context, id int) (*User, error) {
	var user User

	sql := `
    	SELECT id, username, password, rating, is_admin FROM users
    	WHERE id = $1;
	`

	err := s.pool.QueryRow(ctx, sql, id).Scan(&user.Id, &user.UserName, &user.Password, &user.Rating, &user.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		s.logError(err, "storage) FindUserById")
		return nil, err
	}

	log.Println("storage) FindUserById user:", user)
	return &user, nil
}

func (s *storage) UpdateUser(ctx context.Context, u User) (*User, error) {
	var user *User

	if u.Id != 0 {
		user, _ = s.FindUserById(ctx, u.Id)
	} else {
		user, _ = s.FindUserByName(ctx, u.UserName)
	}

	if user.Id == 0 {
		return nil, &NotFoundError{}
	}

	sql := `
    	UPDATE users
		SET username = $2, password = $3, rating = $4, is_admin = $5
    	WHERE id = $1
		RETURNING id;
	`

	commandTag, err := s.pool.Exec(ctx, sql, user.Id, u.UserName, u.Password, u.Rating, u.IsAdmin)
	if err != nil {
		s.logError(err, "storage) UpdateUser error")
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		err = &NotFoundError{}
		s.logError(err, "storage) UpdateUser not found")
		return nil, err
	}

	log.Println("storage) UpdateUser user:", user)
	return user, nil
}

func (s *storage) DeleteUser(ctx context.Context, id int) (bool, error) {
	sql := `DELETE FROM users WHERE id = $1`

	commandTag, err := s.pool.Exec(ctx, sql, id)
	if err != nil {
		s.logError(err, "storage) DeleteUser error")
		return false, err
	}

	if commandTag.RowsAffected() == 0 {
		err = &NotFoundError{}
		s.logError(err, "storage) DeleteUser not found")
		return false, err
	}

	return true, nil
}
