package controllers

import (
	"AuthApi/initializers"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetInvoice(c *gin.Context) {
	id := c.Param("id")

	var invoice []models.Invoice
	initializers.DB.Preload("Status").Preload("InvoiceItems").Preload("InvoiceItems.Product").Preload("PaymentItem").Preload("InvoiceItems.Product.User").Where("user_id", id).Find(&invoice)

	c.JSON(http.StatusOK, invoice)
}

func MetodePembayaran(c *gin.Context) {
	var pembayaran []models.Payment

	initializers.DB.Find(&pembayaran)

	c.JSON(http.StatusOK, pembayaran)
}

func MakeInvoice(c *gin.Context) {
	var Invoice struct {
		UserID      uint   `json:"user_id"`
		TotalAmount uint   `json:"total_amount"`
		Address     string `json:"address"`
	}

	if c.Bind(&Invoice) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	invoices := models.Invoice{
		PaymentID:   1,
		UserID:      Invoice.UserID,
		TotalAmount: Invoice.TotalAmount,
		StatusID:    1,
		Address:     Invoice.Address,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	create := initializers.DB.Create(&invoices).Find(&invoices)

	if create.Error == nil {
		var invoice []models.Invoice

		initializers.DB.Preload("Status").Preload("PaymentItem").Preload("InvoiceItems").Where("id = ?", invoices.ID).Find(&invoice)

		c.JSON(http.StatusOK, invoice)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": create.Error,
		})
	}
}

func MakeInvoiceItem(c *gin.Context) {
	var InvoiceItem1 struct {
		InvoiceID  uint `json:"invoice_id"`
		ProductID  uint `json:"product_id"`
		Quantity   uint `json:"quantity"`
		TotalPrice uint `json:"total_price"`
	}

	if c.Bind(&InvoiceItem1) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}

	invoices := models.InvoiceItem{
		InvoiceID:  InvoiceItem1.InvoiceID,
		ProductID:  InvoiceItem1.ProductID,
		Quantity:   InvoiceItem1.Quantity,
		TotalPrice: InvoiceItem1.TotalPrice,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	create := initializers.DB.Create(&invoices)

	if create.Error == nil {
		var invoice1 models.Invoice
		if err := initializers.DB.Where("id = ?", InvoiceItem1.InvoiceID).First(&invoice1).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Invoice not found",
			})
			return
		}

		if invoice1.PaymentID == 1 {
		} else {
			if err := initializers.DB.Model(&invoice1).Update("status_id", 2).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to update invoice status",
				})
				return
			}
		}

		var invoice []models.Invoice
		initializers.DB.Preload("Status").Preload("InvoiceItems").Preload("InvoiceItems.Product").Preload("PaymentItem").Preload("InvoiceItems.Product.User").Where("user_id").Find(&invoice)

		c.JSON(http.StatusOK, invoice)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": create.Error,
		})
	}
}

func InvoiceStatusNotPaid(c *gin.Context) {
	id := c.Param("id")
	var invoice []models.Invoice
	status := initializers.DB.Preload("Status").Where("user_id = ? AND status_id = ?", id, 1).Find(&invoice)

	if status.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": status.Error,
		})
		return
	}

	c.JSON(http.StatusOK, invoice)
}
