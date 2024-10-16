package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name     string             `json:"name" binding:"required"`
    Email    string             `json:"email" binding:"required,email"`
    Password string             `json:"password,omitempty" binding:"required"`
}
