package database

import (
	"checker/pkg/checker"
	"checker/pkg/client"
	"encoding/json"
	"fmt"
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

func getPrefixedKey(prefix, key string) []byte {
	return []byte(fmt.Sprintf("%s/%s", prefix, key))
}

func (db *Database) put(prefix, key string, data []byte) error {
	db.log.Printf("put %s entry: %s", prefix, key)
	entry := badger.NewEntry(getPrefixedKey(prefix, key), data)
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (db *Database) get(prefix, key string) ([]byte, error) {
	db.log.Printf("get %s entry: %s", prefix, key)
	var val []byte
	err := db.db.View(func(txn *badger.Txn) error {
		get, err := txn.Get(getPrefixedKey(prefix, key))
		if err != nil {
			return err
		}
		val, err = get.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}

type TaskChainEntry struct {
	Type        string               `json:"type"`
	Variant     string               `json:"variant"`
	TaskMessage *checker.TaskMessage `json:"taskMessage"`
	UserA       *client.User         `json:"userA"`
	UserB       *client.User         `json:"userB"`
	Room        string               `json:"room"`
	Password    string               `json:"password"`
	Noise       string               `json:"noise"`
	Timestamp   time.Time            `json:"timestamp"`
}

func (db *Database) PutTaskChainEntry(fi *TaskChainEntry) error {
	fi.Timestamp = time.Now()
	data, err := json.Marshal(fi)
	if err != nil {
		return err
	}
	return db.put("task", fi.TaskMessage.TaskChainId, data)
}

func (db *Database) GetTaskChainEntry(taskChainId string) (*TaskChainEntry, error) {
	data, err := db.get("task", taskChainId)
	if err != nil {
		return nil, err
	}
	var fi TaskChainEntry
	if err := json.Unmarshal(data, &fi); err != nil {
		return nil, err
	}
	return &fi, nil
}

type TeamEntry struct {
	TeamId    uint64    `json:"teamId"`
	PublicKey string    `json:"publicKey"`
	Timestamp time.Time `json:"timestamp"`
}

func (db *Database) PutTeamEntry(te *TeamEntry) error {
	teamIdStr := fmt.Sprintf("%d", te.TeamId)
	te.Timestamp = time.Now()
	data, err := json.Marshal(te)
	if err != nil {
		return err
	}
	return db.put("team", teamIdStr, data)
}

func (db *Database) GetTeamEntry(teamId uint64) (*TeamEntry, error) {
	teamIdStr := fmt.Sprintf("%d", teamId)
	data, err := db.get("team", teamIdStr)
	if err != nil {
		return nil, err
	}
	var te TeamEntry
	if err := json.Unmarshal(data, &te); err != nil {
		return nil, err
	}
	return &te, nil
}

func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		db.log.Error(err)
	}
}
