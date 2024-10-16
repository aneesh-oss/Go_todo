package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a to-do item

type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title" binding:"required"`
	Description string             `json:"description" bson:"description"  binding:"required"`
	Completed   bool               `json:"completed" bson:"completed"`
}

// type Todo struct {
// 	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
// 	Title       string             `json:"title" bson:"title" binding:"required"`
// 	Description string             `json:"description" bson:"description"  binding:"required"`
// 	Completed   bool               `json:"completed" bson:"completed"`
// }

// type Todo struct {
//     ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
//     UserID    primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
//     Title     string             `json:"title" binding:"required"`
//     Completed bool               `json:"completed"`
// }
