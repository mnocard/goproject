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
	tService "github.com/mnocard/go-project/internal/services/task"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Use(auth.AuthUser(aService))

	uService := uService.New(uStorage)
	tService := tService.New(uStorage)
	h := handlers.New(uService, tService)

	admin := r.Group("/admin")
	admin.Use(auth.AdminRequired())
	{
		admin.POST("/changeAdminPassword", h.ChangeAdminPassword)
		admin.POST("/createUser", h.CreateUser)
		admin.POST("/createTask", h.CreateTask)
	}

	user := r.Group("/user")
	{
		user.GET("/getRating", h.GetRating)
		user.GET("/completeTask/:taskId", h.CompleteTask)
	}

	return r.Run()
}
