package main

import (
	"AuthApi/controllers"
	"AuthApi/controllers/admin"
	"AuthApi/initializers"
	"AuthApi/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AdminPower)

	userAuth := router.Group("/userAuth")
	userAuth.Use(middleware.RequiredAuth)

	router.Static("/images", "images/")

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

	router.POST("/save-image", controllers.SaveImage)

	router.Run()
}
