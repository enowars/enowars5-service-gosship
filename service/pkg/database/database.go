package database

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/logrusorgru/aurora/v3"
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
	TypeMessageEntry
)

var configEntryKey = "server-config"

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

func (db *Database) addNewEntry(meta byte, id string, msg proto.Message) error {
	val, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	entry := badger.NewEntry([]byte(id), val).WithMeta(meta)
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (db *Database) GetConfig() (*ConfigEntry, error) {
	db.log.Println("getting config...")
	var ce ConfigEntry
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(configEntryKey))
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
	return db.addNewEntry(TypeConfigEntry, configEntryKey, ce)
}

func (db *Database) AddOrUpdateUser(id string, u *UserEntry) error {
	db.log.Printf("adding/updating user with id: %s", id)
	return db.addNewEntry(TypeUserEntry, id, u)
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

func (db *Database) UpdateUserFingerprint(username, fingerprint string) error {
	db.log.Printf("updating fingerprint for user: %s", username)
	userId, userEntry, err := db.FindUserByPredicate(func(entry *UserEntry) bool {
		return entry.Name == username
	})
	if err != nil {
		return err
	}
	if userId == "" {
		return fmt.Errorf("user %s not found", username)
	}
	userEntry.Fingerprint = fingerprint
	return db.addNewEntry(TypeUserEntry, userId, userEntry)
}

func (db *Database) AddMessageEntry(m *MessageEntry) error {
	return db.addNewEntry(TypeMessageEntry, uuid.NewString(), m)
}

func (db *Database) RenameUser(id, newUsername string) error {
	userId, _, err := db.FindUserByPredicate(func(entry *UserEntry) bool {
		return entry.Name == newUsername
	})
	if err != nil {
		return err
	}
	if userId != "" {
		return fmt.Errorf("username already taken")
	}
	return db.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		if item.UserMeta() != TypeUserEntry {
			return fmt.Errorf("invlid database record type")
		}
		var tmp UserEntry
		err = item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &tmp)
		})
		if err != nil {
			return err
		}
		tmp.Name = newUsername
		updatedVal, err := proto.Marshal(&tmp)
		if err != nil {
			return err
		}
		return txn.SetEntry(badger.NewEntry([]byte(id), updatedVal).WithMeta(TypeUserEntry))
	})
}

type MessageEntries []*MessageEntry

func (m MessageEntries) Len() int {
	return len(m)
}

func (m MessageEntries) Less(i, j int) bool {
	return m[i].Timestamp.AsTime().Before(m[j].Timestamp.AsTime())
}

func (m MessageEntries) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (db *Database) GetRecentMessagesForUserAndRoom(uid, room string) (MessageEntries, error) {
	res := make(MessageEntries, 0)
	pastMarker := time.Now().Add(-10 * time.Minute)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.UserMeta() != TypeMessageEntry {
				continue
			}
			var tmp MessageEntry
			err := it.Item().Value(func(val []byte) error {
				return proto.Unmarshal(val, &tmp)
			})
			if err != nil {
				return err
			}
			if tmp.Timestamp.AsTime().Before(pastMarker) {
				continue
			}
			if tmp.Type == MessageType_DIRECT && tmp.To != uid && tmp.From != uid {
				continue
			}
			if (tmp.Type == MessageType_ROOM_ANNOUNCEMENT || tmp.Type == MessageType_PUBLIC) && tmp.Room != room {
				continue
			}
			res = append(res, &tmp)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(res)
	return res, err
}

func (db *Database) GetRecentDirectMessagesForUser(selfId, uid string) (MessageEntries, error) {
	res := make(MessageEntries, 0)
	pastMarker := time.Now().Add(-24 * time.Hour)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.UserMeta() != TypeMessageEntry {
				continue
			}
			var tmp MessageEntry
			err := it.Item().Value(func(val []byte) error {
				return proto.Unmarshal(val, &tmp)
			})
			if err != nil {
				return err
			}
			if tmp.Type != MessageType_DIRECT {
				continue
			}
			if tmp.Timestamp.AsTime().Before(pastMarker) {
				continue
			}
			if (tmp.To == selfId && tmp.From == uid) || (tmp.From == selfId && tmp.To == uid) {
				res = append(res, &tmp)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(res)
	return res, err
}

func (db *Database) GetUserById(id string) (*UserEntry, error) {
	var ue UserEntry
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		if item.UserMeta() != TypeUserEntry {
			return fmt.Errorf("invalid database record type")
		}
		return item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &ue)
		})
	})
	if err != nil {
		return nil, err
	}
	return &ue, nil
}

func (db *Database) DumpToLog() {
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var logEntry string
				switch item.UserMeta() {
				case TypeConfigEntry:
					var ce ConfigEntry
					err := proto.Unmarshal(val, &ce)
					if err != nil {
						db.log.Error(err)
						return nil
					}
					logEntry = aurora.Sprintf("%s: %s", aurora.Red("TypeConfigEntry"), ce.String())
				case TypeUserEntry:
					var ue UserEntry
					err := proto.Unmarshal(val, &ue)
					if err != nil {
						db.log.Error(err)
						return nil
					}
					logEntry = aurora.Sprintf("%s: %s", aurora.Red("TypeUserEntry"), ue.String())
				case TypeMessageEntry:
					var me MessageEntry
					err := proto.Unmarshal(val, &me)
					if err != nil {
						db.log.Error(err)
						return nil
					}
					logEntry = aurora.Sprintf("%s: %s", aurora.Red("TypeMessageEntry"), me.String())
				default:
					logEntry = aurora.Sprintf(aurora.Red("unknown type"))
				}
				db.log.Println(logEntry)
				return nil
			})
			if err != nil {
				db.log.Error(err)
			}
		}
		return nil
	})
	if err != nil {
		db.log.Error(err)
	}
}

func (db *Database) DumpMessages(emit func(*MessageEntry) error) error {
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.UserMeta() != TypeMessageEntry {
				continue
			}
			err := item.Value(func(val []byte) error {
				var me MessageEntry
				err := proto.Unmarshal(val, &me)
				if err != nil {
					return nil
				}
				return emit(&me)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ResetExceptConfig() {
	err := db.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if item.UserMeta() == TypeConfigEntry {
				continue
			}
			err := txn.Delete(item.Key())
			if err != nil {
				db.log.Error(err)
			}
		}
		return nil
	})
	if err != nil {
		db.log.Error(err)
	}
}

func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		db.log.Error(err)
	}
}
