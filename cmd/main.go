package main

import (
	"os"

	"github.com/Innoractive/chippiphone/handler"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load configratinos
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Templating engine
	engine := html.New("./view", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	mainHandler := handler.MainHandler{}
	app.Get("/", mainHandler.Index)

	searchHandler := handler.SearchHandler{}
	app.Get("/search", searchHandler.Index)

	app.Listen(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))
}
