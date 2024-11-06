// -> Fiber handlers for the API
package api

import (
	"doable-go/internal/doable"
	log "doable-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Handler for /api/todos/get
func GetTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		todos, err := doable.ReadTodos()
		if err != nil {
			return fiber.NewError(500, "Error while getting todos")
		}
		log.Info("[API] Requested todos")
		return c.JSON(todos)
	}
}

// Handler for /api/lists/get
func GetTodoLists() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lists, err := doable.ReadLists()
		if err != nil {
			return fiber.NewError(500, "Error while getting lists")
		}

		log.Info("[API] Todo lists requested")
		return c.JSON(lists)
	}
}

// Handler for /api/todos/get/formatted
func GetFormattedTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		todos, err := doable.ReadTodos()
		if err != nil {
			return fiber.NewError(500, "Error while getting todos for formatted todos")
		}

		lists, err := doable.ReadLists()
		if err != nil {
			return fiber.NewError(500, "Error while getting lists for formatted todos")
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
