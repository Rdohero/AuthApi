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
		Userid    uint
		Productid uint
		Quantity  uint
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
	result := initializers.DB.Where("user_id = ? AND product_id = ?", addCartBody.Userid, addCartBody.Productid).First(&existingCart)
	if result.Error != nil {
		userCart := models.Cart{
			UserID:    addCartBody.Userid,
			ProductID: addCartBody.Productid,
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
	initializers.DB.Preload("Product").Preload("Product.User").Where("user_id = ?", addCartBody.Userid).Find(&cart)

	c.JSON(http.StatusOK, cart)
}

func RemoveCart(c *gin.Context) {
	userID := c.Param("userid")
	productID := c.Param("productid")

	var cart []models.Cart
	del := initializers.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&cart)

	if del.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Terjadi kesalahan dalam mencari pengguna.",
		})
		return
	}

	result := initializers.DB.Preload("Product").Preload("Product.User").Where("user_id = ?", userID).Find(&cart)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Terjadi kesalahan dalam mencari pengguna.",
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}
