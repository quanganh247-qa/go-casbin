package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

var jwtKey = []byte("your_secret_key") // move to env/config in production

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	Role      string
	CreatedAt time.Time
}

// ✅ Handler: Đăng ký user mới
func RegisterUser(db *gorm.DB, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user User
		// Ensure the database schema is ready
		if err := db.AutoMigrate(&User{}); err != nil {
			log.Printf("Migration error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database migration failed"})
			return
		}
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. Check nếu username đã tồn tại
		var existing User
		if err := db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		// 2. Lưu user vào DB (role mặc định là user)
		user = User{
			Username: req.Username,
			Password: req.Password, // TODO: mã hoá password (bcrypt)
			Role:     "user",
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// 3. Thêm Casbin policy: g, username, user
		_, err := enforcer.AddGroupingPolicy(user.Username, "user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
			return
		}

		// 4. Thêm Casbin policy: p, user, /profile, GET
		_, err = enforcer.AddPolicy(user.Username, "/profile", "GET")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign policy"})
			return
		}

		// 4. Tạo JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Username,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenStr, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
			return
		}

		// 5. Trả về token
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"token":   tokenStr,
		})
	}
}
