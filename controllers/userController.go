package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func checkUsernameCriteria(username string) error {
	// check username for only alphaNumeric characters
	var nameAlphaNumeric = true
	for _, char := range username {
		// func IsLetter(r rune) bool, func IsNumber(r rune) bool
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	if nameAlphaNumeric != true {
		// func New(text string) error
		return errors.New("Username must only contain letters and numbers")
	}
	// check username length
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	if nameLength != true {
		return errors.New("Username must be longer than 4 characters and less than 51")
	}
	return nil
}

func checkPasswordCriteria(password string) error {
	var err error
	// variables that must pass for password creation criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
			err = errors.New("Pa")
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true
		// func IsPunct(r rune) bool, func IsSymbol(r rune) bool
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	// check password length
	if 7 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	// create error for any criteria not passed
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces {
		switch false {
		case pswdLowercase:
			err = errors.New("Password must contain atleast one lower case letter")
		case pswdUppercase:
			err = errors.New("Password must contain atleast one uppercase letter")
		case pswdNumber:
			err = errors.New("Password must contain atleast one number")
		case pswdSpecial:
			err = errors.New("Password must contain atleast one special character")
		case pswdLength:
			err = errors.New("Passward length must atleast 12 characters and less than 60")
		case pswdNoSpaces:
			err = errors.New("Password cannot have any spaces")
		}
		return err
	}
	return nil
}

func checkEmailValid(email string) error {
	// check email syntax is valid
	//func MustCompile(str string) *Regexp
	emailRegex, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		fmt.Println(err)
		return errors.New("sorry, something went wrong")
	}
	rg := emailRegex.MatchString(email)
	if rg != true {
		return errors.New("Email address is not a valid syntax, please check again")
	}
	// check email length
	if len(email) < 4 {
		return errors.New("Email length is too short")
	}
	if len(email) > 253 {
		return errors.New("Email length is too long")
	}
	return nil
}

func checkEmailDomain(email string) error {
	i := strings.Index(email, "@")
	host := email[i+1:]
	// func LookupMX(name string) ([]*MX, error)
	_, err := net.LookupMX(host)
	if err != nil {
		err = errors.New("Could not find email's domain server, please check and try again")
		return err
	}
	return nil
}

func GetUserById(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User Tidak ditemukan",
		})
	}
	c.JSON(http.StatusOK, user)
}

func UpdateFullnameById(c *gin.Context) {
	id := c.Param("id")

	var user []models.User
	var inputUser models.User

	c.ShouldBindJSON(&inputUser)

	initializers.DB.Where("id = ?", id).Find(&user)

	initializers.DB.First(&user, id).Update("fullname", inputUser.Fullname)

	c.JSON(http.StatusOK, user)
}

func UpdateUsernameById(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	c.ShouldBindJSON(&user)

	username := user.Username
	checkUsername := checkUsernameCriteria(username)
	if checkUsername != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"StatusUsername": checkUsername.Error(),
		})
	} else {
		initializers.DB.Model(&user).Where("id = ?", id).Update("username", user.Username)
		c.JSON(http.StatusOK, user.Username)
	}
}

func UpdatePhotoProfile(c *gin.Context) {
	id := c.Param("id")

	var user1 []models.User
	initializers.DB.Where("id = ?", id).Find(&user1)

	oldfoto := user1[0].Foto
	os.Remove(oldfoto)

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allowedMIMETypes := []string{"image/jpeg", "image/png", "image/svg"}

	if !IsValidMIMEType(file, allowedMIMETypes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya menerima jpeg, png, dan svg"})
		return
	}
	// Define the path where the file will be saved
	basePath := filepath.Join("images", file.Filename)
	// Create the "images" directory if it doesn't exist
	os.MkdirAll("images", os.ModePerm)
	// Save the file to the defined path
	filePath := generateUniqueFileName(basePath)
	c.SaveUploadedFile(file, filePath)

	user1[0].Foto = filePath

	initializers.DB.Model(&user1).Where("id = ?", id).Update("foto", user1[0].Foto)

	c.JSON(http.StatusOK, user1)
}

func IsValidMIMEType(file *multipart.FileHeader, allowedMIMETypes []string) bool {
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// Membaca tipe MIME dari file
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return false
	}

	// Menggunakan http.DetectContentType untuk mendeteksi tipe MIME
	fileType := http.DetectContentType(buffer)

	// Memeriksa apakah tipe MIME ada dalam daftar yang diizinkan
	for _, allowedType := range allowedMIMETypes {
		if fileType == allowedType {
			return true
		}
	}

	return false
}

func ForgotPassword(c *gin.Context) {
	var ForgotPwd struct {
		Password string
		Email    string
		Otp      string
	}
	c.Bind(&ForgotPwd)

	checkPassword := checkPasswordCriteria(ForgotPwd.Password)
	if checkPassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": checkPassword.Error(),
		})
		return
	}

	var user models.User

	initializers.DB.Where("email = ?", ForgotPwd.Email).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "User not found",
		})
		return
	}

	var tokenString, err = DapatkanOtpString(ForgotPwd.Otp)

	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			HapusOtp(ForgotPwd.Otp)
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Otp Has Been Expired",
			})
			return
		}

		if ForgotPwd.Email == user.Email {
			if ForgotPwd.Otp == strconv.Itoa(int(claims["otp"].(float64))) {
				hash, _ := bcrypt.GenerateFromPassword([]byte(ForgotPwd.Password), 14)

				initializers.DB.First(&user).Update("password", string(hash))

				HapusOtp(ForgotPwd.Otp)

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
		HapusOtp(ForgotPwd.Otp)
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Otp Has Been Expired",
		})
		return
	}
}

func GetCartByUserId(c *gin.Context) {
	var id = c.Param("id")

	var cart []models.Cart
	initializers.DB.Preload("Product").Preload("Product.User").Where("user_id = ?", id).Find(&cart)

	c.JSON(http.StatusOK, cart)
}

func generateUniqueFileName(basePath string) string {
	extension := filepath.Ext(basePath)
	name := strings.TrimSuffix(basePath, extension)

	counter := 1
	for {
		newPath := basePath
		if counter > 1 {
			newPath = fmt.Sprintf("%s_%d%s", name, counter, extension)
		}

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}

		counter++
	}
}
