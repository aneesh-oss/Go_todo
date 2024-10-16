package controllers

import (
	"context"
	// "fmt"
	"log"
	"net/http"

	"github.com/aneesh-oss/todo-app-backend/database"
	"github.com/aneesh-oss/todo-app-backend/models"
	"github.com/gin-gonic/gin"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo"

	logger "github.com/sirupsen/logrus"
)

var todoCollection *mongo.Collection

func init() {

	client, err := database.GetMongoClient()
	if err != nil {
		log.Fatalf("Failed to get MongoDB client: %v", err)
	}
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// client, err := mongo.Connect(context.TODO(), clientOptions)
	// if err != nil {
	// 	panic(err)
	// }
	todoCollection = client.Database("todo_app").Collection("todos")
}

// func CreateTodo(c *gin.Context) {
// 	userIDHex, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
// 		return
// 	}

// 	var todo models.Todo
// 	if err := c.ShouldBindJSON(&todo); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	todo.UserID = userID
// 	todo.Completed = false

// 	_, err = todoCollection.InsertOne(context.TODO(), todo)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating todo"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, todo)
// }

// CreateTodo handles the creation of a new todo
func CreateTodo(c *gin.Context) {
	// Retrieve the userID from the context (set by AuthMiddleware)
	userIDHex, exists := c.Get("userID")
	if !exists {
		log.Println("UserID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert the userID from hex string to ObjectID
	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		log.Printf("Invalid userID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var todo models.Todo
	// Bind the incoming JSON to the Todo struct
	if err := c.ShouldBindJSON(&todo); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Assign the userID to the Todo
	todo.UserID = userID
	todo.Completed = false // Default value

	// Insert the Todo into MongoDB
	result, err := userCollection.InsertOne(context.TODO(), todo)
	if err != nil {
		log.Printf("Error inserting todo into MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating todo"})
		return
	}

	// Assign the insertedID to the Todo struct
	todo.ID = result.InsertedID.(primitive.ObjectID)

	log.Printf("Todo created successfully with ID: %s", todo.ID.Hex())
	c.JSON(http.StatusCreated, todo)
}

func GetTodos(c *gin.Context) {
	userIDHex, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cursor, err := todoCollection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching todos"})
		return
	}
	defer cursor.Close(context.TODO())

	var todos []models.Todo
	if err = cursor.All(context.TODO(), &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing todos"})
		return
	}

	// Log the fetched todos in detail
	for _, todo := range todos {
		logger.Infof("Fetched Todo - ID: %s, Title: %s, Description: %s, Completed: %v",
			todo.ID.Hex(), todo.Title, todo.Description, todo.Completed)
	}

	// // Log the details of the todos in the console
	// fmt.Println("User ID:", userIDHex)
	// fmt.Println("User's Todos:", todos)

	c.JSON(http.StatusOK, todos)
}

func UpdateTodo(c *gin.Context) {
	userIDHex, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	todoIDHex := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var updateData struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{}
	if updateData.Title != "" {
		update["title"] = updateData.Title
	}
	update["completed"] = updateData.Completed

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data to update"})
		return
	}

	result, err := todoCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": todoID, "user_id": userID},
		bson.M{"$set": update},
	)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func DeleteTodo(c *gin.Context) {
	userIDHex, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	todoIDHex := c.Param("id")
	todoID, err := primitive.ObjectIDFromHex(todoIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	result, err := todoCollection.DeleteOne(context.TODO(), bson.M{"_id": todoID, "user_id": userID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
