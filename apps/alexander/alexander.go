package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/tddhit/alexander"
	"github.com/tddhit/tools/log"
)

var (
	url            string
	postFile       string
	repeatTimes    int
	concurrency    int
	connectTimeout int
	readTimeout    int
)

func init() {
	flag.StringVar(&url, "url", "", "url")
	flag.StringVar(&postFile, "p", "", "post file")
	flag.IntVar(&repeatTimes, "n", 1, "repeat n times")
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

	var (
		stats       []*alexander.Stats
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
	)

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		s := &alexander.Stats{
			Min: alexander.MAX_TIME_DURATION,
		}
		stats = append(stats, s)
		go func(s *alexander.Stats) {
			alexander.Request(ctx, url, data, header, repeatTimes, opt, s)
			wg.Done()
		}(s)
	}

	go func() {
		log.Error(http.ListenAndServe(":6062", nil))
	}()

	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan)
		<-signalChan
		cancel()
	}()

	wg.Wait()
	alexander.Summarize(stats)
}
