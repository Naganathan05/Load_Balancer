package main

import (
	redisDB "Load_Balancer_Server/services"
	"fmt"
	// "log"
	// "github.com/gofiber/fiber/v2"
)

func main() {
	// app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, Fiber!")
	// })

	// log.Fatal(app.Listen(":8080"))
	fmt.Print("Hello")
	redisDB.RedisClient()
}
