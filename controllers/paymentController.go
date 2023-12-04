package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Payment(c *gin.Context) {
	var pay struct {
		User    int
		Payment int
		Invoice int
	}

	if c.Bind(&pay) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	var invoice models.Invoice
	result := initializers.DB.First(&invoice, pay.Invoice)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": "Invoice not found",
		})
		return
	}

	updateResult := initializers.DB.Model(&invoice).Updates(models.Invoice{
		PaymentID: pay.Payment,
		StatusID:  2,
	})

	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to update invoice",
		})
		return
	}

	var invoice2 []models.Invoice
	initializers.DB.Preload("Status").Preload("InvoiceItems").Preload("InvoiceItems.Product").Preload("PaymentItem").Preload("InvoiceItems.Product.User").Where("user_id", pay.User).Find(&invoice2)

	c.JSON(http.StatusOK, invoice2)
}
