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

func AuthUser(a Auth) gin.HandlerFunc {
	return func(c *gin.Context) {

		user, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"message": "authorization error"})
			return
		}

		if ok, err := a.Auth(c.Request.Context(), user, password); !ok {
			if err != nil {
				log.Println(err)
			}
			c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"message": "authorization error"})
			return
		}

		c.Set(gin.AuthUserKey, user)
		c.Next()
	}
}
