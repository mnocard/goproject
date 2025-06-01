package config

import (
	"errors"
	"os"
)

func GetConnectionString() (string, error) {
	conn := os.Getenv("CONN_STRING")
	if conn == "" {
		return "", errors.New("connection string not found")
	}
	return conn, nil
}
