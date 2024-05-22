package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveImage(c *gin.Context) {
	var Image struct {
		Image []byte
	}

	fmt.Println(c.Bind(&Image))

	if err := c.Bind(&Image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	Images := models.Image{
		Image: Image.Image,
	}

	create := initializers.DB.Create(&Images)

	if create.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"Succes": create,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": create.Error,
		})
	}
}
