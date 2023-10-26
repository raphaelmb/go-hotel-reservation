package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/raphaelmb/go-hotel-reservation/api/middleware"
	"github.com/raphaelmb/go-hotel-reservation/db/fixtures"
	"github.com/raphaelmb/go-hotel-reservation/types"
)

func TestUserGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		otherUser = fixtures.AddUser(db.Store, "another", "user", false)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "hotel", "anywhere", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 5.5, hotel.ID)

		from           = time.Now()
		till           = from.AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		bookingHandler = NewBookingHandler(db.Store)

		app   = fiber.New()
		route = app.Group("/", middleware.JWTAuthentication(db.User))
	)

	t.Run("user should be able to get booking", func(t *testing.T) {
		route.Get("/:id", bookingHandler.HandleGetBooking)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
		req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 response but got %d", resp.StatusCode)
		}

		var bookingResp *types.Booking
		if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
			t.Fatal(err)
		}

		if booking.UserID != user.ID {
			t.Fatalf("expected user id %s but got %s", booking.UserID, user.ID)
		}
		if booking.ID != bookingResp.ID {
			t.Fatalf("expected booking id %s but got %s", booking.ID, bookingResp.ID)
		}
	})

	t.Run("different user should be able to get booking", func(t *testing.T) {
		route.Get("/:id", bookingHandler.HandleGetBooking)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
		req.Header.Add("X-Api-Token", CreateTokenFromUser(otherUser))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode == http.StatusOK {
			t.Fatalf("expected a non 200 response but got %d", resp.StatusCode)
		}
	})

}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.tearDown(t)

	var (
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "hotel", "anywhere", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 5.5, hotel.ID)

		from           = time.Now()
		till           = from.AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		bookingHandler = NewBookingHandler(db.Store)

		app   = fiber.New()
		admin = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
	)

	t.Run("admin should be able to get bookings", func(t *testing.T) {
		admin.Get("/", bookingHandler.HandleGetBookings)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 response but got %d", resp.StatusCode)
		}
		var bookings []*types.Booking
		if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
			t.Fatal(err)
		}
		if len(bookings) != 1 {
			t.Fatalf("expected 1 but got %d", len(bookings))
		}
		if booking.ID != bookings[0].ID {
			t.Fatalf("expected %s but got %s", booking.ID, bookings[0].ID)
		}
		if booking.UserID != bookings[0].UserID {
			t.Fatalf("expected %s but got %s", booking.UserID, bookings[0].UserID)
		}
	})

	t.Run("non admin should not be able to get bookings", func(t *testing.T) {
		admin.Get("/", bookingHandler.HandleGetBookings)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode == http.StatusOK {
			t.Fatalf("expected a non 200 response but got %d", resp.StatusCode)
		}
	})
}
