package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jaquelb/my-yarn-api/pkg/health"
	"go.mongodb.org/mongo-driver/bson"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", health.GetPing)

	yarns := r.Group("/yarns")

	yarns.GET("/", func(c *gin.Context) {

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
