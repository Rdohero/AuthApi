package middleware

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func RequiredAuth(c *gin.Context) {
	//Get the cookie off req
	var tokenBody struct {
		Token string
	}

	if c.Bind(&tokenBody) != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := tokenBody.Token

	// Decode/validate it
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// Find the user with token sub
		var user []models.User
		initializers.DB.First(&user, claims["sub"])

		if user[0].ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		userMap := map[string]interface{}{
			"ID": user[0].ID,
			// Data pengguna lainnya, jika diperlukan
		}

		c.Set("userMap", userMap)
		// Attack to req
		c.Set("user", user)

		// Continue
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
