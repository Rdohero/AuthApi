package main

import (
	"AuthApi/controllers"
	"AuthApi/controllers/admin"
	"AuthApi/initializers"
	"AuthApi/middleware"
	"AuthApi/models"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func Transaction(c *gin.Context) {
	// Start a transaction
	var Transaksi struct {
		Penjual string
		Pembeli string
		Harga   int
	}

	c.Bind(&Transaksi)

	tx := initializers.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Transaction failed"})
		}
	}()

	// Check for errors during the transaction
	if tx.Error != nil {
		c.JSON(500, gin.H{"error": "Error starting transaction"})
		return
	}

	// Retrieve user1 and user2 from the database
	var penjual, pembeli models.TestTransaction
	tx.First(&pembeli, "name = ?", Transaksi.Pembeli)

	if tx.Error != nil || pembeli.ID == 0 {
		tx.Rollback()
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Check if user1 has sufficient balance
	if pembeli.Balance < Transaksi.Harga {
		tx.Rollback()
		c.JSON(400, gin.H{"error": "Insufficient balance"})
		return
	}

	pembeli.Balance -= Transaksi.Harga
	tx.Save(&pembeli)

	// Check if the users exist

	tx.First(&penjual, "name = ?", Transaksi.Penjual)

	if tx.Error != nil || penjual.ID == 0 {
		tx.Rollback()
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Perform the transaction
	penjual.Balance += Transaksi.Harga

	// Update user1 and user2 balances
	tx.Save(&penjual)

	// Commit the transaction
	tx.Commit()

	c.JSON(200, gin.H{"message": "Transaction successful"})
}

func main() {
	router := gin.Default()
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AdminPower)

	userAuth := router.Group("/userAuth")
	userAuth.Use(middleware.RequiredAuth)

	router.Static("/images", "images/")

	router.POST("/perform-transaction", Transaction)

	router.POST("/signup", controllers.Signup)
	router.POST("/emailve", controllers.OtpEmailVer)
	router.POST("/resendOtp", controllers.ResendOtpEmailPassVer)
	router.POST("/forgot", controllers.ForgotPassword)
	router.POST("/login", controllers.Login)
	userAuth.POST("/getUser", controllers.GetUserById)
	router.PUT("/user/fullname/:id", controllers.UpdateFullnameById)
	router.PUT("/user/username/:id", controllers.UpdateUsernameById)
	router.PUT("/user/foto/:id", controllers.UpdatePhotoProfile)
	adminGroup.DELETE("/user/delete/:id", admin.DeleteUser)

	router.GET("/product", controllers.GetProduct)
	router.GET("/product/:id", controllers.GetProductById)
	router.GET("/product/search", controllers.SearchProduct)

	router.POST("/cart", controllers.AddCart)
	router.DELETE("/cart/:userid/:productid", controllers.RemoveCart)
	router.GET("/user/cart/:id", controllers.GetCartByUserId)

	router.GET("/poster", controllers.GetPoster)

	router.GET("/invoice/:id", controllers.GetInvoice)
	router.POST("/create/invoice", controllers.MakeInvoice)
	router.POST("/create/invoiceitem", controllers.MakeInvoiceItem)
	router.GET("/invoice/notpaid/:id", controllers.InvoiceStatusNotPaid)

	router.GET("/payment", controllers.MetodePembayaran)
	router.POST("/payment", controllers.Payment)

	router.Run()
}
