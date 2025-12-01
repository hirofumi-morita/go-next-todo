package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	ctx := context.Background()

	users, err := GQLClient.GetUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUser(c *gin.Context) {
	userID := c.Param("id")
	ctx := context.Background()

	user, err := GQLClient.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	ctx := context.Background()

	currentUserID, _ := c.Get("user_id")
	targetID, _ := strconv.ParseUint(userID, 10, 64)
	if uint(targetID) == currentUserID.(uint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	deleted, err := GQLClient.DeleteUser(ctx, userID)
	if err != nil || !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func UpdateUserAdmin(c *gin.Context) {
	userID := c.Param("id")
	ctx := context.Background()

	var input struct {
		IsAdmin bool `json:"is_admin"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := GQLClient.UpdateUserAdmin(ctx, userID, input.IsAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
