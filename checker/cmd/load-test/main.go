package main

import (
	"bytes"
	"checker/pkg/checker"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func randomString() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func createTaskMessagePayload(variant uint64) io.Reader {
	taskMessage := &checker.TaskMessage{
		Method:      checker.TaskMessageMethodPutFlag,
		Address:     "127.0.0.1",
		TeamName:    "team",
		Flag:        "ENO" + randomString(),
		VariantId:   variant,
		Timeout:     15000,
		TaskChainId: randomString(),
	}
	rawPayload, err := json.Marshal(taskMessage)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(rawPayload)
}

func sendRequest(variant uint64, client *http.Client) error {
	log.Printf("sending request (%d)...", variant)
	request, err := http.NewRequest("POST", "http://localhost:2002/", createTaskMessagePayload(variant))
	if err != nil {
		return err
	}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	resRaw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Printf("%s", string(resRaw))
	return nil
}

func send100Requests(client *http.Client) {
	for i := 0; i < 100; i++ {
		err := sendRequest(uint64(i%2), client)
		if err != nil {
			log.Println(err)
		}
	}
}

func run() error {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    15 * time.Second,
			DisableCompression: true,
		},
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			send100Requests(client)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
