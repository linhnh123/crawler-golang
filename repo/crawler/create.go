package crawler

import (
	"crawler/model"
	"time"
)

func (r RepoMongo) Create(data model.Crawler) (err error) {
	session, db := r.Session.GetCollection(r.Collection)
	// override data
	data.CreatedAt = time.Now()

	err = db.Insert(data)
	session.Close()
	return
}
