package yarn

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Yarn struct {
	ID          primitive.ObjectID `bson:"_id"`
	Color       *Color             `json:"color"`
	Brand       string             `json:"brand"`
	Name        string             `json:"name"`
	KnitNeedle  *Hook              `json:"knit_needle"`
	CrochetHook *Hook              `json:"crochet_hook"`
	Tex         int                `json:"tex"`
	Length      int                `json:"length"`
	Weight      int                `json:"weight"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type Color struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Hook struct {
	Sizes []float32 `json:"sizes"`
}
