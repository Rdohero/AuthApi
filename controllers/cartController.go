package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AddCart(c *gin.Context) {
	var addCartBody struct {
		ID        uint
		UserID    uint `json:"user_id"`
		ProductID uint `json:"product_id"`
		Quantity  uint `json:"quantity"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	if c.Bind(&addCartBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var existingCart models.Cart
	result := initializers.DB.Where("user_id = ? AND product_id = ?", addCartBody.UserID, addCartBody.ProductID).First(&existingCart)
	if result.Error != nil {
		userCart := models.Cart{
			UserID:    addCartBody.UserID,
			ProductID: addCartBody.ProductID,
			Quantity:  addCartBody.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := initializers.DB.Create(&userCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create cart",
			})
			return
		}
	} else {
		existingCart.Quantity += addCartBody.Quantity
		existingCart.UpdatedAt = time.Now()
		if err := initializers.DB.Save(&existingCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update cart",
			})
			return
		}
	}

	var cart []models.Cart
	initializers.DB.Preload("Product").Preload("User").Where("user_id = ?", addCartBody.UserID).Find(&cart)

	c.JSON(http.StatusOK, cart)
}
