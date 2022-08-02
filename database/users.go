package database

import (
	"context"
	"time"

	"github.com/schattenbrot/simple-login-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *dbRepo) CreateUser(user models.User) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.DB.Collection("users")

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	oid := res.InsertedID.(primitive.ObjectID).Hex()

	return &oid, nil
}

func (m *dbRepo) GetAllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users := []*models.User{}

	collection := m.DB.Collection("users")

	opts := options.Find().SetProjection(bson.M{
		"password": 0,
	})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User

		cursor.Decode(&user)

		users = append(users, &user)
	}

	return users, nil
}

func (m *dbRepo) GetUserById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.DB.Collection("users")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	var user models.User

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *dbRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.DB.Collection("users")

	filter := bson.M{"$or": []interface{}{
		bson.M{"email": username},
		bson.M{"name": username},
	}}

	var user models.User

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *dbRepo) UpdateUser(id string, user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.DB.Collection("users")

	update := bson.M{"$set": user}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}

	return nil
}
