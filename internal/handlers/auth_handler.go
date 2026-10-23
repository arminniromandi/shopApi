package handlers

import (
	"net/http"
	"time"

	"shopProject/internal/config"
	"shopProject/internal/database"
	"shopProject/internal/models"
	"shopProject/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// SendOTP handles sending an OTP to the user's email
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	code := utils.GenerateOTP()
	expiresAt := time.Now().Add(5 * time.Minute)

	// Invalidate previous OTPs for this email
	database.DB.Where("email = ?", req.Email).Delete(&models.OTPCode{})

	otpRecord := models.OTPCode{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&otpRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification code"})
		return
	}

	// Send code via email (or console fallback)
	go func() {
		_ = utils.SendOTPEmail(h.cfg, req.Email, code)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification code sent to your email",
	})
}

// VerifyOTP verifies the OTP and returns a JWT token
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	var otpRecord models.OTPCode
	err := database.DB.Where("email = ? AND code = ?", req.Email, req.Code).First(&otpRecord).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
		return
	}

	if time.Now().After(otpRecord.ExpiresAt) {
		database.DB.Delete(&otpRecord)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Verification code has expired"})
		return
	}

	// Delete used OTP
	database.DB.Delete(&otpRecord)

	// Check if user exists or register them
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		user = models.User{Email: req.Email}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
