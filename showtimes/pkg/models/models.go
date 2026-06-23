package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShowTime представляет данные профиля сеанса
type ShowTime struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Date      string             `bson:"date,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	Movies    []string           `bson:"movies,omitempty"`
}
