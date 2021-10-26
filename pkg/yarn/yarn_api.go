package yarn

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *Yarn) GetYarn(c *gin.Context) {
	var results []*Yarn
	col, err := Db.Find(context.TODO(), bson.D{{}})
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
}
