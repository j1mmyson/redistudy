package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/j1mmyson/redistudy/db"
)

var (
	ListenAddr = "localhost:3000"
	RedisAddr  = "localhost:6379"
)

func main() {
	database, err := db.NewDatabase(RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	router := initRouter(database)
	router.Listen(ListenAddr)
}

func initRouter(database *db.Database) *fiber.App {
	app := fiber.New()

	app.Post("/points", func(c *fiber.Ctx) error {
		var userJson db.User
		if err := c.BodyParser(&userJson); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err := database.SaveUser(&userJson)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"user": userJson,
		})
	})

	app.Get("/points/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		user, err := database.GetUser(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.JSON(fiber.Map{
			"user": user,
		})

	})

	return app
}
