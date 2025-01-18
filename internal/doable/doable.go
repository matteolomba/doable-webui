package doable

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2/log"
)

// Doable todo format struct to parse from .todo json files
type Todo struct {
	ID             string `json:"id"`
	CreationDate   string `json:"creationDate"`
	Title          string `json:"title"`
	IsImportant    bool   `json:"isImportant"`
	IsCompleted    bool   `json:"isCompleted"`
	LastModified   string `json:"lastModified"`
	HadInitialSync bool   `json:"hadInitialSync"`
	ListID         string `json:"listId"`
	HasRecurred    bool   `json:"hasRecurred"`
	CompletedDate  string `json:"completedDate"`
	Description    string `json:"description"`
}

// Doable todo list format struct to parse from .list json files
type TodoList struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HiveIcon struct {
		CodePoint          int    `json:"codePoint"`
		FontFamily         string `json:"fontFamily"`
		MatchTextDirection bool   `json:"matchTextDirection"`
	} `json:"hiveIcon"`
	SelectedHiveIcon struct {
		CodePoint          int    `json:"codePoint"`
		FontFamily         string `json:"fontFamily"`
		MatchTextDirection bool   `json:"matchTextDirection"`
	} `json:"selectedHiveIcon"`
	CreationDate   string `json:"creationDate"`
	LastModified   string `json:"lastModified"`
	HadInitialSync bool   `json:"hadInitialSync"`
	Color          []int  `json:"color"`
	IsHidden       bool   `json:"isHidden"`
}

// ReadTodos reads all the todos from the .todo files in the todos directory and returns them as a slice of Todo with an error if any
func ReadTodos() (todos []Todo, err error) {
	//Get files
	files, err := os.ReadDir("sync/todos")
	if err != nil {
		return
	}

	if len(files) == 0 {
		log.Warn("[Doable] No todos found in the directory")
	} else {
		for _, f := range files {
			//Read file
			tempData, err := os.ReadFile("sync/todos/" + f.Name())
			if err != nil {
				return nil, err
			}

			//Parse JSON
			var tempTodo Todo
			err = json.Unmarshal(tempData, &tempTodo)
			if err != nil {
				return nil, err
			}
			todos = append(todos, tempTodo)
		}
	}
	return todos, nil
}

// ReadLists reads all the lists from the .list files in the lists directory and returns them as a slice of TodoList with an error if any
func ReadLists() (lists []TodoList, err error) {
	//Get files
	files, err := os.ReadDir("sync/lists")
	if err != nil {
		return
	}

	if len(files) == 0 {
		log.Warn("[Doable] No lists found in the directory")
	} else {
		for _, f := range files {
			//Read file
			tempData, err := os.ReadFile("sync/lists/" + f.Name())
			if err != nil {
				return nil, err
			}

			//Parse JSON
			var tempList TodoList
			err = json.Unmarshal(tempData, &tempList)
			if err != nil {
				return nil, err
			}
			lists = append(lists, tempList)
		}
	}
	return lists, nil
}

// GetListName returns the name of the list of the todo given the TodoList slice, if the list is not found it returns "Not found", if the todo has no list it returns "No list"
func (t *Todo) GetListName(lists []TodoList) string {
	if t.ListID == "" {
		return "No list"
	}
	for _, l := range lists {
		if l.ID == t.ListID {
			return l.Name
		}
	}
	return "Not found"
}
