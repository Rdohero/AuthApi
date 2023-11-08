package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func Login(c *gin.Context) {
	var body struct {
		Fullname  string
		Username  string
		Email     string
		Password  string
		Emailuser string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	checkEmailD := checkEmailDomain(body.Emailuser)
	if checkEmailD != nil {
		body.Username = body.Emailuser
	} else {
		body.Email = body.Emailuser
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ? OR username = ?", body.Email, body.Username)

	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if user.Active == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Email has not been verified",
		})
	} else {
		if errPassword != nil && user.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Email atau Username atau Password Salah",
			})
		} else {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": user.ID,
				"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
			})

			tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
			c.JSON(http.StatusOK, gin.H{
				"Token": tokenString,
			})
		}
	}
}
