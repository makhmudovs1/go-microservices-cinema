package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type websiteMovie struct {
	ID        primitive.ObjectID
	Title     string
	Director  string
	Rating    float32
	CreatedOn time.Time
}

type websiteShowTime struct {
	ID        primitive.ObjectID
	Date      string
	CreatedAt time.Time
	Movies    []string
}

type websiteBooking struct {
	ID         primitive.ObjectID
	UserID     string
	ShowtimeID string
	Movies     []string
}
