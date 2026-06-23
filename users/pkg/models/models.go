package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User представляет данные профиля пользователя
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	LastName string             `bson:"lastname,omitempty"`
	Email    string             `bson:"email,omitempty"`
}
