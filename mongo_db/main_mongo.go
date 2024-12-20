package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Item      string             `json:"item"`
	Completed bool               `json:"completed"`
}

var collection *mongo.Collection

func getTodoDB(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	filter := bson.M{"_id": objectId}

	var todoLocal Todo

	if err := collection.FindOne(context.Background(), filter).Decode(&todoLocal); err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ginContext.IndentedJSON(http.StatusOK, todoLocal)
}

func getTodosDB(ginContext *gin.Context) {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		for cursor.Next(context.Background()) {
			var todo Todo
			if err := cursor.Decode(&todo); err != nil {
				ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			todos = append(todos, todo)
		}
	} else {
		// If no documents are found, return an empty slice
		todos = []Todo{}
	}

	ginContext.IndentedJSON(http.StatusOK, todos)
}

func createTodoDB(ginContext *gin.Context) {
	var todo Todo
	if err := ginContext.ShouldBindJSON(&todo); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	if todo.Item == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Item is required"})
		return
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ginContext.IndentedJSON(http.StatusCreated, todo)
}

func toggleTodoStatusDB(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	filter := bson.M{"_id": objectId}
	var todo Todo

	if err := collection.FindOne(context.Background(), filter).Decode(&todo); err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	update := bson.M{"$set": bson.M{"completed": !todo.Completed}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todo.Completed = !todo.Completed

	ginContext.IndentedJSON(http.StatusOK, todo)
}

func updateTodoDB(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	filter := bson.M{"_id": objectId}
	var todoLocal Todo

	if err := collection.FindOne(context.Background(), filter).Decode(&todoLocal); err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var todoData Todo
	if err := ginContext.ShouldBindJSON(&todoData); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	update := bson.M{"$set": bson.M{"item": todoData.Item}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	todoLocal.Item = todoData.Item
	ginContext.IndentedJSON(http.StatusOK, todoLocal)
}

func deleteTodoDB(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	filter := bson.M{"_id": objectId}

	var todoLocal Todo

	if err := collection.FindOne(context.Background(), filter).Decode(&todoLocal); err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ginContext.IndentedJSON(http.StatusOK, todoLocal)
}

func main() {
	MONGODB_URI := "mongodb+srv://qm0uFsC65I2ZUZCt:qm0uFsC65I2ZUZCt@cluster0.t3xsp.mongodb.net/golang_db?retryWrites=true&w=majority&appName=Cluster0"
	clientOption := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection = client.Database("golang_db").Collection("todos")

	router := gin.Default()
	router.GET("/todos", getTodosDB)
	router.POST("/todos", createTodoDB)
	router.GET("/todos/:id", getTodoDB)
	router.PATCH("/todos/:id", toggleTodoStatusDB)
	router.PUT("/todos/:id", updateTodoDB)
	router.DELETE("/todos/:id", deleteTodoDB)
	router.Run("localhost:9191")

	// Ensure we disconnected from DB
	defer client.Disconnect(context.Background())
}
