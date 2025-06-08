package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mnocard/go-project/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	handlers "github.com/mnocard/go-project/internal/handlers"
	auth "github.com/mnocard/go-project/internal/middleware/auth"
	aService "github.com/mnocard/go-project/internal/services/auth"
	uService "github.com/mnocard/go-project/internal/services/user"
	uStorage "github.com/mnocard/go-project/internal/storage"
)

func RunRouter() error {
	r := gin.Default()

	uStorage, err := uStorage.NewStorage()
	if err != nil {
		panic(err)
	}
	defer uStorage.Close()

	aService := aService.New(uStorage)
	r.Use(auth.AuthUser(aService))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	uService := uService.New(uStorage)
	h := handlers.New(uService)
	r.POST("changeAdminPassword", h.ChangeAdminPassword)

	return r.Run()
}
