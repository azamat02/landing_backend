package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mborders/logmatic"
	"log"
	"os"
	"time"
)

var DB *pgxpool.Pool
var l = logmatic.NewLogger()

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
}

func HandleForm(c *fiber.Ctx) error {
	//Get data of form
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	id := 0

	err = DB.QueryRow(context.Background(), "insert into forms_data (name, organization, phone, email, send_date) values($1,$2,$3,$4,$5) returning id",
		data["name"],
		data["organization"],
		data["phone"],
		data["email"],
		time.Now()).Scan(&id)

	if err != nil {
		l.Error("ERROR: %s", err)
		return err
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Connect() {
	DB_HOST, _ := os.LookupEnv("DB_HOST")
	DB_USER, _ := os.LookupEnv("DB_USER")
	DB_PASS, _ := os.LookupEnv("DB_PASS")
	DB_NAME, _ := os.LookupEnv("DB_NAME")
	DB_PORT, _ := os.LookupEnv("DB_PORT")
	DB_SSL, _ := os.LookupEnv("DB_SSL")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME, DB_SSL)
	connConfig, err := pgxpool.ParseConfig(dbUrl)

	if err != nil {
		l.Error("%s", err)
	}

	connConfig.ConnConfig.PreferSimpleProtocol = true

	//conn, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	conn, err := pgxpool.Connect(context.Background(), dbUrl)

	DB = conn
	if err != nil {
		l.Error("Unable to connect to database: %v\n, err", err)
		os.Exit(1)
	}

}

func Setup(app *fiber.App) {
	//Form route
	app.Post("api/form", HandleForm)
}

func main() {
	app := fiber.New()

	Setup(app)

	app.Use(logger.New(logger.ConfigDefault))

	app.Get("/dashboard", monitor.New())

	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	Connect()

	log.Fatal(app.Listen(":4000"))
}
