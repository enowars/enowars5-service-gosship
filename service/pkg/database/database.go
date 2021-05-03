package database

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var databasePath = "./db"

func init() {
	if dbp, ok := os.LookupEnv("DATABASE_PATH"); ok {
		databasePath = dbp
	}
}

const (
	TypeConfigEntry byte = iota
	TypeUserEntry
	TypeDirectMessageEntry
)

var configEntryKey = []byte("server-config")

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

func (db *Database) GetConfig() (*ConfigEntry, error) {
	db.log.Println("getting config...")
	var ce ConfigEntry
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(configEntryKey)
		if err != nil {
			return err
		}
		if item.UserMeta() != TypeConfigEntry {
			return fmt.Errorf("invalid config entry type")
		}
		return item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &ce)
		})
	})
	if err != nil {
		return nil, err
	}
	return &ce, nil
}

func (db *Database) SetConfig(ce *ConfigEntry) error {
	db.log.Println("updating config...")
	val, err := proto.Marshal(ce)
	if err != nil {
		return err
	}
	entry := badger.NewEntry(configEntryKey, val).WithMeta(TypeConfigEntry)
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (db *Database) AddOrUpdateUser(id string, u *UserEntry) error {
	db.log.Printf("adding/updating user with id: %s\n", id)
	val, err := proto.Marshal(u)
	if err != nil {
		return err
	}
	entry := badger.NewEntry([]byte(id), val).WithMeta(TypeUserEntry)
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (db *Database) FindUserByPredicate(predicate func(entry *UserEntry) bool) (string, *UserEntry, error) {
	var ue *UserEntry
	var usedId string
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.UserMeta() != TypeUserEntry {
				continue
			}
			var tmp UserEntry
			err := it.Item().Value(func(val []byte) error {
				return proto.Unmarshal(val, &tmp)
			})
			if err != nil {
				return err
			}
			if predicate(&tmp) {
				ue = &tmp
				usedId = string(item.Key())
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	return usedId, ue, nil
}

func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		db.log.Error(err)
	}
}
