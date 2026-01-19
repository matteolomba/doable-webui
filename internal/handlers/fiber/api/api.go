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

		// Validation
		if todo.Title == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Todo title is required")
		}
		if len(todo.Title) > 100 {
			return fiber.NewError(fiber.StatusBadRequest, "Todo title must be less than 100 characters")
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

// Handler for PUT /api/todos/:id
func UpdateTodo() func(c *fiber.Ctx) error {
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

		// Parse the update data
		var updateData map[string]interface{}
		if err := c.BodyParser(&updateData); err != nil {
			log.Error("[API] Error while parsing update data", "error", err)
			return fiber.NewError(fiber.StatusBadRequest, "Error while parsing update data")
		}

		// Update fields
		if title, ok := updateData["title"].(string); ok {
			if title == "" {
				return fiber.NewError(fiber.StatusBadRequest, "Todo title cannot be empty")
			}
			if len(title) > 100 {
				return fiber.NewError(fiber.StatusBadRequest, "Todo title must be less than 100 characters")
			}
			todo.Title = title
		}
		if description, ok := updateData["description"].(string); ok {
			todo.Description = description
		}
		if listID, ok := updateData["listId"].(string); ok {
			todo.ListID = listID
		}

		todo.LastModified = time.Now().Format("2006-01-02T15:04:05.000")

		// Save the todo
		err = todo.Save()
		if err != nil {
			log.Error("[API] Error while saving updated todo", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while saving updated todo")
		}

		log.Info("[API] Updated todo", "id", id)
		return c.JSON(todo)
	}
}

// ============================================================================
// LISTS HANDLERS
// ============================================================================

// Handler for POST /api/lists
func CreateList() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var list doable.TodoList
		if err := c.BodyParser(&list); err != nil {
			log.Error("[API] Error while parsing list from body", "error", err)
			return fiber.NewError(fiber.StatusBadRequest, "Error while parsing list from body")
		}

		// Validate required fields
		if list.Name == "" {
			log.Info("[API] Bad request: List name is required")
			return fiber.NewError(fiber.StatusBadRequest, "List name is required")
		}
		if len(list.Name) > 50 {
			return fiber.NewError(fiber.StatusBadRequest, "List name must be less than 50 characters")
		}

		// Format for creation (sets ID, timestamps)
		err := list.FormatForCreation()
		if err != nil {
			log.Error("[API] Error while formatting the list for creation", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while formatting the list for creation")
		}

		// Save the list
		err = list.Save()
		if err != nil {
			log.Error("[API] Error while saving the new list", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while saving the new list")
		}

		log.Info("[API] Created new list with name \"" + list.Name + "\" (" + list.ID + ")")
		return c.Status(fiber.StatusCreated).JSON(list)
	}
}

// Handler for PUT /api/lists/:id
func UpdateList() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			log.Info("[API] No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}

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

		// Parse the update data
		var updateData map[string]interface{}
		if err := c.BodyParser(&updateData); err != nil {
			log.Error("[API] Error while parsing update data", "error", err)
			return fiber.NewError(fiber.StatusBadRequest, "Error while parsing update data")
		}

		// Update fields
		if name, ok := updateData["name"].(string); ok && name != "" {
			if len(name) > 50 {
				return fiber.NewError(fiber.StatusBadRequest, "List name must be less than 50 characters")
			}
			list.Name = name
		}
		if color, ok := updateData["color"].([]interface{}); ok && len(color) == 3 {
			// Convert interface{} to int for RGB
			colorInts := make([]int, 3)
			for i, v := range color {
				if f, ok := v.(float64); ok {
					colorInts[i] = int(f)
				}
			}
			list.Color = colorInts
		}
		if icon, ok := updateData["icon"].(string); ok {
			list.Icon = icon
		}

		list.LastModified = time.Now().Format("2006-01-02T15:04:05.000")

		// Save the list
		err = list.Save()
		if err != nil {
			log.Error("[API] Error while saving updated list", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while saving updated list")
		}

		log.Info("[API] Updated list", "id", id, "name", list.Name)
		return c.JSON(list)
	}
}

// Handler for DELETE /api/lists/:id
func DeleteList() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			log.Info("[API] Bad request: No id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No id provided")
		}

		list, err := doable.ReadList(id)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Not found: List with id " + id + " does not exist")
				return fiber.NewError(fiber.StatusNotFound, "List with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while reading list with id "+id, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while reading list with id "+id)
			}
		}

		// Delete the list
		err = list.Delete()
		if err != nil {
			if os.IsNotExist(err) {
				log.Info("[API] Not found: List with id " + id + " does not exist")
				return fiber.NewError(fiber.StatusNotFound, "List with id "+id+" does not exist")
			} else {
				log.Error("[API] Error while deleting list with id "+id, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while deleting list with id "+id)
			}
		}

		log.Info("[API] Deleted list", "id", id, "name", list.Name)
		return c.Status(fiber.StatusOK).JSON(list)
	}
}

