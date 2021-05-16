package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"gosship/pkg/database"

	"github.com/dgraph-io/badger/v3"
	gossh "golang.org/x/crypto/ssh"
)

func GetHostSigner(db *database.Database) (gossh.Signer, error) {
	config, err := db.GetConfig()
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	}
	if err != badger.ErrKeyNotFound {
		return gossh.NewSignerFromSigner(ed25519.PrivateKey(config.PrivateKey))
	}

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	err = db.SetConfig(&database.ConfigEntry{PrivateKey: privateKey})
	if err != nil {
		return nil, err
	}
	return gossh.NewSignerFromSigner(privateKey)
}

func GetRoomConfig(db *database.Database) (*database.RoomConfigEntry, error) {
	config, err := db.GetRoomConfig()
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	}
	if err != badger.ErrKeyNotFound {
		return config, nil
	}

	config = &database.RoomConfigEntry{Rooms: map[string]string{
		"default": "",
	}}
	err = db.SetRoomConfig(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
