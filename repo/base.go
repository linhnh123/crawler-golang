package repo

import "crawler/config"

type Mongo struct {
	Session    *config.Mongo
	Collection string
}
