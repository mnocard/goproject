package storage

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

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return "not found"
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
	err := s.createUserTable(ctx)
	if err != nil {
		s.logError(err, "storage) checkTable - create user table")
	}

	err = s.createTaskTable(ctx)
	if err != nil {
		s.logError(err, "storage) checkTable - create task table")
	}

	user, err := s.FindUserByName(ctx, "admin")
	if err != nil && err != pgx.ErrNoRows {
		s.logError(err, "storage) checkTable - find admin")
		panic(err)
	}

	if user != nil {
		log.Println("storage) checkTable - admin is already exist")
		return
	}

	sql := `
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

func (s *storage) createUserTable(ctx context.Context) error {
	sql := `
	    CREATE TABLE IF NOT EXISTS users
        (
            id          serial    NOT NULL PRIMARY KEY,
            username    text      NOT NULL,
            password    text      NOT NULL,
			is_admin    bool      NOT NULL,
			UNIQUE(username)
        );`
	_, err := s.pool.Exec(ctx, sql)
	return err
}

func (s *storage) createTaskTable(ctx context.Context) error {
	sql := `
	    CREATE TABLE IF NOT EXISTS tasks
        (
            id                serial     NOT NULL PRIMARY KEY,
            user_id           integer    NOT NULL,
            points            integer    NOT NULL,
			parent_task_id    integer    NOT NULL,
        );`
	_, err := s.pool.Exec(ctx, sql)
	return err
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
