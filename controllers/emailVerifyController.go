package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

func ResendOtpEmailPassVer(c *gin.Context) {
	var Resend struct {
		Email string
	}

	c.Bind(&Resend)

	var user, _ = getUserByEmail(Resend.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "User not found",
		})
		return
	}

	otp := rand.Intn(900000) + 100000

	otpStr := fmt.Sprintf("%06d", otp)

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": Resend.Email,
		"otp":   otp,
		"exp":   time.Now().Add(time.Minute * 2).Unix(),
	}).SignedString([]byte(os.Getenv("SECRET")))

	SimpanOtp(otpStr, token)

	subject := "Email Verificaion"
	HTMLbody :=
		`<html>
			<h1>Code to Verify Email</h1>
			<p>` + otpStr + `</p>
		</html>`
	to := []string{Resend.Email}
	// SMTP - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// Set up authentication information
	auth := smtp.PlainAuth("", "crafterman79@gmail.com", "njij loxz hjry lacv", host)
	// Construct the message
	msg := []byte(
		"Subject: " + subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-s8\"\r\n" +
			"\r\n" +
			HTMLbody)
	err := smtp.SendMail(address, auth, "crafterman79@gmail.com", to, msg)

	if err != nil {
		fmt.Println("Error sending email")
	}

	c.JSON(http.StatusOK, gin.H{
		"Status": "Resend Code Succes",
	})
}

func OtpEmailVer(c *gin.Context) {
	var Otp struct {
		Email string
		Otp   string
	}
	c.Bind(&Otp)

	var token2, err1 = DapatkanOtpString(Otp.Otp)

	if token2 == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err1.Error(),
		})
		return
	}

	token, _ := jwt.Parse(token2, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			HapusOtp(Otp.Otp)
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Otp Has Been Expired",
			})
			return
		}

		user, err := getUserByEmail(claims["email"].(string))

		var userUpdate []models.User

		if Otp.Email == user.Email {
			if Otp.Otp == strconv.Itoa(int(claims["otp"].(float64))) {
				// update user.Active to true
				err = MakeActive(user.ID)

				HapusOtp(Otp.Otp)

				initializers.DB.First(&userUpdate, user.ID).Update("token_string", nil)
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
				"Error": "Otp Not Valid",
			})
		}
	} else {
		HapusOtp(Otp.Otp)
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Otp Has Been Expired",
		})
		return
	}
}

func getUserByEmail(email string) (*models.User, error) {
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
