package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/controllers"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.POST("/signup", controllers.SignupHandler)

	auth.POST("/login", controllers.Login)

	auth.POST("/verify_otp", controllers.VerifyOtp)

	auth.POST("/forgot_password", controllers.ForgotPassword)

	auth.POST("/reset_password", controllers.ResetPassword)

	auth.POST("/resend_otp", controllers.ResendOtpHandler)

	auth.POST("/refresh", controllers.RefreshTokenHandler)

	auth.POST("/logout", controllers.Logout)
}
