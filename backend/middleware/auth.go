package middleware

import (
        "net/http"
        "os"
        "strings"

        "github.com/gin-gonic/gin"
        "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                secret := os.Getenv("SESSION_SECRET")
                if secret == "" {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
                        c.Abort()
                        return
                }

                authHeader := c.GetHeader("Authorization")
                if authHeader == "" {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
                        c.Abort()
                        return
                }

                tokenString := strings.TrimPrefix(authHeader, "Bearer ")
                if tokenString == authHeader {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
                        c.Abort()
                        return
                }

                token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                                return nil, jwt.ErrSignatureInvalid
                        }
                        return []byte(secret), nil
                })

                if err != nil || !token.Valid {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
                        c.Abort()
                        return
                }

                claims, ok := token.Claims.(jwt.MapClaims)
                if !ok {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
                        c.Abort()
                        return
                }

                userIDFloat, ok := claims["user_id"].(float64)
                if !ok {
                        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
                        c.Abort()
                        return
                }
                userID := uint(userIDFloat)

                isAdmin, ok := claims["is_admin"].(bool)
                if !ok {
                        isAdmin = false
                }

                c.Set("user_id", userID)
                c.Set("is_admin", isAdmin)
                c.Next()
        }
}

func AdminMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
                isAdmin, exists := c.Get("is_admin")
                if !exists || !isAdmin.(bool) {
                        c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
                        c.Abort()
                        return
                }
                c.Next()
        }
}
