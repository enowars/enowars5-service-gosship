package database

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
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
	TypeRoomConfigEntry
)

var (
	configEntryKey     = "server-config"
	roomConfigEntryKey = "room-config"
)

var entryToPrefix = map[byte]string{
	TypeConfigEntry:     "config",
	TypeUserEntry:       "user",
	TypeMessageEntry:    "message",
	TypeRoomConfigEntry: "config",
}

func getKeyWithPrefix(meta byte, keyParts ...string) []byte {
	return []byte(fmt.Sprintf("%s/%s", entryToPrefix[meta], strings.Join(keyParts, "/")))
}

func stripKeyPrefix(meta byte, key []byte) string {
	return strings.Replace(string(key), entryToPrefix[meta]+"/", "", 1)
}

type Database struct {
	log          *logrus.Logger
	db           *badger.DB
	gcTickerStop chan struct{}
}

func NewDatabase(log *logrus.Logger) (*Database, error) {
	opts := badger.DefaultOptions(databasePath).WithLogger(log).WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		log:          log,
		db:           db,
		gcTickerStop: make(chan struct{}, 1),
	}, nil
}

func (db *Database) addNewEntry(meta byte, id string, msg proto.Message) error {
	val, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	entry := badger.NewEntry(getKeyWithPrefix(meta, id), val).WithMeta(meta)

	switch meta {
	case TypeMessageEntry:
		// messages expire after 30 minutes
		entry.WithTTL(30 * time.Minute)
	case TypeUserEntry:
		// users expire after 2 hours after the last login
		entry.WithTTL(2 * time.Hour)
	}

	return db.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(entry)
	})
}

func (db *Database) GetConfig() (*ConfigEntry, error) {
	db.log.Println("getting config...")
	var ce ConfigEntry
	err := db.db.View(func(txn *badger.Txn) error {

		item, err := txn.Get(getKeyWithPrefix(TypeConfigEntry, configEntryKey))
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

func (db *Database) GetRoomConfig() (*RoomConfigEntry, error) {
	db.log.Println("getting room config...")
	var rce RoomConfigEntry
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getKeyWithPrefix(TypeRoomConfigEntry, roomConfigEntryKey))
		if err != nil {
			return err
		}
		if item.UserMeta() != TypeRoomConfigEntry {
			return fmt.Errorf("invalid config entry type")
		}
		return item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &rce)
		})
	})
	if err != nil {
		return nil, err
	}
	return &rce, nil
}

func (db *Database) SetRoomConfig(rce *RoomConfigEntry) error {
	db.log.Println("updating room config...")
	return db.addNewEntry(TypeRoomConfigEntry, roomConfigEntryKey, rce)
}

func (db *Database) UpdateRooms(rooms map[string]*RoomEntry) error {
	return db.SetRoomConfig(&RoomConfigEntry{Rooms: rooms})
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
		prefix := getKeyWithPrefix(TypeUserEntry)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
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
				usedId = stripKeyPrefix(TypeUserEntry, item.Key())
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
		item, err := txn.Get(getKeyWithPrefix(TypeUserEntry, id))
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
		return txn.SetEntry(badger.NewEntry(getKeyWithPrefix(TypeUserEntry, id), updatedVal).WithMeta(TypeUserEntry))
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
	pastMarker := time.Now().Add(-15 * time.Minute)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getKeyWithPrefix(TypeMessageEntry)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
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
	pastMarker := time.Now().Add(-15 * time.Minute)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getKeyWithPrefix(TypeMessageEntry)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
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
		item, err := txn.Get(getKeyWithPrefix(TypeUserEntry, id))
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

func (db *Database) DumpDirectMessages(username string, emit func(*MessageEntry) error) error {
	id, _, err := db.FindUserByPredicate(func(entry *UserEntry) bool {
		return entry.Name == username
	})
	if err != nil {
		return err
	}

	if id == "" {
		return fmt.Errorf("username not found in database")
	}

	err = db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getKeyWithPrefix(TypeMessageEntry)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if item.UserMeta() != TypeMessageEntry {
				continue
			}
			err := item.Value(func(val []byte) error {
				var me MessageEntry
				err := proto.Unmarshal(val, &me)
				if err != nil {
					return err
				}
				if me.Type != MessageType_DIRECT {
					return nil
				}
				if me.To != id && me.From != id {
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

type UserEntries []*UserEntry

func (m UserEntries) Len() int {
	return len(m)
}

func (m UserEntries) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}

func (m UserEntries) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (db *Database) DumpUsers() (UserEntries, error) {
	res := make(UserEntries, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getKeyWithPrefix(TypeUserEntry)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
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

func (db *Database) runGC() {
	for {
		err := db.db.RunValueLogGC(0.7)
		if err != nil {
			db.log.Error(err)
			break
		}
	}
}

func (db *Database) RunGarbageCollection() {
	d := 5 * time.Minute
	db.log.Printf("starting database garbage collection every %s...", d)
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			db.log.Debug("running database garbage collection...")
			db.runGC()
		case <-db.gcTickerStop:
			db.log.Println("stopped database garbage collection")
			return
		}
	}
}

func (db *Database) Close() {
	db.gcTickerStop <- struct{}{}
	db.log.Println("closing database...")
	err := db.db.Close()
	if err != nil {
		db.log.Error(err)
	}
}
