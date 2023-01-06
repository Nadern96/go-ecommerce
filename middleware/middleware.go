package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadern96/go-ecommerce/helpers"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		log.Println("token = ", clientToken)
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
