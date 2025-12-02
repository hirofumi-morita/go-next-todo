package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateGroupInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Color       *string `json:"color"`
}

type UpdateGroupInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

func GetGroups(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ctx := context.Background()

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
	groups, err := GQLClient.GetGroups(ctx, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func GetGroup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupID := c.Param("id")
	ctx := context.Background()

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
	group, err := GQLClient.GetGroup(ctx, groupID, userIDStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}

func CreateGroup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ctx := context.Background()

	var input CreateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

	var description *string
	if input.Description != "" {
		description = &input.Description
	}

	group, err := GQLClient.CreateGroup(ctx, userIDStr, input.Name, description, input.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group": group})
}

func UpdateGroup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupID := c.Param("id")
	ctx := context.Background()

	var input UpdateGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

	group, err := GQLClient.UpdateGroup(ctx, groupID, userIDStr, input.Name, input.Description, input.Color)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}

func DeleteGroup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupID := c.Param("id")
	ctx := context.Background()

	userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
	deleted, err := GQLClient.DeleteGroup(ctx, groupID, userIDStr)
	if err != nil || !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}
