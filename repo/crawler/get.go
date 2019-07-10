package crawler

import (
	"crawler/model"

	"github.com/globalsign/mgo/bson"
)

func (r RepoMongo) GetByLink(link string) (result model.Crawler, err error) {
	session, db := r.Session.GetCollection(r.Collection)
	err = db.Find(bson.M{"link": link}).One(&result)
	session.Close()
	return
}
