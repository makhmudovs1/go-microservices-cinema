package mongodb

import (
	"context"
	"errors"

	"cinema/bookings/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookingModel представляет сессию базы данных mgo для модели бронирования
type BookingModel struct {
	C *mongo.Collection
}

// Метод All используется для получения всех записей из таблицы bookings
func (m *BookingModel) All() ([]models.Booking, error) {
	// Определяем переменные
	ctx := context.TODO()
	b := []models.Booking{}

	// Ищем все бронирования
	bookingCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = bookingCursor.All(ctx, &b)
	if err != nil {
		return nil, err
	}

	return b, err
}

// FindByID используется для поиска записи бронирования по id
func (m *BookingModel) FindByID(id string) (*models.Booking, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Ищем бронирование по id
	var booking = models.Booking{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&booking)
	if err != nil {
		// Проверяем, что бронирование не найдено
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &booking, nil
}

// Метод Insert используется для добавления новой записи бронирования
func (m *BookingModel) Insert(booking models.Booking) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), booking)
}

// Метод Delete используется для удаления записи бронирования
func (m *BookingModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}
