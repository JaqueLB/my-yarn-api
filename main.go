package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]Yarn)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("my-yarn-api").Collection("yarns")
}

type Yarn struct {
	ID          string    `json:"_id`
	Color       *Color    `json:"color"`
	Brand       string    `json:"brand"`
	Name        string    `json:"name"`
	KnitNeedle  *Hook     `json:"knit_needle"`
	CrochetHook *Hook     `json:"crochet_hook"`
	Tex         int       `json:"tex"`
	Length      int       `json:"length"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Color struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Hook struct {
	Sizes []float32 `json:"sizes"`
}

func createYarn(yarn *Yarn) error {
	_, err := collection.InsertOne(ctx, yarn)
	return err
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	yarns := r.Group("/yarns")

	yarns.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"all": db})
	})

	yarns.GET("/:name", func(c *gin.Context) {
		var result *Yarn
		filter := bson.D{{"name", c.Params.ByName("name")}}
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"search": filter, "status": "no value"})
		}
		c.JSON(http.StatusNotFound, gin.H{"search": c.Params, "value": result})
	})

	yarns.PATCH("/:uuid", func(c *gin.Context) {
		uuid := c.Params.ByName("uuid")
		var body *Yarn
		_, ok := db[uuid]
		if ok && c.Bind(&body) == nil {
			db[uuid] = *body
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	yarns.POST("/", func(c *gin.Context) {
		var yarn *Yarn
		err := c.Bind(&yarn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		yarn.ID = uuid.NewString()
		yarn.CreatedAt = time.Now()
		yarn.UpdatedAt = time.Now()

		err = createYarn(yarn)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	yarns.DELETE("/:uuid", func(c *gin.Context) {
		uuid := c.Params.ByName("uuid")
		delete(db, uuid)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
