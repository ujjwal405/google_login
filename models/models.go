package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserSignup struct {
	ID       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email" json:"email" validate:"email,required"`
	Password string             `bson:"password" json:"password" validate:"required,min=5,max=10"`
	Phone    string             `bson:"phone" json:"phone"`
	Username string             `bson:"user_name" json:"username"`
	User_id  string             `bson:"user_id" json:"user_id"`
	Isvalid  bool               `json:"isvalid"`
}
type UserData struct {
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}
type ContextKey string
