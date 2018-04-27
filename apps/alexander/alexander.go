package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tddhit/alexander"
)

var (
	url            string
	postFile       string
	requestNum     int
	concurrency    int
	connectTimeout int
	readTimeout    int
)

func init() {
	flag.StringVar(&url, "url", "", "url")
	flag.StringVar(&postFile, "p", "", "post file")
	flag.IntVar(&requestNum, "n", 1, "the number of requests")
	flag.IntVar(&concurrency, "c", 1, "concurrency")
	flag.IntVar(&connectTimeout, "ct", 1000, "connect timeout")
	flag.IntVar(&readTimeout, "rt", 3000, "read timeout")
	flag.Parse()
}

func main() {
	data, err := alexander.Load(postFile)
	if err != nil {
		log.Fatal(err)
	}

	opt := alexander.Option{
		HTTPVersion:     "1.1",
		ConnectTimeout:  time.Duration(connectTimeout) * time.Millisecond,
		ReadTimeout:     time.Duration(readTimeout) * time.Millisecond,
		IdleConnTimeout: 1000 * time.Millisecond,
		KeepAlive:       30 * time.Second,
		MaxIdleConns:    10,
	}

	header := make(http.Header)
	header.Add("Content-Type", "application/json")

	var wg sync.WaitGroup
	wg.Add(concurrency)
	go func() {
		alexander.Request(url, data, header, requestNum, opt)
		wg.Done()
	}()
	wg.Wait()
}
