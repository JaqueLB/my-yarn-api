package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var db = make(map[string]Yarn)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017/")
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
	ID          primitive.ObjectID `bson:"_id"`
	Color       *Color             `json:"color, omitempty"`
	Brand       string             `json:"brand, omitempty"`
	Name        string             `json:"name, omitempty"`
	KnitNeedle  *Hook              `json:"knit_needle, omitempty"`
	CrochetHook *Hook              `json:"crochet_hook, omitempty"`
	Tex         int                `json:"tex, omitempty"`
	Length      int                `json:"length, omitempty"`
	Weight      int                `json:"weight, omitempty"`
	CreatedAt   time.Time          `json:"created_at, omitempty"`
	UpdatedAt   time.Time          `json:"updated_at, omitempty"`
}

type Color struct {
	Name string `json:"name, omitempty"`
	Code string `json:"code, omitempty"`
}

type Hook struct {
	Sizes []float32 `json:"sizes, omitempty"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	yarns := r.Group("/yarns")

	yarns.GET("/", func(c *gin.Context) {
		var results []*Yarn
		col, err := collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		for col.Next(context.TODO()) {
			var elem Yarn
			err := col.Decode(&elem)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			results = append(results, &elem)
		}
		if err := col.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		// Close the cursor once finished
		col.Close(context.TODO())
		c.JSON(http.StatusOK, gin.H{"yarns": results})
	})

	yarns.GET("/:uuid", func(c *gin.Context) {
		var result *Yarn
		filter := bson.D{{"id", c.Params.ByName("uuid")}}
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"search": filter, "status": "no value"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"search": c.Params, "value": result})
	})

	yarns.PATCH("/:uuid", func(c *gin.Context) {
		updatedAt := time.Now()
		filter := bson.D{{"id", c.Params.ByName("uuid")}}
		requestData, err := c.GetRawData()
		var m map[string]interface{}
		err = json.Unmarshal(requestData, &m)
		m["updated_at"] = updatedAt
		newData, err := json.Marshal(m)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		updateData := bson.D{}
		bson.UnmarshalExtJSON(newData, true, &updateData)
		update := bson.D{{"$set", updateData}}
		updateResponse, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"found":   updateResponse.MatchedCount,
			"updated": updateResponse.ModifiedCount,
		})
	})

	yarns.POST("/", func(c *gin.Context) {
		now := time.Now()
		requestData, err := c.GetRawData()
		var m map[string]interface{}
		err = json.Unmarshal(requestData, &m)
		m["created_at"] = now
		m["updated_at"] = now
		m["id"] = uuid.NewString()
		newData, err := json.Marshal(m)

		_, err = collection.InsertOne(ctx, newData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	yarns.DELETE("/:uuid", func(c *gin.Context) {
		_, err := collection.DeleteOne(context.TODO(), bson.D{{"id", c.Params.ByName("uuid")}})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
