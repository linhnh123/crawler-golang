package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/globalsign/mgo"
)

type Mongo struct {
	Name       string        `mapstructure:"name"`
	Address    []string      `mapstructure:"address"`
	RepSetName string        `mapstructure:"repset_name"`
	DBName     string        `mapstructure:"dbname"`
	User       string        `mapstructure:"user"`
	Pass       string        `mapstructure:"pass"`
	Timeout    time.Duration `mapstructure:"timeout"`
	IsSSL      bool          `mapstructure:"is_ssl"`
	Source     string        `mapstructure:"source"`
	PoolLimit  int           `mapstructure:"pool_limit"`
	ReadPref   string        `mapstructure:"read_pref"`
	Session    *mgo.Session
}

type Mongos map[string]*Mongo

var (
	onceMongo      map[string]*sync.Once
	onceMongoMutex = sync.RWMutex{}
)

func init() {
	onceMongo = make(map[string]*sync.Once)
}

func mapReadPref(readPref string) mgo.Mode {
	switch readPref {
	case "Primary":
		return mgo.Primary
	case "PrimaryPreferred":
		return mgo.PrimaryPreferred
	case "Secondary":
		return mgo.Secondary
	case "SecondaryPreferred":
		return mgo.SecondaryPreferred
	case "Nearest":
		return mgo.Nearest
	case "Eventual":
		return mgo.Eventual
	case "Monotonic":
		return mgo.Monotonic
	case "Strong":
		return mgo.Strong
	default:
		return mgo.Primary
	}
}

func (adapters Mongos) Get(name string) (result *Mongo) {
	if adapter, ok := adapters[name]; ok {
		result = adapter
	} else {
		panic("Không tìm thấy config Mongo " + name)
	}
	return
}

func (config *Mongo) Init() {
	if onceMongo[config.Name] == nil {
		onceMongo[config.Name] = &sync.Once{}
	}
	onceMongo[config.Name].Do(func() {
		onceMongoMutex.Lock()
		log.Printf("[%s][%s] Mongo [connecting]\n", config.Name, config.Address)

		//create dialInfo
		dialInfo := &mgo.DialInfo{
			Addrs:          config.Address,
			ReplicaSetName: config.RepSetName,
			Database:       config.DBName,
			Username:       config.User,
			Password:       config.Pass,
			Timeout:        config.Timeout * time.Second,
			Source:         config.Source,
			PoolLimit:      config.PoolLimit,
		}

		if config.IsSSL {
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				conn, err := tls.Dial("tcp", addr.String(), &tls.Config{})
				return conn, err
			}
		}

		//connect to DB
		mySession, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			log.Printf("[%s][%s] Mongo [error]: %v\n", config.Name, config.Address, err)
			time.Sleep(1 * time.Second)
			onceMongo[config.Name] = &sync.Once{}
			onceMongoMutex.Unlock()
			config.Init()
			return
		}
		config.Session = mySession
		config.Session.SetMode(mapReadPref(config.ReadPref), true)
		log.Printf("[%s][%s] Mongo [connected]\n", config.Name, config.Address)
		onceMongoMutex.Unlock()
	})
}

// Collection ..
type Collection struct {
	*mgo.Collection
	mongo Mongo
}

// GetCollection func
func (config Mongo) GetCollection(name string) (session *mgo.Session, collection Collection) {
	if config.Session != nil {
		session = config.Session.Copy()
		collection = Collection{session.DB(config.DBName).C(name), config}
	} else {
		panic(fmt.Errorf("[%s] Chưa init Mongo", config.Name))
	}
	return
}