// Handler for PUT /api/lists/:id/todos/reassign
// Body: {"action": "delete|reassign-none|reassign", "newListId": "id-for-reassign"}
func BulkUpdateListTodos() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		listID := c.Params("id")
		if listID == "" {
			log.Info("[API] Bad request: No list id provided")
			return fiber.NewError(fiber.StatusBadRequest, "No list id provided")
		}

		var reqData map[string]interface{}
		if err := c.BodyParser(&reqData); err != nil {
			log.Error("[API] Error while parsing bulk update request", "error", err)
			return fiber.NewError(fiber.StatusBadRequest, "Error while parsing bulk update request")
		}

		action, ok := reqData["action"].(string)
		if !ok || action == "" {
			log.Info("[API] Bad request: action is required")
			return fiber.NewError(fiber.StatusBadRequest, "action is required")
		}

		todos, err := doable.GetTodos()
		if err != nil {
			log.Error("[API] Error while getting todos", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error while getting todos")
		}

		// Filter todos that belong to this list
		var todosInList []*doable.Todo
		for i := range todos {
			if todos[i].ListID == listID {
				todosInList = append(todosInList, &todos[i])
			}
		}

		if len(todosInList) == 0 {
			log.Info("[API] No todos found in list", "listID", listID)
			return c.JSON(fiber.Map{"message": "No todos found in list", "action": action, "count": 0})
		}

		switch action {
		case "delete":
			// Delete all todos in this list
			for _, todo := range todosInList {
				err := todo.Delete()
				if err != nil {
					log.Error("[API] Error while deleting todo", "id", todo.ID, "error", err)
					return fiber.NewError(fiber.StatusInternalServerError, "Error while deleting todo "+todo.ID)
				}
			}
			log.Info("[API] Deleted todos from list", "listID", listID, "count", len(todosInList))
			return c.JSON(fiber.Map{"message": "Deleted all todos", "action": action, "count": len(todosInList)})

		case "reassign-none":
			// Reassign all todos to no list (empty string)
			for _, todo := range todosInList {
				todo.ListID = ""
				todo.LastModified = time.Now().Format("2006-01-02T15:04:05.000")
				err := todo.Save()
				if err != nil {
					log.Error("[API] Error while saving todo", "id", todo.ID, "error", err)
					return fiber.NewError(fiber.StatusInternalServerError, "Error while saving todo "+todo.ID)
				}
			}
			log.Info("[API] Reassigned todos to no list", "listID", listID, "count", len(todosInList))
			return c.JSON(fiber.Map{"message": "Reassigned all todos to no list", "action": action, "count": len(todosInList)})

		case "reassign":
			// Reassign all todos to another list
			newListID, ok := reqData["newListId"].(string)
			if !ok || newListID == "" {
				log.Info("[API] Bad request: newListId is required for reassign action")
				return fiber.NewError(fiber.StatusBadRequest, "newListId is required for reassign action")
			}

			// Verify new list exists
			if _, err := doable.ReadList(newListID); err != nil {
				if os.IsNotExist(err) {
					log.Info("[API] Not found: List with id " + newListID + " does not exist")
					return fiber.NewError(fiber.StatusNotFound, "List with id "+newListID+" does not exist")
				}
				log.Error("[API] Error while reading list", "id", newListID, "error", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Error while reading list")
			}

			for _, todo := range todosInList {
				todo.ListID = newListID
				todo.LastModified = time.Now().Format("2006-01-02T15:04:05.000")
				err := todo.Save()
				if err != nil {
					log.Error("[API] Error while saving todo", "id", todo.ID, "error", err)
					return fiber.NewError(fiber.StatusInternalServerError, "Error while saving todo "+todo.ID)
				}
			}
			log.Info("[API] Reassigned todos to new list", "fromListID", listID, "toListID", newListID, "count", len(todosInList))
			return c.JSON(fiber.Map{"message": "Reassigned all todos to new list", "action": action, "newListId": newListID, "count": len(todosInList)})

		default:
			log.Info("[API] Bad request: Unknown action", "action", action)
			return fiber.NewError(fiber.StatusBadRequest, "Unknown action: "+action)
		}
	}
}
