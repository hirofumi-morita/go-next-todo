package handlers

import (
        "net/http"
        "os"
        "time"
        "todo-app/config"
        "todo-app/models"

        "github.com/gin-gonic/gin"
        "github.com/golang-jwt/jwt/v5"
)

type RegisterInput struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
        var input RegisterInput
        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        var existingUser models.User
        if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
                c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
                return
        }

        user := models.User{Email: input.Email}
        if err := user.HashPassword(input.Password); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
                return
        }

        var userCount int64
        config.DB.Model(&models.User{}).Count(&userCount)
        if userCount == 0 {
                user.IsAdmin = true
        }

        if err := config.DB.Create(&user).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
                return
        }

        token, err := generateToken(user)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{
                "message": "User registered successfully",
                "token":   token,
                "user": gin.H{
                        "id":       user.ID,
                        "email":    user.Email,
                        "is_admin": user.IsAdmin,
                },
        })
}

func Login(c *gin.Context) {
        var input LoginInput
        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        var user models.User
        if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
                return
        }

        if !user.CheckPassword(input.Password) {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
                return
        }

        token, err := generateToken(user)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "message": "Login successful",
                "token":   token,
                "user": gin.H{
                        "id":       user.ID,
                        "email":    user.Email,
                        "is_admin": user.IsAdmin,
                },
        })
}

func GetMe(c *gin.Context) {
        userID, _ := c.Get("user_id")
        
        var user models.User
        if err := config.DB.First(&user, userID).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{
                "user": gin.H{
                        "id":       user.ID,
                        "email":    user.Email,
                        "is_admin": user.IsAdmin,
                },
        })
}

func generateToken(user models.User) (string, error) {
        secret := os.Getenv("SESSION_SECRET")
        if secret == "" {
                return "", jwt.ErrTokenSignatureInvalid
        }

        claims := jwt.MapClaims{
                "user_id":  user.ID,
                "email":    user.Email,
                "is_admin": user.IsAdmin,
                "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString([]byte(secret))
}
