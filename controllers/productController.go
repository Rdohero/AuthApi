package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProduct(c *gin.Context) {
	var product []models.Product

	initializers.DB.Find(&product)

	c.JSON(http.StatusOK, product)
}
