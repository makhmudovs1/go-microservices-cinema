package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Booking представляет данные профиля бронирования
type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"userid"`
	ShowtimeID string             `bson:"showtimeid"`
	Movies     []string           `bson:"movies"`
}
