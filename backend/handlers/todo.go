package handlers

import (
        "context"
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
)

type CreateTodoInput struct {
        Title       string  `json:"title" binding:"required"`
        Description string  `json:"description"`
        GroupID     *string `json:"group_id"`
}

type UpdateTodoInput struct {
        Title       string  `json:"title"`
        Description string  `json:"description"`
        Completed   *bool   `json:"completed"`
        GroupID     *string `json:"group_id"`
}

func GetTodos(c *gin.Context) {
        userID, _ := c.Get("user_id")
        ctx := context.Background()

        userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
        todos, err := GQLClient.GetTodos(ctx, userIDStr)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func GetTodo(c *gin.Context) {
        userID, _ := c.Get("user_id")
        todoID := c.Param("id")
        ctx := context.Background()

        userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
        todo, err := GQLClient.GetTodo(ctx, todoID, userIDStr)
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func CreateTodo(c *gin.Context) {
        userID, _ := c.Get("user_id")
        ctx := context.Background()

        var input CreateTodoInput
        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
        todo, err := GQLClient.CreateTodo(ctx, userIDStr, input.Title, input.Description, input.GroupID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"todo": todo})
}

func UpdateTodo(c *gin.Context) {
        userID, _ := c.Get("user_id")
        todoID := c.Param("id")
        ctx := context.Background()

        var input UpdateTodoInput
        if err := c.ShouldBindJSON(&input); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)

        var title, description *string
        if input.Title != "" {
                title = &input.Title
        }
        if input.Description != "" {
                description = &input.Description
        }

        todo, err := GQLClient.UpdateTodo(ctx, todoID, userIDStr, title, description, input.Completed, input.GroupID)
        if err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"todo": todo})
}

func DeleteTodo(c *gin.Context) {
        userID, _ := c.Get("user_id")
        todoID := c.Param("id")
        ctx := context.Background()

        userIDStr := strconv.FormatUint(uint64(userID.(uint)), 10)
        deleted, err := GQLClient.DeleteTodo(ctx, todoID, userIDStr)
        if err != nil || !deleted {
                c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
