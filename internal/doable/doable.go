package doable

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "doable-go/pkg/logger"

	"github.com/fsnotify/fsnotify"
)

var (
	cachedTodos []Todo     //Cached todos to avoid reading from files every time
	cachedLists []TodoList //Cached lists to avoid reading from files every time
)

const (
	STARTUP = 0 //Used for the loadFiles function to know if it is called on startup
)

func Init() {
	if _, err := os.Stat("sync"); os.IsNotExist(err) {
		os.Mkdir("sync", 0755)
	}
	if _, err := os.Stat("sync/todos"); os.IsNotExist(err) {
		os.Mkdir("sync/todos", 0755)
	}
	if _, err := os.Stat("sync/lists"); os.IsNotExist(err) {
		os.Mkdir("sync/lists", 0755)
	}

	//Watch for file changes
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		done := make(chan bool)

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Info("[FSNotify] Change detected, will rescan files", "event", event.Name)
						loadFiles()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					fmt.Println("Errore:", err)
				}
			}
		}()

		err = filepath.Walk("./sync", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		<-done
	}()

	//Load files on startup
	go loadFiles(STARTUP)
	log.Info("[Doable] Initialized")
}

// -> Todo related struct and functions
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

// GetTodos returns todos, if they are not cached it reads them from the files and caches them then returns them
func GetTodos() (todos []Todo, err error) {
	if cachedTodos != nil {
		return cachedTodos, nil
	}
	todos, err = ReadTodos()
	if err != nil {
		return nil, err
	}
	cachedTodos = todos
	return todos, nil
}

// ReadTodo reads a single todo from the .todo file in the todos directory, given the id, and returns it as a Todo struct
func ReadTodo(id string) (todo Todo, err error) {
	//Check if file exists
	if _, err = os.Stat("sync/todos/" + id + ".todo"); os.IsNotExist(err) {
		return
	}

	//Read file
	tempData, err := os.ReadFile("sync/todos/" + id + ".todo")
	if err != nil {
		return Todo{}, err
	}

	//Parse JSON
	var tempTodo Todo
	if err = json.Unmarshal(tempData, &tempTodo); err != nil {
		return
	}
	return tempTodo, nil
}

// ReadTodos reads all the todos from the .todo files in the todos directory and returns them as a slice of Todo with an error if any
func ReadTodos() (todos []Todo, err error) {
	//Get files
	files, err := os.ReadDir("sync/todos")
	if err != nil {
		return nil, err
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

// SaveTodo saves a todo to a .todo file in the todos directory, given the Todo struct
func SaveTodo(todo Todo) error {
	//Parse JSON
	data, err := json.MarshalIndent(todo, "", "\t")
	if err != nil {
		return err
	}

	//Write file
	err = os.WriteFile("sync/todos/"+todo.ID+".todo", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// -> Todo list related struct and functions
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

// GetLists returns lists, if they are not cached it reads them from the files and caches them then returns them
func GetLists() (lists []TodoList, err error) {
	if cachedLists != nil {
		return cachedLists, nil
	}
	lists, err = ReadLists()
	if err != nil {
		return nil, err
	}
	cachedLists = lists
	return lists, nil
}

// ReadList reads a single list from the .list file in the lists directory, given the id, and returns it as a TodoList struct
func ReadList(id string) (list TodoList, err error) {
	//Check if file exists
	if _, err = os.Stat("sync/lists/" + id + ".list"); os.IsNotExist(err) {
		return
	}

	//Read file
	tempData, err := os.ReadFile("sync/lists/" + id + ".list")
	if err != nil {
		return TodoList{}, err
	}

	//Parse JSON
	var tempList TodoList
	if err = json.Unmarshal(tempData, &tempList); err != nil {
		return
	}
	return tempList, nil
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

// loadFiles loads the todos and lists from the files
func loadFiles(mode ...int) {
	var modeString string
	if len(mode) > 0 && mode[0] == STARTUP {
		modeString = "startup"
	} else {
		modeString = "file change"
	}
	log.Info("[Doable] Caching (file load) on " + modeString + " started")

	var err error
	cachedLists, err = ReadLists()
	if err != nil {
		log.Error("[Doable] Error reading lists on " + modeString + " -> " + err.Error())
	}
	cachedTodos, err = ReadTodos()
	if err != nil {
		log.Error("[Doable] Error reading todos on " + modeString + " -> " + err.Error())
	}

	log.Info("[Doable] Caching (file load) on " + modeString + " done")
}
