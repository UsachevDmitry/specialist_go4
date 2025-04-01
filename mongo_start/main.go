package main

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"fmt"
	"log"
)

func main() {
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
	
	defer client.Disconnect(ctx)

	fmt.Println("Connected to MongoDB!")

	//Get all database names
	dbNames, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dbNames)

	// Create new database and collection
	exampleDB := client.Database("exdb")
	fmt.Printf("%T\n", exampleDB)

	exampleCollection := exampleDB.Collection("example")
	fmt.Printf("%T\n", exampleCollection)

	// //Delete full collection
	// defer exampleCollection.Drop(ctx)

	newDoc := bson.D{
		{Key: "strEx", Value: "Hello, Mongo!"},
		{Key: "intEx", Value: "12"},
		{Key: "strSlice", Value: []string{"first", "second", "third"}},
	}
	//=======
	// CREATE
	//=======
	// insert new document
	r, err := exampleCollection.InsertOne(ctx, newDoc)
	if err != nil {
		log.Fatal(err)
	}

	//Print new document "_id"
	fmt.Println(r.InsertedID, r.Acknowledged)


	//insert many documents
	newDocs := []any{ 
		bson.D{
			{Key: "strEx", Value: "Hello, Mongo!2"},
			{Key: "intEx", Value: "124"},
			{Key: "strSlice", Value: []string{"first2", "second2", "third2"}},
		},
		bson.D{
			{Key: "strEx", Value: "Hello, Mongo!3"},
			{Key: "intEx", Value: "1245"},
			{Key: "strSlice", Value: []string{"first3", "second3", "third3"}},
		},
	}
	rs, err := exampleCollection.InsertMany(ctx, newDocs)
	if err != nil {
		log.Fatal(err)
	}
	//Print new documents
	fmt.Println(rs.InsertedIDs)

	//====
	//READ
	//====

	//find document by ObjectID
	c := exampleCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})

	var exampleResult bson.M
	err = c.Decode(&exampleResult)
	if err != nil {
		log.Fatal(err)
	}

	//Print document
	fmt.Printf("\nItem with ID: %v, containing the following:\n", exampleResult["_id"])
	fmt.Println("Key: strEx", exampleResult["strEx"])
	fmt.Println("Key: intEx", exampleResult["intEx"])
	fmt.Println("Key: strSlice", exampleResult["strSlice"])

    // find document by value of ObjectID
	ObjectID, err := bson.ObjectIDFromHex("67ea80b5cdb4be322b262075")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ObjectID)


	second_sr := exampleCollection.FindOne(ctx, bson.M{"_id": bson.M{"$eq": ObjectID}})

	var secondResult bson.M
	err = second_sr.Decode(&secondResult)
	if err != nil {
		log.Fatal(err)
	}

	//Print document
	fmt.Printf("\nItem with ID: %v, containing the following:\n", secondResult["_id"])
	fmt.Println("Key: strEx", secondResult["strEx"])
	fmt.Println("Key: intEx", secondResult["intEx"])
	fmt.Println("Key: strSlice", secondResult["strSlice"])

	//======
	//Find
	//======

	//Get all documents
	allExamples, err := exampleCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var resExamples []bson.M
	if err := allExamples.All(ctx, &resExamples); err != nil {
		log.Fatal(err)
	}
	for _, e := range resExamples {
		fmt.Printf("\nItem with ID: %v, containing the following:\n", e["_id"])
		fmt.Println("Key: strEx", e["strEx"])
		fmt.Println("Key: intEx", e["intEx"])
		fmt.Println("Key: strSlice", e["strSlice"])
	}

	//======
	//Update
	//======
	rUpd, err := exampleCollection.UpdateOne(
		ctx, 
		bson.M{"_id": r.InsertedID},
		bson.D{
			{Key: "$set", Value: bson.M{"strEx": "Change string"}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rUpd.ModifiedCount)

	//Check New Data
	// find document by ObjectID
	srUpd := exampleCollection.FindOne(ctx, bson.M{"_id": r.InsertedID})
	var exampleUpd bson.M
	err = srUpd.Decode(&exampleUpd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nItem witch ID: %v, containing the following change:/n", exampleUpd["_id"])
	fmt.Println("Key: str:Ex", exampleUpd["strEx"])

	//update many
	manyUpd, err := exampleCollection.UpdateMany(ctx,
		bson.D{
			{Key: "intEx", Value: bson.D{{Key: "$gt", Value: 60}}},
		},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "intEx", Value: 60}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(manyUpd.ModifiedCount)

	//======
	// Delete
	//=======
	rDel, err := exampleCollection.DeleteOne(ctx, bson.M{"_id": r.InsertedID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Count of deleted documents", rDel.DeletedCount)
}