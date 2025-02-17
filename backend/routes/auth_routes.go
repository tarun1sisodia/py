package routes

import (
	"github.com/gin-gonic/gin"
	"smart-attendance/controllers"
)

// AuthRoutes defines the authentication routes.
func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register/teacher", controllers.RegisterTeacher)
		authGroup.POST("/register/student", controllers.RegisterStudent)
		authGroup.POST("/verify-otp", controllers.VerifyOTP)
		authGroup.POST("/login/teacher", controllers.LoginTeacher)
		authGroup.POST("/login/student", controllers.LoginStudent)
		authGroup.POST("/reset-password", controllers.ResetPassword)
	}

	// Firebase Authentication Routes
	firebaseAuthGroup := router.Group("/firebaseauth")
	{
		firebaseAuthGroup.Use(controllers.VerifyFirebaseToken())
		firebaseAuthGroup.POST("/createorupdate", controllers.CreateOrUpdateUser)
	}
}
