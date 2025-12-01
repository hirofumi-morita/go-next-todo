package main

import (
	"log"
	"todo-app/config"
	"todo-app/graph"
	"todo-app/handlers"
	"todo-app/middleware"
	"todo-app/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	if err := config.DB.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")

	resolver := graph.NewResolver(config.DB)
	gqlClient := graph.NewClient(resolver)
	handlers.InitGraphQLClient(gqlClient)
	log.Println("GraphQL layer initialized successfully")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", handlers.GetMe)

			todos := protected.Group("/todos")
			{
				todos.GET("", handlers.GetTodos)
				todos.GET("/:id", handlers.GetTodo)
				todos.POST("", handlers.CreateTodo)
				todos.PUT("/:id", handlers.UpdateTodo)
				todos.DELETE("/:id", handlers.DeleteTodo)
			}

			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", handlers.GetAllUsers)
				admin.GET("/users/:id", handlers.GetUser)
				admin.DELETE("/users/:id", handlers.DeleteUser)
				admin.PATCH("/users/:id", handlers.UpdateUserAdmin)
			}
		}
	}

	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
