package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"time"

	"shopProject/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateOTP produces a 6-digit verification code
func GenerateOTP() string {
	digits := "0123456789"
	b := make([]byte, 6)
	_, _ = io.ReadFull(rand.Reader, b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}

// SendOTPEmail sends the verification code or logs to stdout if SMTP is not configured
func SendOTPEmail(cfg *config.Config, toEmail, otpCode string) error {
	if cfg.SMTPUser == "" || cfg.SMTPPass == "" {
		log.Printf("\n=========================================\n[DEV OTP SIMULATION]\nSending OTP to: %s\nCode: %s (Expires in 5 minutes)\n=========================================\n", toEmail, otpCode)
		return nil
	}

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	subject := "Subject: Verification Code - Shop App\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; direction: rtl; text-align: right;">
			<h2>کد تأیید ورود به فروشگاه</h2>
			<p>کد یک‌بار مصرف شما برای ورود:</p>
			<h1 style="color: #2b6cb0; letter-spacing: 5px;">%s</h1>
			<p>این کد به مدت ۵ دقیقه معتبر است.</p>
		</div>
	`, otpCode)

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	return smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{toEmail}, msg)
}

// GenerateJWT creates a signed JWT token
func GenerateJWT(userID uint, email, secret string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days valid
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT verifies the token and parses claims
func ValidateJWT(tokenStr, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
