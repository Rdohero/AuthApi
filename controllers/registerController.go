package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

func Signup(c *gin.Context) {
	// Get the email/pass off req body
	var body struct {
		Fullname string
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	// check password criteria
	checkPassword := checkPasswordCriteria(body.Password)
	if checkPassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": checkPassword.Error(),
		})
	}

	// check username criteria
	checkUsername := checkUsernameCriteria(body.Username)
	if checkUsername != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": checkUsername.Error(),
		})
	}

	// check email is valid
	checkEmail := checkEmailValid(body.Email)
	if checkEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": checkEmail.Error(),
		})
	}

	//check if email domain exists
	checkEmailD := checkEmailDomain(body.Email)
	if checkEmailD != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": checkEmailD.Error(),
		})
	}
	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash password",
		})

		return
	}

	if checkUsername == nil && checkPassword == nil && checkEmail == nil && checkEmailD == nil {
		otp := rand.Intn(900000) + 100000

		otpStr := fmt.Sprintf("%06d", otp)

		token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": body.Email,
			"otp":   otp,
			"exp":   time.Now().Add(time.Minute * 2).Unix(),
		}).SignedString([]byte(os.Getenv("SECRET")))

		SimpanOtp(otpStr, token)
		// Create the user
		user := models.User{Fullname: body.Fullname, Username: body.Username, Email: body.Email, Password: string(hash), Active: false}
		result := initializers.DB.Create(&user)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Email is already in use",
			})

			return
		}

		subject := "Email Verificaion"
		HTMLbody :=
			`<html>
			<h1>Code to Verify Email</h1>
			<p>` + otpStr + `</p>
		</html>`
		to := []string{body.Email}
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
		fmt.Println("Check for sent email!")

		// Respond
		c.JSON(http.StatusOK, gin.H{
			"Status": "Succes, Check for sent email!",
		})
	}
}
