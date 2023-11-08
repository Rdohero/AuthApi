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
	router.POST("/login", controllers.Login)
	userAuth.POST("/getUser", controllers.GetUserById)
	router.PUT("/user/fullname/:id", controllers.UpdateFullnameById)
	router.PUT("/user/username/:id", controllers.UpdateUsernameById)
	router.PUT("/user/foto/:id", controllers.UpdatePhotoProfile)
	adminGroup.DELETE("/user/delete/:id", admin.DeleteUser)
	router.GET("/user/cart/:id", controllers.GetCardByUserId)

	router.GET("/product", controllers.GetProduct)

	router.Run()
}
