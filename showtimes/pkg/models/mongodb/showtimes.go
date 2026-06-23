package mongodb

import (
	"context"
	"errors"

	"cinema/showtimes/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ShowTimeModel представляет сессию базы данных mgo для модели сеанса.
type ShowTimeModel struct {
	C *mongo.Collection
}

// Метод All используется для получения всех записей из таблицы showtimes.
func (m *ShowTimeModel) All() ([]models.ShowTime, error) {
	// Определяем переменные
	ctx := context.TODO()
	st := []models.ShowTime{}

	// Ищем все сеансы
	showtimeCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = showtimeCursor.All(ctx, &st)
	if err != nil {
		return nil, err
	}

	return st, err
}

// FindByID используется для поиска записи сеанса по id
func (m *ShowTimeModel) FindByID(id string) (*models.ShowTime, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Ищем сеанс по id
	var showtime = models.ShowTime{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&showtime)
	if err != nil {
		// Проверяем, что сеанс не найден
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &showtime, nil
}

// FindByDate используется для поиска записи сеанса по дате
func (m *ShowTimeModel) FindByDate(date string) (*models.ShowTime, error) {
	// Ищем сеанс по дате
	var showtime = models.ShowTime{}
	err := m.C.FindOne(context.TODO(), bson.M{"date": date}).Decode(&showtime)
	if err != nil {
		// Проверяем, что сеанс не найден
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &showtime, nil
}

// Метод Insert используется для добавления новой записи сеанса
func (m *ShowTimeModel) Insert(showtime models.ShowTime) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), showtime)
}

// Метод Delete используется для удаления записи сеанса
func (m *ShowTimeModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}
