package main

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID      		int    `json:"id"`
	Item    		string `json:"item"`
	Completed 	bool   `json:"completed"`
}

var (
	todos = []todo{}
	nextID   int = 1
	idMutex sync.Mutex
)

func generateID() int {
	idMutex.Lock()
	defer idMutex.Unlock()
	id := nextID
	nextID++
	return id
}

func getTodoById(s string) (*todo, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return nil, errors.New("invalid todo id")
	}

	for i, t := range todos {
		if t.ID == id { 
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}

func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, todo)
}

func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}

func createTodo(context *gin.Context) {
	var todo todo
	if err := context.ShouldBindJSON(&todo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if todo.Item == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Item is required"})
		return
	}

	todo.ID = generateID()
	todos = append(todos, todo)
	context.IndentedJSON(http.StatusCreated, todo)
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	todo.Completed = !todo.Completed

	context.IndentedJSON(http.StatusOK, todo)
}

func updateTodo(context *gin.Context) {
	id := context.Param("id")
	todoLocal, err := getTodoById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var todoData todo
	if err := context.ShouldBindJSON(&todoData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	todoLocal.Item = todoData.Item
	todoLocal.Completed = todoData.Completed

	context.IndentedJSON(http.StatusOK, todoLocal)
}


func deleteTodoById(s string) (*todo, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return nil, errors.New("invalid todo id")
	}

	for i, t := range todos {
		if t.ID == id {
			deletedTodo := t
			todos = append(todos[:i], todos[i+1:]...)
			return &deletedTodo, nil
		}
	}
	return nil, errors.New("todo not found")
}

func deleteTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := deleteTodoById(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, todo)
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.POST("/todos", createTodo)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run("localhost:9191")
}
