package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/raphaelmb/go-hotel-reservation/api"
	"github.com/raphaelmb/go-hotel-reservation/db"
	"github.com/raphaelmb/go-hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
		ctx           = context.Background()
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("james token ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin token ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "hotel name", "Brazil", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 299.99, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)

	for i := 0; i < 100; i++ {
		fixtures.AddHotel(store, fmt.Sprintf("Hotel-%d", i), fmt.Sprintf("Localtion-%d", i), rand.Intn(5)+1, nil)
	}
}
