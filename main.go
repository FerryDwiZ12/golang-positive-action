package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

// GetHandler handles GET requests
func getHandler(c *fiber.Ctx, db *sql.DB) error {
	var todos []fiber.Map

	// Fetch and display todos
	rows, err := db.Query(`SELECT action, total_point FROM "public"."actions"`)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).SendString("An error occurred while retrieving data")
	}

	for rows.Next() {
		var action string
		var totalPoint int
		err := rows.Scan(&action, &totalPoint)
		if err != nil {
			log.Fatal(err)
			return c.Status(fiber.StatusInternalServerError).SendString("An error occurred while scanning data")
		}
		todos = append(todos, fiber.Map{
			"Action":     action,
			"TotalPoint": totalPoint,
		})
	}

	return c.Render("index", fiber.Map{
		"Title": "Positive Action",
		"Todos": todos,
	})
}

// PostHandler handles POST requests
func postHandler(c *fiber.Ctx, db *sql.DB) error {
	// Process the form submission
	action := c.FormValue("action")

	// Generate a random totalPoint between 1 and 10
	totalPoint := rand.Intn(10) + 1

	createdAt := time.Now()

	// SQL statement for insertion
	sqlStatement := `
		INSERT INTO "public"."actions" (action, total_point, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	// Execute the SQL statement
	var id int
	err := db.QueryRow(sqlStatement, action, totalPoint, createdAt).Scan(&id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("An error occurred while inserting data")
	}

	// Redirect to the home page after successful submission
	return c.Redirect("/")
}


func main() {
	connStr := "postgresql://postgres:root@localhost:5432/db_positive_action?sslmode=disable"
	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// GET handler
	app.Get("/", func(c *fiber.Ctx) error {
		return getHandler(c, db)
	})

	// POST handler
	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
