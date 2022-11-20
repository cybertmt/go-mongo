package main

import (
	"context"
	"log"
	"time"

	c "gomongo/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	/*
	   Connect to my cluster
	*/
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + c.MongoHost + c.MongoPort))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(databases)

	/*
	   Define my document struct
	*/
	type Post struct {
		Title string `bson:"title,omitempty"`
		Body  string `bson:"body,omitempty"`
	}

	/*
	   Get my collection instance
	*/
	collection := client.Database("cybertest").Collection("cybercollection")

	/*
	   Insert documents
	*/
	docs := []interface{}{
		bson.D{{"title", "World"}, {"body", "Hello World"}},
		bson.D{{"title", "Mars"}, {"body", "Hello Mars"}},
		bson.D{{"title", "Pluto"}, {"body", "Hello Pluto"}},
	}

	res, insertErr := collection.InsertMany(ctx, docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	log.Println(res)
	/*
	   Iterate a cursor and print it
	*/
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var posts []Post
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}
	log.Println(posts)

	_, currErr = collection.DeleteMany(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)
}
