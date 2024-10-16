package main

import (
	"log"
	"os"

	"github.com/aneesh-oss/todo-app-backend/controllers"
	"github.com/aneesh-oss/todo-app-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	router := gin.Default()

	// Public routes
	router.POST("/signup", controllers.SignUp)
	router.POST("/signin", controllers.SignIn)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Profile
		protected.GET("/profile", controllers.GetProfile)
		protected.PUT("/profile", controllers.UpdateProfile)

		// Todos
		protected.POST("/todos", controllers.CreateTodo)
		protected.GET("/todos", controllers.GetTodos)
		protected.PUT("/todos/:id", controllers.UpdateTodo)
		protected.DELETE("/todos/:id", controllers.DeleteTodo)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
}

// package main

// import (
//     "github.com/gin-gonic/gin"
//     "github.com/aneesh-oss/todo-app-backend/controllers"
//     "github.com/aneesh-oss/todo-app-backend/middleware"
// )

// func main() {
//     router := gin.Default()

//     // Public routes
//     router.POST("/signup", controllers.SignUp)
//     router.POST("/signin", controllers.SignIn)

//     // Protected routes
//     protected := router.Group("/")
//     protected.Use(middleware.AuthMiddleware())
//     {
//         // Profile
//         protected.GET("/profile", controllers.GetProfile)
//         protected.PUT("/profile", controllers.UpdateProfile)

//         // Todos
//         protected.POST("/todos", controllers.CreateTodo)
//         protected.GET("/todos", controllers.GetTodos)
//         protected.PUT("/todos/:id", controllers.UpdateTodo)
//         protected.DELETE("/todos/:id", controllers.DeleteTodo)
//     }

//     router.Run(":8080")
// }
