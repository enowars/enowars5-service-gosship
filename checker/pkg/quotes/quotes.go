package quotes

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"io"
	"math/rand"
	"time"
)

/*
Data imported from: https://github.com/vinitshahdeo/inspirational-quotes/blob/18962f9fcd5a2cbe291b205727ee316f05611ef2/data/data.json
LICENSE: MIT / Copyright (c) 2018 Vinit Shahdeo
*/
//go:embed data.csv
var dataCSV []byte

type Quote struct {
	Text string
	From string
}

var allQuotes = make([]*Quote, 0)
var pRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {
	r := csv.NewReader(bytes.NewReader(dataCSV))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		allQuotes = append(allQuotes, &Quote{
			Text: record[0],
			From: record[1],
		})
	}
}

func GetRandom() *Quote {
	return allQuotes[pRand.Intn(len(allQuotes))]
}
