package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	Auth(ctx context.Context, uName, password string) (bool, error)
}

const isAdminKey string = "isAdmin"

func AuthUser(a Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("AuthUser(a Auth) gin.HandlerFunc c.Request", c.Request)

		user, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(http.StatusNetworkAuthenticationRequired, gin.H{"message": "authorization error 1"})
			return
		}

		log.Println("AuthUser(a Auth) gin.HandlerFunc", user, password)

		isAdmin, err := a.Auth(c.Request.Context(), user, password)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusNetworkAuthenticationRequired, gin.H{"message": "authorization error 2"})
			return
		}

		c.Set(isAdminKey, isAdmin)
		c.Set(gin.AuthUserKey, user)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := c.GetBool(isAdminKey)
		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"message": "you need to be admin 3"})
			return
		}

		c.Next()
	}
}
