package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type Crawler struct {
	ID        bson.ObjectId `json:"ID,omitempty" bson:"_id,omitempty"`
	Title     string        `json:"title,omitempty" bson:"title,omitempty"`
	Content   string        `json:"content,omitempty" bson:"content,omitempty"`
	Link      string        `json:"link,omitempty" bson:"link,omitempty"`
	CreatedAt time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
