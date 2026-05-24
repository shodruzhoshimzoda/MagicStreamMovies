package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Подключение к MongoDB
func DBConnect() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error downloading .env file: ", err)

	}
	mongoDb := os.Getenv("MONGODB_URI")

	if mongoDb == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}

	fmt.Println("MONGODB_URI is: ", mongoDb)
	clientOptions := options.Client().ApplyURI(mongoDb)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client

}

var Client *mongo.Client = DBConnect()

func OpenCollection(collectionName string) *mongo.Collection {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error downloading .env file: ", err)

	}

	database := os.Getenv("DATABASE_NAME") // Название базы данных
	fmt.Println("DATABASE_NAME: ", database)
	collection := Client.Database(database).Collection(collectionName) // Загрузка коллекций

	if collection == nil {
		return nil
	}

	return collection

}
