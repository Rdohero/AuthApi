package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

func OtpEmailVer(c *gin.Context) {
	var Otp struct {
		Otp string
	}
	c.Bind(&Otp)
	var otpStore = TokenString

	token, _ := jwt.Parse(otpStore, func(token *jwt.Token) (interface{}, error) {
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

		user, err := getUserByUsername(claims["email"].(string))
		// Find the user with token sub
		if Otp.Otp == strconv.Itoa(int(claims["otp"].(float64))) {
			// update user.Active to true
			err = MakeActive(user.ID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"Error": "Please try verification email again",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"Status": "Succes",
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Otp Not Valid",
			})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Please try link in verification email again",
		})
	}
}

func getUserByUsername(email string) (*models.User, error) {
	var u models.User
	result := initializers.DB.Where("email = ?", email).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

func MakeActive(userID uint) error {
	var u models.User
	if err := initializers.DB.First(&u, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("User not found")
		}
		return err
	}

	u.Active = true
	if err := initializers.DB.Save(&u).Error; err != nil {
		return err
	}
	return nil
}
