package main

import (
	"doable-go/internal/config"
	"doable-go/internal/doable"
	"doable-go/internal/handlers/fiber/api"
	"doable-go/internal/handlers/fiber/frontend"
	log "doable-go/pkg/logger"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func init() {
	log.Init(config.Envs.LogLevel)
	log.Info("[Logger] Initialized")
	doable.Init()
}

func main() {
	//test()
	//TODO: Add lists and todos gathering on startup + watch for changes on files

	//-> Web server setup
	//Render engine
	engine := html.New("internal/frontend/templates", ".html")

	//Fiber app
	app := fiber.New(fiber.Config{
		Views:       engine,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	//Enable fiber logger only if log level is set to debug
	if strings.ToUpper(config.Envs.LogLevel) == "DEBUG" {
		app.Use(
			fiberLogger.New(fiberLogger.Config{
				Format: "[${ip}] ${status} - ${method} ${path} in ${latency}\n",
				Output: os.Stderr,
			}),
		)
	}

	//-> Frontend
	//Static assets
	app.Static("/css/", "./internal/frontend/assets/css/")
	app.Static("/js/", "./internal/frontend/assets/js/")
	app.Static("/fonts/", "./internal/frontend/assets/fonts/")

	//Base route
	app.Get("/", frontend.Index())

	//-> API
	//Routes
	app.Get("/api/lists/:id", api.GetTodoList())
	app.Get("/api/lists", api.GetTodoLists())
	app.Get("/api/todos/formatted", api.GetFormattedTodos())
	app.Get("/api/todos/:id/formatted", api.GetFormattedTodo())
	app.Put("/api/todos/:id/check", api.CheckTodo())
	app.Get("/api/todos/:id", api.GetTodo())
	app.Get("/api/todos", api.GetTodos())
	app.Post("/api/todos", api.CreateTodo())

	//Start server
	log.Fatal(app.Listen(":80"))
}

func test() {
	todos, err := doable.ReadTodos()
	if err != nil {
		log.Error(err)
	}

	lists, err := doable.ReadLists()
	if err != nil {
		log.Error(err)
	}

	for _, t := range todos {
		if !t.IsCompleted {
			if t.ListID != "" {
				fmt.Println(t.Title + " - " + t.GetListName(lists))
			} else {
				fmt.Println(t.Title)
			}
		}
	}
}
