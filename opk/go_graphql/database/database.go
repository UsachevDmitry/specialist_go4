package database

import (
	"context"
	"go_graphql/graph/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)


func timePtr(t time.Time) *time.Time {
	return &t
}

type DB struct {
	client *mongo.Client
}

func Connect() *DB{
	ctx := context.TODO()
    // URI с логином и паролем
    uri := "mongodb://root:P%40ssw0rd@127.0.0.1:27017/pets?authSource=admin&authMechanism=SCRAM-SHA-256"

    // Настройка клиента
    clientOptions := options.Client().ApplyURI(uri)

    // Подключение
    client, err := mongo.Connect(clientOptions)
    if err != nil {
        log.Fatal(err)
    }
	// Проверка подключения
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB {
		client: client,
	}
}

func collectionHelper(db *DB, collectionName string) *mongo.Collection {
	return db.client.Database("blog_post").Collection(collectionName)
}

func (db *DB) GetPost(id string) *model.Post {
	collection := collectionHelper(db, "ports")
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	Id, err:= bson.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": Id}

	var post model.Post
	
	err = collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}
	return &post
} 

func (db *DB) CreatePost(postInfo *model.NewPost) *model.Post {
	collection := collectionHelper(db, "ports")
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	postInfo.PublishedAt = timePtr(time.Now())
	postInfo.PublishedAt = timePtr(time.Now())

	result, err := collection.InsertOne(ctx, postInfo)
	if err != nil {
		log.Fatal(err)
	}


	newPost := &model.Post{
		ID: result.InsertedID.(bson.ObjectID).Hex(),
		Title: postInfo.Title,
		Content: postInfo.Content,
		Author: *postInfo.Author,
		Hero: *postInfo.Hero,
		PublishedAt: *postInfo.PublishedAt,
		UpdatedAt: *postInfo.UpdatedAt,
	}
	return newPost
}