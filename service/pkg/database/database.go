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
	TypeRoomEntry
	TypeIndexEntry
)

const (
	configEntryKey = "server-config"
)

var entryToPrefix = map[byte]string{
	TypeConfigEntry:  "config",
	TypeUserEntry:    "user",
	TypeMessageEntry: "message",
	TypeRoomEntry:    "room",
	TypeIndexEntry:   "index",
}

type Index string

const (
	IndexUserName          = Index("user/name")
	IndexUserFingerprint   = Index("user/fingerprint")
	IndexDirectMessageUser = Index("dm/user")
)

var DefaultRoom = &RoomEntry{}

func getKeyWithPrefix(meta byte, keyParts ...string) []byte {
	return []byte(fmt.Sprintf("%s/%s", entryToPrefix[meta], strings.Join(keyParts, "/")))
}

func stripKeyPrefix(meta byte, key []byte) string {
	return strings.Replace(string(key), entryToPrefix[meta]+"/", "", 1)
}

func getIndexKey(index Index, key string) []byte {
	return getKeyWithPrefix(TypeIndexEntry, string(index), key)
}

type Database struct {
	log          *logrus.Logger
	db           *badger.DB
	gcTickerStop chan struct{}
}

func NewDatabase(log *logrus.Logger) (*Database, error) {
	opts := badger.DefaultOptions(databasePath).
		WithLogger(log).
		WithLoggingLevel(badger.WARNING)
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

func createIndexEntry(index Index, key string, ttl time.Duration, ref []byte) *badger.Entry {
	return badger.NewEntry(getIndexKey(index, key), ref).WithMeta(TypeIndexEntry).WithTTL(ttl)
}

func createNewEntries(meta byte, id string, msg proto.Message) ([]*badger.Entry, error) {
	val, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	entry := badger.NewEntry(getKeyWithPrefix(meta, id), val).
		WithMeta(meta).
		WithDiscard()

	entries := []*badger.Entry{entry}
	switch meta {
	case TypeMessageEntry:
		// messages expire after 30 minutes
		ttl := 30 * time.Minute
		entry.WithTTL(ttl)
		m := msg.(*MessageEntry)
		if m.Type == MessageType_DIRECT {
			entries = append(entries,
				createIndexEntry(IndexDirectMessageUser, m.From+"/"+id[3:], ttl, entry.Key),
				createIndexEntry(IndexDirectMessageUser, m.To+"/"+id[3:], ttl, entry.Key))
		}
	case TypeUserEntry:
		// users expire after 1 hour after the last login
		ttl := 1 * time.Hour
		entry.WithTTL(ttl)
		u := msg.(*UserEntry)
		entries = append(entries,
			createIndexEntry(IndexUserName, u.Name, ttl, entry.Key),
			createIndexEntry(IndexUserFingerprint, u.Fingerprint, ttl, entry.Key))
	case TypeRoomEntry:
		// rooms expire after 1 hour
		entry.WithTTL(1 * time.Hour)
	}
	return entries, nil
}

func (db *Database) addNewEntry(meta byte, id string, msg proto.Message) error {
	entries, err := createNewEntries(meta, id, msg)
	if err != nil {
		return err
	}

	return db.db.Update(func(txn *badger.Txn) error {
		for _, e := range entries {
			if err := txn.SetEntry(e); err != nil {
				return err
			}
		}
		return nil
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

func (db *Database) GetRoom(room string) (*RoomEntry, error) {
	if room == "default" {
		return DefaultRoom, nil
	}
	var re RoomEntry
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getKeyWithPrefix(TypeRoomEntry, room))
		if err != nil {
			return err
		}
		if item.UserMeta() != TypeRoomEntry {
			return fmt.Errorf("invalid room entry type")
		}
		return item.Value(func(val []byte) error {
			return proto.Unmarshal(val, &re)
		})
	})
	if err != nil {
		return nil, err
	}
	return &re, nil
}

func (db *Database) AddRoom(room string, re *RoomEntry) error {
	return db.addNewEntry(TypeRoomEntry, room, re)
}

func (db *Database) GetAllRooms(namesOnly bool) (map[string]*RoomEntry, error) {
	res := make(map[string]*RoomEntry)
	res["default"] = DefaultRoom
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getKeyWithPrefix(TypeRoomEntry)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if item.UserMeta() != TypeRoomEntry {
				continue
			}
			roomName := stripKeyPrefix(TypeRoomEntry, item.Key())
			if namesOnly {
				res[roomName] = nil
				continue
			}
			var tmp RoomEntry
			err := it.Item().Value(func(val []byte) error {
				return proto.Unmarshal(val, &tmp)
			})
			if err != nil {
				return err
			}
			res[roomName] = &tmp
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (db *Database) AddOrUpdateUser(id string, u *UserEntry) error {
	db.log.Printf("adding/updating user with id: %s", id)
	return db.addNewEntry(TypeUserEntry, id, u)
}

func (db *Database) FindUserByIndex(index Index, searchKey string) (string, *UserEntry, error) {
	var ue *UserEntry
	var usedId string
	err := db.db.View(func(txn *badger.Txn) error {
		indexItem, err := txn.Get(getIndexKey(index, searchKey))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		userKey, err := indexItem.ValueCopy(nil)
		if err != nil {
			return err
		}
		userItem, err := txn.Get(userKey)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		if userItem.UserMeta() != TypeUserEntry {
			return fmt.Errorf("invlid database record type")
		}
		var tmp UserEntry
		err = userItem.Value(func(val []byte) error {
			return proto.Unmarshal(val, &tmp)
		})
		if err != nil {
			return err
		}
		ue = &tmp
		usedId = stripKeyPrefix(TypeUserEntry, userKey)
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	return usedId, ue, nil
}

func (db *Database) AddMessageEntry(m *MessageEntry) error {
	prefix := "unknown"
	switch m.Type {
	case MessageType_PUBLIC, MessageType_ROOM_ANNOUNCEMENT:
		prefix = "room/" + m.Room
	case MessageType_DIRECT:
		prefix = "dm"
	case MessageType_ANNOUNCEMENT:
		prefix = "announcement"
	}
	id := fmt.Sprintf("%s/%s", prefix, uuid.NewString())
	return db.addNewEntry(TypeMessageEntry, id, m)
}

func (db *Database) RenameUser(id, newUsername string) error {
	userId, _, err := db.FindUserByIndex(IndexUserName, newUsername)
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
		err = txn.Delete(getIndexKey(IndexUserName, tmp.Name))
		if err != nil {
			return err
		}

		tmp.Name = newUsername

		entries, err := createNewEntries(TypeUserEntry, id, &tmp)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := txn.SetEntry(entry); err != nil {
				return err
			}
		}
		return nil
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

func (db *Database) GetRecentMessagesForRoom(room string, skipAnnouncements bool) (MessageEntries, error) {
	res := make(MessageEntries, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixes := [][]byte{getKeyWithPrefix(TypeMessageEntry, "room", room)}
		if !skipAnnouncements {
			prefixes = append(prefixes, getKeyWithPrefix(TypeMessageEntry, "announcement"))
		}
		for _, prefix := range prefixes {
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
				if (tmp.Type == MessageType_ROOM_ANNOUNCEMENT || tmp.Type == MessageType_PUBLIC) && tmp.Room != room {
					continue
				}
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

func (db *Database) GetRecentDirectMessagesForUser(selfId, uid string) (MessageEntries, error) {
	res := make(MessageEntries, 0)
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getIndexKey(IndexDirectMessageUser, selfId)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			indexItem := it.Item()
			messageKey, err := indexItem.ValueCopy(nil)
			if err != nil {
				return err
			}
			messageItem, err := txn.Get(messageKey)
			if err != nil {
				if err == badger.ErrKeyNotFound {
					// index pointed to expired message, let's ignore it
					continue
				}
				return err
			}
			if messageItem.UserMeta() != TypeMessageEntry {
				continue
			}
			var tmp MessageEntry
			err = messageItem.Value(func(val []byte) error {
				return proto.Unmarshal(val, &tmp)
			})
			if err != nil {
				return err
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
	id, _, err := db.FindUserByIndex(IndexUserName, username)
	if err != nil {
		return err
	}

	if id == "" {
		return fmt.Errorf("username not found in database")
	}

	err = db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := getIndexKey(IndexDirectMessageUser, id)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			indexItem := it.Item()
			messageKey, err := indexItem.ValueCopy(nil)
			if err != nil {
				return err
			}
			messageItem, err := txn.Get(messageKey)
			if err != nil {
				if err == badger.ErrKeyNotFound {
					// index pointed to expired message, let's ignore it
					continue
				}
				return err
			}
			if messageItem.UserMeta() != TypeMessageEntry {
				continue
			}
			err = messageItem.Value(func(val []byte) error {
				var me MessageEntry
				if err := proto.Unmarshal(val, &me); err != nil {
					return err
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

func (db *Database) GetAllUsers() (UserEntries, error) {
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
