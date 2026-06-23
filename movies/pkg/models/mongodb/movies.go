package mongodb

import (
	"context"
	"errors"

	"cinema/movies/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MovieModel представляет сессию базы данных mgo для модели фильма.
type MovieModel struct {
	C *mongo.Collection
}

// Метод All используется для получения всех записей из таблицы movies.
func (m *MovieModel) All() ([]models.Movie, error) {
	// Определяем переменные
	ctx := context.TODO()
	mm := []models.Movie{}

	// Ищем все фильмы
	movieCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = movieCursor.All(ctx, &mm)
	if err != nil {
		return nil, err
	}

	return mm, err
}

// FindByID используется для поиска записи фильма по id
func (m *MovieModel) FindByID(id string) (*models.Movie, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Ищем фильм по id
	var movie = models.Movie{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&movie)
	if err != nil {
		// Проверяем, что фильм не найден
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &movie, nil
}

// Метод Insert используется для добавления новой записи фильма
func (m *MovieModel) Insert(movie models.Movie) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), movie)
}

// Метод Delete используется для удаления записи фильма
func (m *MovieModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}
