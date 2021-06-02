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

const Timeout = 15000

func createTaskMessagePayload(variant uint64) io.Reader {
	taskMessage := &checker.TaskMessage{
		Method:      checker.TaskMessageMethodPutFlag,
		Address:     "127.0.0.1",
		TeamName:    "team",
		Flag:        "ENO" + randomString(),
		VariantId:   variant,
		Timeout:     Timeout,
		TaskChainId: randomString(),
	}
	rawPayload, err := json.Marshal(taskMessage)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(rawPayload)
}

func sendRequest(group, cnt int, client *http.Client) error {
	variant := uint64(cnt % 2)
	log.Printf("[%02d:%04d:var(%d)]: sending request...", group, cnt, variant)
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
	log.Printf("[%02d:%04d:var( %d)]: %s", group, cnt, variant, string(resRaw))
	return nil
}

func send100Requests(group int, client *http.Client) {
	for i := 0; i < 100; i++ {
		err := sendRequest(group, i, client)
		if err != nil {
			log.Println(err)
		}
	}
}

func run() error {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    Timeout * time.Millisecond,
			DisableCompression: true,
		},
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(group int) {
			send100Requests(group, client)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
