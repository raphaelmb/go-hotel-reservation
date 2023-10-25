package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/raphaelmb/go-hotel-reservation/db/fixtures"
)

func TestGetBookins(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "hotel", "anywhere", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 5.5, hotel.ID)

	from := time.Now()
	till := from.AddDate(0, 0, 2)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
	fmt.Println(booking)
}
