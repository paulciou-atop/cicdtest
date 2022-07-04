package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (mongoDb MongoDB) CreateOne(dbInf DbInfo, ctx context.Context, document interface{}) (oidHex string, err error) {
	log.Println("createOne=======")

	client := mongoDb.DbIni()
	dbName := dbInf.DbName
	tbName := dbInf.TbName
	collection := client.Database(dbName).Collection(tbName)
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}
	return oid.Hex(), nil

}

func (mongoDb MongoDB) DbIni() *mongo.Client {
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")) //sholud be read by setting
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SessionStatus Service Started")
	//collection = client.Database("nms").Collection("device_data_store")

	//collection = client.Database(dbName).Collection(tbName)

	return client
}

func DbClose(client *mongo.Client) {
	fmt.Println("Closing MongoDB Connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB : %v", err)
	}

	fmt.Println("Stopping the server")
}
func (mongoDb MongoDB) ListDb() (err error) {

	return nil

}
func (mongoDb MongoDB) ReadDb() (err error) {

	return nil

}

func (mongoDb MongoDB) CreatMany() (err error) {

	return nil

}

func (mongoDb MongoDB) UpdateDb() (err error) {

	return nil

}

func (mongoDb MongoDB) DeleteRecord() (err error) {

	return nil

}
