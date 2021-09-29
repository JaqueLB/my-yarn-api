package main

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]Yarn)

type Yarn struct {
	Color       *Color `json:"color"`
	Brand       string `json:"brand"`
	Name        string `json:"name"`
	KnitNeedle  *Hook  `json:"knit_needle"`
	CrochetHook *Hook  `json:"crochet_hook"`
	Tex         int    `json:"tex"`
	Length      int    `json:"length"`
	Weight      int    `json:"weight"`
}

type Color struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Hook struct {
	Sizes []float32 `json:"sizes"`
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

	yarns.GET("/:uuid", func(c *gin.Context) {
		uuid := c.Params.ByName("uuid")
		value, ok := db[uuid]
		if ok {
			c.JSON(http.StatusOK, gin.H{"yarn": uuid, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"yarn": uuid, "status": "no value"})
		}
	})

	yarns.POST("/:uuid", func(c *gin.Context) {
		uuid := c.Params.ByName("uuid")
		var body *Yarn
		_, ok := db[uuid]
		if ok && c.Bind(&body) == nil {
			db[uuid] = *body
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	yarns.PUT("/", func(c *gin.Context) {
		var body *Yarn

		if c.Bind(&body) == nil {
			uuid := uuid.NewString()
			db[uuid] = *body
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
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
