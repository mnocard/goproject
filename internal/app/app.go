package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handlers "github.com/mnocard/go-project/internal/handlers"
	uService "github.com/mnocard/go-project/internal/services/user"
	uStorage "github.com/mnocard/go-project/internal/storage"
)

func RunRouter() error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))

	uStorage, err := uStorage.NewStorage()
	if err != nil {
		panic(err)
	}

	uService := uService.New(uStorage)
	h := handlers.New(uService)
	authorized.POST("admin", h.ChangeAdminPassword)

	return r.Run()
}
