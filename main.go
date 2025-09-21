package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/clem-kay/mini-trello/config"
	"github.com/clem-kay/mini-trello/models"
	"github.com/clem-kay/mini-trello/routes"
	"github.com/clem-kay/mini-trello/utils"
)

func main() {
	setup()
	log.Println("Mini Trello application started successfully!")

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetEnv("ALLOW_ORIGINS", "*"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.RegisterAuthorRoutes(app)

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",

		ReadinessProbe: func(c *fiber.Ctx) bool {
			return config.DB != nil && config.DB.Exec("SELECT 1").Error == nil
		},
		ReadinessEndpoint: "/ready",
	}))

	port := utils.GetEnv("PORT", "3000")

	log.Println("Server is running on port " + port)
	log.Fatal(app.Listen(":" + port)) // Use log.Fatal to catch startup errors
}

func setup() {
	log.Println("Initializing the mini trello application...")
	log.Println("Setting up database connections...")

	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	log.Println("Database connection established successfully.")

	config.DB.AutoMigrate(&models.User{}, &models.Board{})
	log.Println("Database migrated successfully.")
}
