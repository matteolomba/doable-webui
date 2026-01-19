package doable

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "doable-go/pkg/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
)

var (
	cachedTodos     []Todo     // Cached todos to avoid reading from files every time
	cachedLists     []TodoList // Cached lists to avoid reading from files every time
	mu              sync.RWMutex
	ErrTitleMissing error = errors.New("doable: title (required) of the todo is missing")
	ErrDoesNotExist error = errors.New("doable: the required todo does not exist")
)

const (
	STARTUP = 0 // Used for the loadFiles function to know if it is called on startup
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

	// Watch for file changes
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
					// Only reload on write/create/remove/rename - standard save is write+rename(sometimes)
					// With atomic writes, we often see Create (temp) -> Write (temp) -> Rename (to target) -> Chmod
					// fsnotify might catch the Rename/Create on the target directory.
					if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Rename == fsnotify.Rename {
						// Simple debounce or just reload safely
						// log.Info("[FSNotify] Change detected, will rescan files", "event", event.Name)
						// For now, reload. It grabs Lock inside.
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

	// Load files on startup
	go loadFiles(STARTUP)
	log.Info("[Doable] Initialized")
}

// atomicWriteFile writes data to a temp file and renames it to filename for atomic writes
func atomicWriteFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	tmpFile, err := os.CreateTemp(dir, "tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up if something fails before rename

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), filename)
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
	mu.RLock()
	// Return a copy to be safe, or just the slice (if caller doesn't modify)
	// For this app, returning the cached slice is okay as long as we don't modify elements in place without lock
	if cachedTodos != nil {
		// Make a copy to prevent concurrent map read/write if the caller acts on it while reload happens
		todosCopy := make([]Todo, len(cachedTodos))
		copy(todosCopy, cachedTodos)
		mu.RUnlock()
		return todosCopy, nil
	}
	mu.RUnlock()

	// If not cached, load (this shouldn't happen often after init, but strictly speaking loadFiles handles the caching)
	// But let's reuse ReadTodos which goes to disk
	return ReadTodos()
}

// ReadTodo reads a single todo from the .todo file in the todos directory, given the id, and returns it as a Todo struct
func ReadTodo(id string) (todo Todo, err error) {
	// Directly read file implies bypassing cache? Or fallback?
	// The original code reads from disk. Let's keep it that way but simpler.

	path := "sync/todos/" + id + ".todo"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}

	// Read file - atomic read not strictly necessary but good practice to handle partial writes?
	// Atomic write ensures we don't read partials.
	tempData, err := os.ReadFile(path)
	if err != nil {
		return Todo{}, err
	}

	// Parse JSON
	var tempTodo Todo
	if err = json.Unmarshal(tempData, &tempTodo); err != nil {
		return
	}
	return tempTodo, nil
}

// ReadTodos reads all the todos from the .todo files in the todos directory and returns them as a slice of Todo with an error if any
func ReadTodos() (todos []Todo, err error) {
	// Get files
	files, err := os.ReadDir("sync/todos")
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		// log.Warn("[Doable] No todos found in the directory")
	} else {
		for _, f := range files {
			if f.IsDir() || filepath.Ext(f.Name()) != ".todo" {
				continue
			}
			// Read file
			tempData, err := os.ReadFile("sync/todos/" + f.Name())
			if err != nil {
				return nil, err
			}

			// Parse JSON
			var tempTodo Todo
			err = json.Unmarshal(tempData, &tempTodo)
			if err != nil {
				// Log error and continue with other files
				log.Warn("[Doable] Skipping corrupt todo file", "file", f.Name(), "error", err)
				continue
			}
			todos = append(todos, tempTodo)
		}
	}
	return todos, nil
}

// FormatForCreation prepares the fields of the todo in creation with the given title, if the title is empty it returns an error, otherwise it creates the todo with a new UUID, sets the required fields and returns nil
func (t *Todo) FormatForCreation() (err error) {
	if t.Title == "" {
		return ErrTitleMissing
	}

	// Create ID, until it is unique
	isValid := false
	for !isValid {
		t.ID = uuid.New().String()
		if _, err = os.Stat("sync/todos/" + t.ID + ".todo"); err != nil { // error means probably not exist
			isValid = true
		}
	}

	// Set other fields
	t.IsCompleted = false
	t.IsImportant = false
	t.HasRecurred = false
	if t.ListID == "null" { // Fix for potential "null" string from JSON
		t.ListID = ""
	}
	t.CompletedDate = ""
	formattedTime := time.Now().Format("2006-01-02T15:04:05.000")
	t.CreationDate = formattedTime
	t.LastModified = formattedTime
	return nil
}

// GetListName returns the name of the list of the todo given the TodoList slice, if the list is not found it returns "Not found", if the todo has no list it returns "No list"
func (t *Todo) GetListName(li any) string {
	if t.ListID == "" {
		return "No list"
	}

	if list, ok := li.(TodoList); ok {
		return list.Name
	} else if lists, ok := li.([]TodoList); ok {
		for _, l := range lists {
			if l.ID == t.ListID {
				return l.Name
			}
		}
	}
	return "Not found"
}

// Save saves a todo to a .todo file in the todos directory, given the Todo struct
func (t *Todo) Save() (err error) {
	var data []byte
	// Parse JSON
	data, err = json.MarshalIndent(t, "", "\t")
	if err != nil {
		return err
	}

	// Write file atomically
	err = atomicWriteFile("sync/todos/"+t.ID+".todo", data)
	if err != nil {
		return err
	}

	// Update cache immediately to ensure consistency
	updateTodoCache(*t)
	return nil
}

