package mongodb

import (
	"context"
	"errors"

	"cinema/users/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserModel представляет сессию базы данных mgo для модели пользователя.
type UserModel struct {
	C *mongo.Collection
}

// Метод All используется для получения всех записей из таблицы users.
func (m *UserModel) All() ([]models.User, error) {
	// Определяем переменные
	ctx := context.TODO()
	uu := []models.User{}

	// Ищем всех пользователей
	userCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = userCursor.All(ctx, &uu)
	if err != nil {
		return nil, err
	}

	return uu, err
}

// FindByID используется для поиска записи пользователя по id
func (m *UserModel) FindByID(id string) (*models.User, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Ищем пользователя по id
	var user = models.User{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&user)
	if err != nil {
		// Проверяем, что пользователь не найден
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &user, nil
}

// Метод Insert используется для добавления нового пользователя
func (m *UserModel) Insert(user models.User) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), user)
}

// Метод Delete используется для удаления пользователя
func (m *UserModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}
