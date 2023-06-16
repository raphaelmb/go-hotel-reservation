package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raphaelmb/go-hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "John",
		LastName:  "Rambo",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("User")
}
