package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminPower(c *gin.Context) {
	var body struct {
		Admin    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		c.Abort()
		return
	}

	var admin = "Ridho"
	var password = "Ridho8297@8927"

	if admin != body.Admin {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "admin atau password salah",
		})
		c.Abort()
		return
	}

	if password != body.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "admin atau password salah",
		})
		c.Abort()
		return
	}

	c.Next()
}
