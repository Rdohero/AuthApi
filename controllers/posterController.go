package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPoster(c *gin.Context) {
	var poster []models.Poster

	initializers.DB.Find(&poster)

	c.JSON(http.StatusOK, poster)
}
