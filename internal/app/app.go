package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var users = map[string]string{
	"admin": "admin",
}

func GetRouter() error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))

	// r.POST("/admin", func(c *gin.Context) {
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		log.Println("admin start")
		// user := "admin"
		log.Println("request:", c.Request)
		log.Println("user", user)

		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			log.Println("c.Bind(&json) == nil")
			log.Println("json.Value", json.Value)

			users[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
		log.Println("admin end")
	})

	return r.Run()
}
