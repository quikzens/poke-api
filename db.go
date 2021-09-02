package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB, Ctx = connectDB()

func connectDB() (*mongo.Database, context.Context) {
	// Open mongodb connection
	var uri = "mongodb://localhost:27017"
	// Declare Context type object for managing multiple API requests
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ping the primary / Check the connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	var db = client.Database("pokemon")
	fmt.Println("database connected...")
	return db, ctx
}
