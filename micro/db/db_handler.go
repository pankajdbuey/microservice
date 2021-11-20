package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	//i IDB
	collection *mongo.Collection
	ctx        context.Context
	cancel     context.CancelFunc
)

// type IDB interface{
// 	Get(id primitive.ObjectID) (*Pet, error)
// }

type Pet struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type      string             `json:"type,omitempty" bson:"type,omitempty"`
	Breed     string             `json:"breed,omitempty" bson:"breed,omitempty"`
	BirthDate *time.Time         `json:"birthdate,omitempty" bson:"birthdate,omitempty"`
}

func init() {
	// i = &Pet{}
	//var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db_server:27017"))
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	if err != nil {
		fmt.Println("Error: while connecting to database", err)
		return
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Error: not able to ping server, ", err)
		return
	}
	collection = client.Database("petdatastore").Collection("pet")

}

func (*Pet) Insert(p *Pet) (primitive.ObjectID, error) {
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// defer cancel()

	if collection != nil {
		res, err := collection.InsertOne(context.Background(), p)
		if err != nil {
			fmt.Printf("Error: while inserting %v in db, err = %v", p, err)
			return primitive.ObjectID{}, err
		}
		fmt.Println("Successfully inserted in DB", p)
		return res.InsertedID.(primitive.ObjectID), nil
	}

	return primitive.ObjectID{}, errors.New("did not get db instance, please check if the db is up and running")

}

func (p *Pet) GetAll() ([]Pet, error) {
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// defer cancel()
	var pets []Pet
	if collection != nil {
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			fmt.Printf("Error: while fetching data , err = %v", err)
			return nil, err
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var pet Pet
			cursor.Decode(&pet)
			pets = append(pets, pet)
		}
		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return pets, nil
	}

	return nil, errors.New("did not get db instance, please check if the db is up and running")

}

func (p *Pet) Get(id primitive.ObjectID) (*Pet, error) {
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// defer cancel()
	var pet Pet
	if collection != nil {
		err := collection.FindOne(context.Background(), Pet{ID: id}).Decode(&pet)
		if err != nil {
			return nil, err
		}
		return &pet, nil
	}

	return nil, errors.New("did not get db instance, please check if the db is up and running")

}

func (p *Pet) Delete(id primitive.ObjectID) (int64, error) {
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// defer cancel()
	if collection != nil {
		r, err := collection.DeleteOne(context.Background(), Pet{ID: id})
		if err != nil {
			return 0, err
		}
		return r.DeletedCount, nil
	}

	return 0, errors.New("did not get db instance, please check if the db is up and running")

}

func (*Pet) Update(p *Pet) (*Pet, error) {
	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// defer cancel()
	if collection != nil {
		u := bson.D{
			{"$set", bson.D{{"type", p.Type}, {"breed", p.Breed}, {"birthdate", p.BirthDate}}},
		}
		f := bson.M{"_id": p.ID}
		_, err := collection.UpdateOne(
			context.Background(), f,
			u,
		)
		if err != nil {
			return nil, err
		}
		return p, nil

	}

	return nil, errors.New("did not get db instance, please check if the db is up and running")

}
