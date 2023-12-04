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

func GetProductById(c *gin.Context) {
	var id = c.Param("id")
	var product []models.Product

	initializers.DB.Find(&product, "id = ?", id)

	c.JSON(http.StatusOK, product)
}

func SearchProduct(c *gin.Context) {
	query := c.Query("q")

	var product []models.Product

	initializers.DB.Where("name LIKE ?", "%"+query+"%").Find(&product)

	c.JSON(http.StatusOK, product)
}
