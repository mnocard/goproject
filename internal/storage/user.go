package user

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mnocard/go-project/internal/services/config"
)

type notFoundError struct {
}

func (e *notFoundError) Error() string {
	return "not found"
}

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type storage struct {
	pool *pgxpool.Pool
}

func NewStorage() (*storage, error) {
	ctx := context.Background()

	connString, err := config.GetConnectionString()
	if err != nil {
		log.Println("connString error", err)
		return nil, err
	}

	log.Println("connString ", connString)

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}
	log.Println("db connected successfully")
	storage := storage{pool: dbpool}
	storage.checkTable(ctx)
	log.Println("db users exists")

	return &storage, nil
}

func (s *storage) checkTable(ctx context.Context) {
	sql := `
	    CREATE TABLE IF NOT EXISTS users
        (
            id          bigserial    NOT NULL PRIMARY KEY,
            username    text         NOT NULL,
            password    text         NOT NULL,
			is_admin    bool         NOT NULL,
			UNIQUE(username)
        );`
	_, err := s.pool.Exec(ctx, sql)
	if err != nil {
		s.logError(err, "storage) checkTable - create table")
	}

	user, err := s.FindByName(ctx, "admin")
	if err != nil && err != pgx.ErrNoRows {
		s.logError(err, "storage) checkTable - find admin")
		panic(err)
	}

	if user != nil {
		log.Println("storage) checkTable - admin is already exist")
		return
	}

	sql = `
        INSERT INTO users (username, password, is_admin)
        VALUES ($1, $2, true)
        RETURNING id;
	`

	commandTag, err := s.pool.Exec(ctx, sql, "admin", "admin")
	if err != nil {
		s.logError(err, "storage) checkTable - create admin")
	}

	if commandTag.RowsAffected() == 0 {
		log.Println("storage) checkTable - create admin ")
	}
}

func (s *storage) Close() {
	s.pool.Close()
}

func (s *storage) logError(err error, caller string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		log.Println(caller, pgErr.Message)
		log.Println(caller, pgErr.Code)
		return
	}
	log.Println(caller, err)
}

func (s *storage) Create(ctx context.Context, u *User) (int, error) {
	user, err := s.FindByName(ctx, u.UserName)
	if err == nil {
		return user.Id, nil
	}

	sql := `
        INSERT INTO users (username, password, is_admin)
        VALUES ($1, $2, $3)
        RETURNING id;
	`

	var id int
	err = s.pool.QueryRow(ctx, sql, u.UserName, u.Password, u.IsAdmin).Scan(&id)
	if err != nil {
		s.logError(err, "storage) Create")
		return 0, err
	}

	return user.Id, nil
}

func (s *storage) FindByName(ctx context.Context, uName string) (*User, error) {
	var user User

	sql := `
    	SELECT id, username, password, is_admin FROM users
    	WHERE username = $1;
	`

	err := s.pool.QueryRow(ctx, sql, uName).Scan(&user.Id, &user.UserName, &user.Password, &user.IsAdmin)
	if err != nil {
		s.logError(err, "storage) FindByName")
		return nil, err
	}

	log.Println("storage) FindByName user:", user)
	return &user, nil
}

func (s *storage) FindById(ctx context.Context, id int) (*User, error) {
	var user User

	sql := `
    	SELECT id, username, password, is_admin FROM user-db
    	WHERE id = $1;
	`

	err := s.pool.QueryRow(ctx, sql, id).Scan(&user.Id, &user.UserName, &user.Password, &user.IsAdmin)
	if err != nil {
		s.logError(err, "storage) FindById")
		return nil, err
	}

	log.Println("storage) FindById user:", user)
	return &user, nil
}

func (s *storage) Update(ctx context.Context, u User) (*User, error) {
	var user *User

	if u.Id != 0 {
		user, _ = s.FindById(ctx, u.Id)
	} else {
		user, _ = s.FindByName(ctx, u.UserName)
	}

	if user.Id == 0 {
		return nil, &notFoundError{}
	}

	sql := `
    	UPDATE users
		SET username = $2, password = $3, is_admin = $4
    	WHERE id = $1
		RETURNING id;
	`

	commandTag, err := s.pool.Exec(ctx, sql, user.Id, u.UserName, u.Password, u.IsAdmin)
	if err != nil {
		s.logError(err, "storage) Update error")
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		err = &notFoundError{}
		s.logError(err, "storage) Update not found")
		return nil, err
	}

	log.Println("storage) Update user:", user)
	return user, nil
}

func (s *storage) Delete(ctx context.Context, id int) (bool, error) {
	sql := `DELETE FROM users WHERE id = $1`

	commandTag, err := s.pool.Exec(ctx, sql, id)
	if err != nil {
		s.logError(err, "storage) Delete error")
		return false, err
	}

	if commandTag.RowsAffected() == 0 {
		err = &notFoundError{}
		s.logError(err, "storage) Delete not found")
		return false, err
	}

	return true, nil
}
