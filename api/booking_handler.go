package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raphaelmb/go-hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound()
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"cancelled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{Type: "msg", Msg: "updated"})
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceNotFound()
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return ErrResourceNotFound()
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnauthorized()
	}

	if booking.UserID != user.ID {
		return ErrUnauthorized()
	}

	return c.JSON(booking)
}
