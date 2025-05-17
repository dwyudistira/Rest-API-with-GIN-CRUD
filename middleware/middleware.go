package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	SECRET = "secret"
)

func AuthValidate(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Token required",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Cek apakah metode signing adalah HMAC
		if _, invalid := token.Method.(*jwt.SigningMethodHMAC); !invalid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not authorized",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	fmt.Println("Token verified")
	c.Next()
}
