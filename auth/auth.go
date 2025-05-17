package auth

import (
	"RestApi/model"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
)

const (
	USER     = "admin"
	PASSWORD = "pass1234"
	SECRET   = "secret"
)

func LoginHandler(c *gin.Context) {
	var user model.Creadential

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	if user.Username != USER {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User Invalid",
		})
		return
	}else if user.Password != PASSWORD {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Password Invalid",
		})
		return
	} else{

		// Token
		claims := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Issuer:    "Test",
			IssuedAt:  time.Now().Unix(),
		}
	
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(SECRET))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"token":   signedToken,
		})
	}
}
