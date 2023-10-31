package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/raphaelmb/go-hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDB) tearDown(t *testing.T) {
	dbName := os.Getenv(db.MongoDBNameEnvName)
	if err := tdb.client.Database(dbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}
	dbURI := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	return &testDB{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
