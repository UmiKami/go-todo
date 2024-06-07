package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

func main() {
	// Code here
	fmt.Println("Hello World!")
	app := fiber.New()

	api := app.Group("/api")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	} 

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello GO!"})
	})

	api.Get("/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	api.Post("/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Task == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Please provide a task to do."})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(fiber.Map{"msg": "Todo has been created!", "todo": *todo})
	})

	api.Put("/todos/:id", func(c *fiber.Ctx) error { 
		id := c.Params("id")

		var body map[string]interface{}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				modified := false

				if completed, exists := body["completed"]; exists {
					todos[i].Completed = completed.(bool)
					modified = true
				}

				if task, exists := body["task"]; exists {
					todos[i].Task = task.(string)
					modified = true
				}

				if modified {
					return c.Status(200).JSON(fiber.Map{"msg": "Todo with ID " + id + " has been updated!", "todo": todos[i]})
				}

				return c.Status(400).JSON(fiber.Map{"error": "No changes were detected."})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Invalid Todo ID"})
	})

	api.Delete("/todos/:id", func(c *fiber.Ctx) error { 
		id := c.Params("id")

		var body map[string]interface{}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				
				return c.Status(200).JSON(fiber.Map{"msg": "Todo with ID " + id + " has been deleted!"})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Invalid Todo ID"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
