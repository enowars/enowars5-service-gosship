package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

const randomPoolSize = 512
const sessionTokenSize = 32

var randomNumberPool []byte
var pubKey ed25519.PublicKey

func init() {
	randomNumberPool = make([]byte, randomPoolSize)
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
	*dest = randomNumberPool[*pos%randomPoolSize]
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
	tokenHash := sha256.New()
	tokenHash.Write(token)
	return hex.EncodeToString(tokenHash.Sum(nil))
}

func VerifySignature(challenge, signature []byte) bool {
	return ed25519.Verify(pubKey, challenge, signature)
}

func CreateAuthChallenge() (string, []byte) {
	challenge := make([]byte, 512)
	_, _ = rand.Read(challenge)
	challengeId := sha256.New()
	challengeId.Write(challenge)
	return hex.EncodeToString(challengeId.Sum(nil)), challenge
}
