// -> Fiber handlers for the API
package api

import (
	"doable-go/internal/doable"
	log "doable-go/pkg/logger"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler for GET /api/todos/:id
func GetTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id != "" {
			todo, err := doable.ReadTodo(id)
			if err != nil {
				if os.IsNotExist(err) {
					log.Info("[API] Todo with id " + id + " does not exist")
					return fiber.NewError(fiber.StatusNotFound, "Todo with id "+id+" does not exist")
				} else {
					log.Error("[API] Error while reading todo with id "+id, "error", err)
					return fiber.NewError(fiber.StatusInternalServerError, "Error while reading todo with id "+id)
				}
			}
			log.Info("[API] Single todo requested", "id", id)
			return c.JSON(todo)
		} else {
			log.Info("[API] No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}
	}
}

// Handler for GET /api/todos
func GetTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		todos, err := doable.GetTodos()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while getting todos. Specific error: "+err.Error())
		}
		log.Info("[API] Todos requested")
		return c.JSON(todos)
	}
}

// Handler for GET /api/lists/:id
func GetTodoList() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id != "" {
			list, err := doable.ReadList(id)
			if err != nil {
				if os.IsNotExist(err) {
					log.Info("[API] List with id " + id + " does not exist")
					return fiber.NewError(fiber.StatusNotFound, "List with id "+id+" does not exist")
				} else {
					log.Error("[API] Error while reading list with id "+id, "error", err)
					return fiber.NewError(fiber.StatusInternalServerError, "Error while reading list with id "+id)
				}
			}
			log.Info("[API] Single list requested", "id", id)
			return c.JSON(list)
		} else {
			log.Info("[API] No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}
	}
}

// Handler for GET /api/lists
func GetTodoLists() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lists, err := doable.GetLists()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while getting lists. Specific error: "+err.Error())
		}

		log.Info("[API] Todo lists requested")
		return c.JSON(lists)
	}
}

// Handler for GET /api/todos/:id/formatted
func GetFormattedTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id != "" {
			todo, err := doable.ReadTodo(id)
			if err != nil {
				if os.IsNotExist(err) {
					log.Info("[API] Todo with id " + id + " does not exist")
					return fiber.NewError(fiber.StatusNotFound, "Todo with id "+id+" does not exist")
				} else {
					log.Error("[API] Error while reading todo with id "+id, "error", err)
					return fiber.NewError(fiber.StatusInternalServerError, "Error while reading todo with id "+id)
				}
			}

			list, err := doable.ReadList(todo.ListID)
			if err != nil && !os.IsNotExist(err) {
				log.Error("[API] Error while reading list with id "+todo.ListID, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while reading list with id "+todo.ListID)
			}

			if todo.ListID != "" {
				todo.ListID = todo.GetListName(list)
			}

			log.Info("[API] Requested formatted todo", "id", id)
			return c.JSON(todo)
		} else {
			log.Info("[API] No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}
	}
}

// Handler for GET /api/todos/formatted
func GetFormattedTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		todos, err := doable.GetTodos()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while getting todos for formatted todos. Specific error: "+err.Error())
		}

		lists, err := doable.GetLists()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error while getting lists for formatted todos. Specific error: "+err.Error())
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

// Handler for PUT /api/todos/:id/check
func CheckTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			log.Info("[API] No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}

		todo, err := doable.ReadTodo(id)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Todo with id " + id + " does not exist")
				return fiber.NewError(fiber.StatusNotFound, "Todo with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while reading todo with id "+id, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while reading todo with id "+id)
			}
		}
		if !todo.IsCompleted {
			todo.IsCompleted = true
			todo.LastModified = time.Now().Format("2006-01-02T15:04:05.000")

			// Save the todo
			err := todo.Save()
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Error while saving todo")
			}

			log.Info("[API] Todo checked as completed", "id", todo.ID)
			return c.Status(fiber.StatusNoContent).Send([]byte(""))
		} else {
			log.Info("[API] Todo is already completed", "id", todo.ID)
			return fiber.NewError(fiber.StatusBadRequest, "Todo with id "+todo.ID+" is already completed")
		}
	}
}

// Handler for PUT /api/todos/:id/uncheck
func UncheckTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			log.Info("[API] No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}

		todo, err := doable.ReadTodo(id)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Todo with id " + id + " does not exist")
				return fiber.NewError(fiber.StatusNotFound, "Todo with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while reading todo with id "+id, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while reading todo with id "+id)
			}
		}
		if todo.IsCompleted {
			todo.IsCompleted = false
			todo.LastModified = time.Now().Format("2006-01-02T15:04:05.000")

			// Save the todo
			err := todo.Save()
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Error while saving todo")
			}

			log.Info("[API] Todo unchecked", "id", todo.ID)
			return c.Status(fiber.StatusNoContent).Send([]byte(""))
		} else {
			log.Info("[API] Todo is already unchecked", "id", todo.ID)
			return fiber.NewError(fiber.StatusBadRequest, "Todo with id "+todo.ID+" is already unchecked")
		}
	}
}

// Handler for POST /api/todos
func CreateTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var todo doable.Todo
		if err := c.BodyParser(&todo); err != nil {
			log.Error("[API] Error while parsing todo from body", "error", err)
			return fiber.NewError(fiber.StatusBadRequest, "Error while parsing todo from body")
		}

		err := todo.FormatForCreation()
		if err != nil {
			log.Error("[API] Error while formatting the todo for creation", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while formatting the todo for creation")
		}

		err = todo.Save()
		if err != nil {
			log.Error("[API] Error while saving the new todo", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while saving the new todo")
		}

		log.Info("[API] Created new todo with title \"" + todo.Title + "\" (" + todo.ID + ")")
		return c.Status(fiber.StatusCreated).JSON(todo)
	}
}

// Handler for DELETE /api/todos/:id
func DeleteTodo() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			log.Info("[API] Bad request: No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}

		todo, err := doable.ReadTodo(id)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Not found: Todo with id " + id + " does not exist")
				return fiber.NewError(fiber.StatusNotFound, "Todo with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while reading todo with id "+id, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while reading todo with id "+id)
			}
		}

		err = todo.Delete()
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Not found: Todo with id " + id + " does not exist")
				return fiber.NewError(fiber.StatusNotFound, "Todo with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while deleting todo with id "+id, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while deleting todo with id "+id)
			}
		}

		// Return no content status
		log.Info("[API] Deleted todo", "id", id, "todo", todo)
		return c.Status(fiber.StatusOK).JSON(todo)
	}
}
