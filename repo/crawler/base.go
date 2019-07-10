package crawler

import (
	"crawler/config"
	"crawler/repo"
	"sync"
)

type RepoMongo repo.Mongo

var (
	instance *RepoMongo
	once     sync.Once
)

var cfg = config.GetConfig()

// New ..
func New() *RepoMongo {
	once.Do(func() {
		instance = &RepoMongo{
			Session:    cfg.Mongo.Get("crawler"),
			Collection: "Crawler",
		}
		session, db := instance.Session.GetCollection(instance.Collection)
		db.EnsureIndexKey("link")
		session.Close()
	})
	return instance
}
