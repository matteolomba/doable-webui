// -> Fiber handlers for the API
package api

import (
	"doable-go/internal/doable"
	log "doable-go/pkg/logger"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler for /api/todos
func GetTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		todos, err := doable.GetTodos()
		if err != nil {
			return fiber.NewError(500, "Error while getting todos. Error: "+err.Error())
		}
		log.Info("[API] Todos requested")
		return c.JSON(todos)
	}
}

// Handler for /api/lists
func GetTodoLists() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lists, err := doable.GetLists()
		if err != nil {
			return fiber.NewError(500, "Error while getting lists. Error: "+err.Error())
		}

		log.Info("[API] Todo lists requested")
		return c.JSON(lists)
	}
}

// Handler for /api/todos/formatted
func GetFormattedTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		todos, err := doable.GetTodos()
		if err != nil {
			return fiber.NewError(500, "Error while getting todos for formatted todos. Error: "+err.Error())
		}

		lists, err := doable.GetLists()
		if err != nil {
			return fiber.NewError(500, "Error while getting lists for formatted todos. Error: "+err.Error())
		}

		// Format only completed todos
		var formattedTodos []doable.Todo
		for _, t := range todos {
			if !t.IsCompleted {
				if t.ListID != "" {
					t.ListID = t.GetListName(lists)
				}
				formattedTodos = append(formattedTodos, t)
			}
		}

		log.Info("[API] Requested formatted todos")
		return c.JSON(formattedTodos)
	}
}

// Handler for /api/todos/check
func CheckTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.FormValue("id")
		if id == "" {
			log.Info("[API] Bad request: No id provided")
			return fiber.NewError(400, "No id provided")
		}

		todo, err := doable.ReadTodo(id)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Bad request: Todo with id " + id + " does not exist")
				return fiber.NewError(400, "Todo with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while reading todo with id "+id, "error", err)
				return fiber.NewError(500, "Error while reading todo with id "+id)
			}
		}
		if !todo.IsCompleted {
			todo.IsCompleted = true
			todo.LastModified = time.Now().Format("2006-01-02T15:04:05.000")

			// Save the todo
			err := doable.SaveTodo(todo)
			if err != nil {
				return fiber.NewError(500, "Error while saving todo")
			}

			log.Info("[API] Todo \"" + todo.Title + "\" (" + todo.ID + ") checked as completed")
			return c.SendString("Todo \"" + todo.Title + "\" (" + todo.ID + ") checked as completed")
		} else {
			log.Info("[API] Bad request: Todo \"" + todo.Title + "\" (" + todo.ID + ") is already completed")
			return fiber.NewError(400, "Todo \""+todo.Title+"\" ("+todo.ID+") is already completed")
		}
	}
}
