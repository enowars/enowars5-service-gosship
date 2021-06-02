package main

import (
	"bytes"
	"checker/pkg/checker"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{})
}

func randomString() string {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

const Timeout = 15000

func createPutFlagPayload(variant uint64) io.Reader {
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

func createPutNoisePayload() io.Reader {
	taskMessage := &checker.TaskMessage{
		Method:      checker.TaskMessageMethodPutNoise,
		Address:     "127.0.0.1",
		TeamName:    "team",
		VariantId:   0,
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
	var payload io.Reader
	var pType string
	if group%2 == 0 {
		pType = "F"
		payload = createPutFlagPayload(variant)
	} else {
		pType = "N"
		payload = createPutNoisePayload()
	}
	logPrefix := fmt.Sprintf("[%02d:%03d:%s:var(%d)]:", group, cnt, pType, variant)

	start := time.Now()
	log.Printf("%s sending request...", logPrefix)

	request, err := http.NewRequest("POST", "http://localhost:2002/", payload)
	if err != nil {
		return err
	}
	res, err := client.Do(request)
	if err != nil {
		return err
	}

	var checkerRes checker.ResultMessage
	if err := json.NewDecoder(res.Body).Decode(&checkerRes); err != nil {
		return err
	}

	log.Printf("%s done(%dms): %s", logPrefix, time.Since(start).Milliseconds(), checkerRes.Result)
	if checkerRes.Result != checker.ResultOk {
		log.Errorf("%s %s", logPrefix, *checkerRes.Message)
	}

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
			IdleConnTimeout:    Timeout * time.Millisecond,
			DisableCompression: true,
		},
		Timeout: Timeout * time.Millisecond,
	}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
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
