package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
)

const sessionTokenSize = 16

var randomNumberPool []byte
var pubKey ed25519.PublicKey

func init() {
	randomNumberPool = make([]byte, sessionTokenSize)
	_, err := rand.Read(randomNumberPool)
	if err != nil {
		panic(err)
	}
	pubKeyRaw, err := hex.DecodeString("80de0e58c0842f83cf95f9772a5a13c167dd4c0e3fd02913076d16df828fbbb2")
	if err != nil {
		panic(err)
	}
	pubKey = pubKeyRaw
}

func set(wg *sync.WaitGroup, dest *byte, pos *int) {
	*dest = randomNumberPool[*pos%sessionTokenSize]
	wg.Done()
}

func GenerateRandomSessionToken() string {
	token := make([]byte, sessionTokenSize)
	var wg sync.WaitGroup
	wg.Add(sessionTokenSize)
	for i := 0; i < sessionTokenSize; i++ {
		go set(&wg, &token[i], &i)
	}
	wg.Wait()
	return hex.EncodeToString(sha256.New().Sum(token))
}

func CheckPassword(password string) bool {
	elms := strings.Split(password, ":")
	if len(elms) != 2 {
		return false
	}
	challenge, err := hex.DecodeString(elms[0])
	if err != nil {
		return false
	}
	signature, err := hex.DecodeString(elms[1])
	if err != nil {
		return false
	}
	return ed25519.Verify(pubKey, challenge, signature)
}
