package controllers

import (
	"context"
	"log"
	"net/http"

	// "os"

	"github.com/aneesh-oss/todo-app-backend/database"
	"github.com/aneesh-oss/todo-app-backend/models"
	"github.com/aneesh-oss/todo-app-backend/utils"
	"github.com/gin-gonic/gin"

	// "github.com/joho/godotenv" // Import godotenv
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Declare a global variable for the user collection
var userCollection *mongo.Collection

// Initialize the MongoDB connection
func init() {

	// // Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Warning: No .env file found. Using environment variables set in the system.")
	// }

	// // Retrieve the MongoDB URI from environment variables
	// mongoURI := os.Getenv("MONGODB_URI")
	// if mongoURI == "" {
	// 	log.Fatal("Error: MONGODB_URI is not set in the environment variables.")
	// }

	// // Set client options
	// clientOptions := options.Client().ApplyURI(mongoURI)

	// // Connect to MongoDB
	// client, err := mongo.Connect(context.TODO(), clientOptions)
	// if err != nil {
	// 	log.Fatalf("Error connecting to MongoDB: %v", err)
	// }

	// // Ping the MongoDB server to verify connection
	// err = client.Ping(context.TODO(), nil)
	// if err != nil {
	// 	log.Fatalf("Error pinging MongoDB: %v", err)
	// }

	// log.Println("Successfully connected to MongoDB Atlas!")

	// Initialize the user collection

	client, err := database.GetMongoClient()
	if err != nil {
		log.Fatalf("Failed to get MongoDB client: %v", err)
	}
	userCollection = client.Database("todo_app").Collection("users")
}

// SignUp handles user registration
func SignUp(c *gin.Context) {
	var user models.User

	// Bind incoming JSON to the User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the user already exists based on email
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// User with the given email already exists
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Hash the user's password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	// Insert the new user into the MongoDB collection
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error inserting user into MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// SignIn handles user authentication
func SignIn(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind incoming JSON to the credentials struct
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find the user in the database based on email
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		// User not found
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		// Password does not match
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Respond with the JWT token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// package controllers

// import (
// 	"context"
// 	"net/http"

// 	// Import godotenv
// 	"github.com/aneesh-oss/todo-app-backend/models"
// 	"github.com/aneesh-oss/todo-app-backend/utils"
// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"golang.org/x/crypto/bcrypt"
// )

// var userCollection *mongo.Collection

// func init() {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		panic(err)
// 	}
// 	userCollection = client.Database("todo_app").Collection("users")
// }

// func SignUp(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if user already exists
// 	var existingUser models.User
// 	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
// 	if err == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
// 		return
// 	}

// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
// 		return
// 	}
// 	user.Password = string(hashedPassword)

// 	// Insert user
// 	_, err = userCollection.InsertOne(context.TODO(), user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
// }

// func SignIn(c *gin.Context) {
// 	var credentials struct {
// 		Email    string `json:"email" binding:"required,email"`
// 		Password string `json:"password" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&credentials); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Find user
// 	var user models.User
// 	err := userCollection.FindOne(context.TODO(), bson.M{"email": credentials.Email}).Decode(&user)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	// Compare password
// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	// Generate JWT
// 	token, err := utils.GenerateJWT(user.ID.Hex())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }
