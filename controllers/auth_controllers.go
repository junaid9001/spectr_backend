package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/models"
	"github.com/junaid9001/spectr_backend/services"
	"github.com/junaid9001/spectr_backend/utils"
	"gorm.io/gorm"
)

// register new user
func SignupHandler(c *gin.Context) {
	var creds struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	if err := config.DB.Where("email=?", creds.Email).First(&existingUser).Error; err == nil {
		c.String(http.StatusConflict, "email already registered")
		return
	}

	hashPassword, err := utils.HashPassword(creds.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "error creating account")
		return
	}

	user := models.User{
		Name:           creds.Name,
		Email:          creds.Email,
		HashedPassword: hashPassword,
		Role:           "user",
		IsBlocked:      false,
		IsVerified:     false,
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, "error while creating account")
		return
	}

	//generate otp also saves in db and sends email

	if _, err := services.GenerateOtp(user.ID, user.Email, "signup"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate OTP"})
		return
	}

	if err := config.DB.Model(&models.AppStats{}).Where("id=?", 1).
		UpdateColumn("total_users", gorm.Expr("total_users + ?", 1)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Account created. Please check your email for the verification code.",
	})

}

// Login handler----------------------
func Login(c *gin.Context) {
	var creds struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=4"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	err := config.DB.Where("email=?", creds.Email).First(&existingUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if !existingUser.IsVerified && existingUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "please verify your email to continue"})
		return
	}

	if !utils.CompareHashAndPass(existingUser.HashedPassword, creds.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if existingUser.IsBlocked {
		c.JSON(http.StatusForbidden, gin.H{"error": "your account is blocked due to suspicious activity"})
		return
	}

	accessToken, err2 := utils.GenerateAccessToken(existingUser.ID, existingUser.Role)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, hashedToken, err3 := utils.GenerateRefreshToken()
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	expirestAt := time.Now().Add(7 * 24 * time.Hour)

	if err := utils.SaveRefreshToken(config.DB, existingUser.ID, hashedToken, expirestAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}

	//save rt in cookie

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(time.Until(expirestAt).Seconds()), //max age 7 days
		"/",                                   //path
		"",
		false,
		false,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"role":         existingUser.Role,
		"user_id":      existingUser.ID,
		"access_token": accessToken,
	})
}

// verify otp-------------------------------
func VerifyOtp(c *gin.Context) {
	var creds struct {
		Email   string `json:"email" binding:"required,email"`
		Otp     string `json:"otp" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User

	if err := config.DB.Where("email=?", creds.Email).First(&user).Error; err != nil {
		c.String(http.StatusUnauthorized, "signup first")
		return
	}

	//check if otp is correct , marks otp used and and set user is verified to true

	valid, err := services.ValidateOtp(user.ID, creds.Otp, creds.Purpose)

	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// response
	responseMsg := "OTP verified successfully"
	if creds.Purpose == "signup" {
		responseMsg = "User verified successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": responseMsg,
	})

}

// forgot password hander--------------------
func ForgotPassword(c *gin.Context) {

	var creds struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	var user models.User

	if err := config.DB.Where("email=?", creds.Email).First(&user).Error; err != nil {
		c.String(http.StatusNotFound, "user not found")
		return
	}

	if _, err := services.GenerateOtp(user.ID, user.Email, "reset_password"); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP sent to your email for password reset",
	})
}

//rest password handler (validates otp and update password)-----------------

func ResetPassword(c *gin.Context) {
	var creds struct {
		Email       string `json:"email" binding:"required,email"`
		NewPassword string `json:"new_password" binding:"required,min=4"`
		Otp         string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User

	if err := config.DB.Where("email=?", creds.Email).First(&user).Error; err != nil {
		c.String(http.StatusNotFound, "invalid email")
		return
	}

	valid, err := services.ValidateOtp(user.ID, creds.Otp, "reset_password")
	if !valid {
		c.String(http.StatusBadRequest, "invalid or expired token")
		return
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	hashedPass, err1 := utils.HashPassword(creds.NewPassword)
	if err1 != nil {
		c.String(http.StatusInternalServerError, err1.Error())
		return
	}

	if err := config.DB.Model(models.User{}).Where("email=?", user.Email).
		Updates(map[string]interface{}{"hashed_password": hashedPass, "updated_at": time.Now()}).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	//delete otp after successfull reset

	config.DB.Where("user_id=? AND purpose=?", user.ID, "reset_password").Delete(&models.Otp{})

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password reset successfully",
	})
}

//resend otp handdler ,send new otp to email--------------

func ResendOtpHandler(c *gin.Context) {
	var input struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.String(http.StatusBadRequest, "invalid credentials")
		return
	}

	var user models.User

	if err := config.DB.Where("email=?", input.Email).First(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "invalid credentials")
		return
	}

	if user.IsVerified && input.Purpose == "signup" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already verified"})
		return
	}

	//stores and mail otp
	if _, err := services.GenerateOtp(user.ID, user.Email, input.Purpose); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "OTP resent successfully. Please check your email.",
	})
}

//new accessToken handler (/refresh)

func RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	rt, err1 := utils.ValidateRefreshToken(config.DB, refreshToken)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	var user models.User

	if err := config.DB.Where("id=?", rt.UserId).First(&user).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, err := utils.GenerateAccessToken(rt.UserId, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"access_token": accessToken,
	})

}

//logout (remove refresh token from db

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token required"})
		return
	}

	if err := utils.DeleteRefreshToken(config.DB, refreshToken); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
