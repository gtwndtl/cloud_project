package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/services"
	"example.com/se/metrics"
)

type (
	Authen struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	signUp struct {
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Age       uint8     `json:"age"`
		Password  string    `json:"password"`
		Role      string    `json:"role"`
		BirthDay  time.Time `json:"birthday"`
		GenderID  uint      `json:"gender_id"`
	}
)

// SignUp handles user registration
func SignUp(c *gin.Context) {
	var payload signUp
	if err := c.ShouldBindJSON(&payload); err != nil {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var userCheck entity.Users
	result := db.Where("email = ?", payload.Email).First(&userCheck)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if userCheck.ID != 0 {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
		return
	}

	hashedPassword, _ := config.HashPassword(payload.Password)
	user := entity.Users{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Age:       payload.Age,
		Password:  hashedPassword,
		Role:      payload.Role,
		BirthDay:  payload.BirthDay,
		GenderID:  payload.GenderID,
	}

	if err := db.Create(&user).Error; err != nil {
		metrics.UserSignupFailuresTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Successful signup
	metrics.UserSignupsTotal.Inc()

	// Update total users gauge
	var count int64
	db.Model(&entity.Users{}).Count(&count)
	metrics.UsersTotal.Set(float64(count))

	c.JSON(http.StatusCreated, gin.H{"message": "Sign-up successful"})
}

// SignIn handles user login
func SignIn(c *gin.Context) {
	var payload Authen
	var user entity.Users

	if err := c.ShouldBindJSON(&payload); err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := config.DB().Where("email = ?", payload.Email).First(&user).Error
	if err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	// Successful login
	metrics.UserLoginsTotal.Inc()

	jwtWrapper := services.JwtWrapper{
		SecretKey:       "SvNQpBN8y3qlVrsGAYYWoJJk56LtzFHx",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}
	token, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		metrics.UserLoginFailuresTotal.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token_type": "Bearer",
		"token":      token,
		"id":         user.ID,
	})
}
