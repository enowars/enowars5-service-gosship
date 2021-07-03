package main

import (
	"fmt"
	"os"

	"checker/service/database"

	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	TypeConfigEntry byte = iota
	TypeUserEntry
	TypeMessageEntry
	TypeRoomConfigEntry
	TypeIndexEntry
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{})
}

func run(dbPath string) error {
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLogger(log).WithReadOnly(true))
	if err != nil {
		return err
	}
	defer db.Close()
	return db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			switch item.UserMeta() {
			case TypeIndexEntry:
				log.Infof("IDX(%s): %s", key, string(val))
			case TypeUserEntry:
				var tmp database.UserEntry
				err = proto.Unmarshal(val, &tmp)
				if err != nil {
					return fmt.Errorf("could not unmarshal %s: %w", key, err)
				}
				log.Infof("USR(%s): %s", key, tmp.String())
			case TypeMessageEntry:
				var tmp database.MessageEntry
				err = proto.Unmarshal(val, &tmp)
				if err != nil {
					return fmt.Errorf("could not unmarshal %s: %w", key, err)
				}
				log.Infof("MSG(%s): %s", key, tmp.String())
			case TypeRoomConfigEntry:
				var tmp database.RoomConfigEntry
				err = proto.Unmarshal(val, &tmp)
				if err != nil {
					return fmt.Errorf("could not unmarshal %s: %w", key, err)
				}
				log.Infof("CFG(%s): %s", key, tmp.String())
			case TypeConfigEntry:
				var tmp database.ConfigEntry
				err = proto.Unmarshal(val, &tmp)
				if err != nil {
					return fmt.Errorf("could not unmarshal %s: %w", key, err)
				}
				log.Infof("CFG(%s): %s", key, tmp.String())
			default:
				log.Warnf("unknown key: %s", key)
			}
		}
		return nil
	})
}

func main() {
	dbPath := "./db"
	if len(os.Args) >= 2 {
		dbPath = os.Args[1]
	}
	log.Infof("using db path: %s", dbPath)
	if err := run(dbPath); err != nil {
		log.Fatal(err)
	}
}