// Delete deletes a todo from the .todo file in the todos directory, given the id
func (t *Todo) Delete() (err error) {
	// Check if file exists
	path := "sync/todos/" + t.ID + ".todo"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return ErrDoesNotExist
	}

	// Delete file
	if err = os.Remove(path); err != nil {
		return err
	}

	// Update cache
	deleteTodoCache(t.ID)
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
	Icon           string `json:"icon"` // New field for custom icon (SVG name)
}

// FormatForCreation formats a TodoList for creation (generates ID and timestamps)
func (tl *TodoList) FormatForCreation() error {
	if tl.Name == "" {
		return errors.New("list name is required")
	}

	tl.ID = uuid.New().String()
	now := time.Now().Format("2006-01-02T15:04:05.000")
	tl.CreationDate = now
	tl.LastModified = now
	tl.HadInitialSync = false

	return nil
}

// Save saves a TodoList to its .list file
func (tl *TodoList) Save() error {
	// Marshal to JSON
	data, err := json.MarshalIndent(tl, "", "\t")
	if err != nil {
		return err
	}

	// Write to file atomically
	err = atomicWriteFile("sync/lists/"+tl.ID+".list", data)
	if err != nil {
		return err
	}

	// Update cache immediately
	updateListCache(*tl)

	log.Info("[Doable] Saved list", "id", tl.ID, "name", tl.Name)
	return nil
}

// Delete deletes a TodoList file
func (tl *TodoList) Delete() error {
	// Delete the file
	err := os.Remove("sync/lists/" + tl.ID + ".list")
	if err != nil {
		return err
	}

	// Refresh the cache
	deleteListCache(tl.ID)

	log.Info("[Doable] Deleted list", "id", tl.ID, "name", tl.Name)
	return nil
}

// --- Cache Helpers ---

func updateTodoCache(t Todo) {
	mu.Lock()
	defer mu.Unlock()

	// Check if it exists to update, otherwise append
	for i, cached := range cachedTodos {
		if cached.ID == t.ID {
			cachedTodos[i] = t
			return
		}
	}
	cachedTodos = append(cachedTodos, t)
}

func deleteTodoCache(id string) {
	mu.Lock()
	defer mu.Unlock()

	for i, cached := range cachedTodos {
		if cached.ID == id {
			// Remove from slice
			cachedTodos = append(cachedTodos[:i], cachedTodos[i+1:]...)
			return
		}
	}
}

func updateListCache(l TodoList) {
	mu.Lock()
	defer mu.Unlock()

	for i, cached := range cachedLists {
		if cached.ID == l.ID {
			cachedLists[i] = l
			return
		}
	}
	cachedLists = append(cachedLists, l)
}

func deleteListCache(id string) {
	mu.Lock()
	defer mu.Unlock()

	for i, cached := range cachedLists {
		if cached.ID == id {
			cachedLists = append(cachedLists[:i], cachedLists[i+1:]...)
			return
		}
	}
}

// GetLists returns lists, if they are not cached it reads them from the files and caches them then returns them
func GetLists() (lists []TodoList, err error) {
	mu.RLock()
	if cachedLists != nil {
		listsCopy := make([]TodoList, len(cachedLists))
		copy(listsCopy, cachedLists)
		mu.RUnlock()
		return listsCopy, nil
	}
	mu.RUnlock()
	return ReadLists()
}

// ReadList reads a single list from the .list file in the lists directory, given the id, and returns it as a TodoList struct
func ReadList(id string) (list TodoList, err error) {
	// Check if file exists
	path := "sync/lists/" + id + ".list"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}

	// Read file
	tempData, err := os.ReadFile(path)
	if err != nil {
		return TodoList{}, err
	}

	// Parse JSON
	var tempList TodoList
	if err = json.Unmarshal(tempData, &tempList); err != nil {
		return
	}
	return tempList, nil
}

// ReadLists reads all the lists from the .list files in the lists directory and returns them as a slice of TodoList with an error if any
func ReadLists() (lists []TodoList, err error) {
	// Get files
	files, err := os.ReadDir("sync/lists")
	if err != nil {
		return
	}

	if len(files) == 0 {
		log.Warn("[Doable] No lists found in the directory")
	} else {
		for _, f := range files {
			if f.IsDir() || filepath.Ext(f.Name()) != ".list" {
				continue
			}
			// Read file
			tempData, err := os.ReadFile("sync/lists/" + f.Name())
			if err != nil {
				return nil, err
			}

			// Parse JSON
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
	mu.Lock()
	defer mu.Unlock()

	var modeString string
	if len(mode) > 0 && mode[0] == STARTUP {
		modeString = "startup"
	} else {
		modeString = "file change"
	}
	log.Debug("[Doable] Caching (file load) on " + modeString + " started")

	var err error
	cachedLists, err = ReadLists()
	if err != nil {
		log.Error("[Doable] Error reading lists on " + modeString + " -> " + err.Error())
	}
	cachedTodos, err = ReadTodos()
	if err != nil {
		log.Error("[Doable] Error reading todos on " + modeString + " -> " + err.Error())
	}

	log.Debug("[Doable] Caching (file load) on " + modeString + " done")
}
