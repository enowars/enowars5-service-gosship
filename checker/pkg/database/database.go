package database

import (
	"checker/pkg/checker"
	"checker/pkg/client"
	"encoding/json"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
)

var databasePath = "./db"

func init() {
	if dbp, ok := os.LookupEnv("DATABASE_PATH"); ok {
		databasePath = dbp
	}
}

type Database struct {
	log *logrus.Logger
	db  *badger.DB
}

func NewDatabase(log *logrus.Logger) (*Database, error) {
	opts := badger.DefaultOptions(databasePath).WithLogger(log).WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		log: log,
		db:  db,
	}, nil
}

type FlagInfo struct {
	Variant     string               `json:"variant"`
	TaskMessage *checker.TaskMessage `json:"taskMessage"`
	UserA       *client.User         `json:"userA"`
	UserB       *client.User         `json:"userB"`
	Room        string               `json:"room"`
	Password    string               `json:"password"`
	Timestamp   time.Time            `json:"timestamp"`
}

func (db *Database) PutFlagInfo(fi *FlagInfo) error {
	fi.Timestamp = time.Now()
	data, err := json.Marshal(fi)
	if err != nil {
		return err
	}
	db.log.Println(string(data)) //debug
	entry := badger.NewEntry([]byte(fi.TaskMessage.TaskChainId), data)
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (db *Database) GetFlagInfo(taskChainId string) (*FlagInfo, error) {
	var fi FlagInfo
	err := db.db.View(func(txn *badger.Txn) error {
		get, err := txn.Get([]byte(taskChainId))
		if err != nil {
			return err
		}
		return get.Value(func(val []byte) error {
			return json.Unmarshal(val, &fi)
		})
	})
	if err != nil {
		return nil, err
	}
	return &fi, nil
}

func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		db.log.Error(err)
	}
}
