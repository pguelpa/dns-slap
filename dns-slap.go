package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var concurrency = flag.Int("concurrency", 10, "How many concurrent lookups to try")
var iterations = flag.Int("iterations", 100, "How many times to lookup in each concurrent process")
var thresholdMilli = flag.Int("threshold", 500, "How long to wait (in milliseconds) on a single lookup before considering it a failure")

type result struct {
	duration time.Duration
	err      error
}

func main() {
	flag.Parse()

	var results = make(chan *result, *concurrency**iterations)
	var wg sync.WaitGroup

	if len(flag.Args()) != 1 {
		fmt.Println("Usage of", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	host := flag.Args()[0]

	fmt.Printf("Starting %d workers with %d lookups each ...\n", *concurrency, *iterations)
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go work(i, host, *iterations, &wg, results)
	}

	wg.Wait()
	close(results)

	fmt.Printf("Workers finished, calculating results\n\n")

	var totalDuration float64
	var totalCount int
	var errorCount int
	var errorMap = make(map[string]int)

	for result := range results {
		totalDuration += result.duration.Seconds()
		totalCount++
		if result.err != nil {
			errorMap[result.err.Error()]++
			errorCount++
		}
	}

	averageDuration := totalDuration / float64(totalCount)
	fmt.Printf("Ran %d lookups in an average time of %f seconds\n", totalCount, averageDuration)
	fmt.Printf("Found %d errors\n", errorCount)
	for err, count := range errorMap {
		fmt.Printf("\t%s returned %d times\n", err, count)
	}
}

func work(index int, host string, count int, wg *sync.WaitGroup, results chan<- *result) {
	var start time.Time
	var duration time.Duration
	var err error

	threshold := time.Millisecond * time.Duration(*thresholdMilli)

	for i := 0; i < count; i++ {
		start = time.Now()
		_, err = net.LookupHost(host)
		duration = time.Since(start)
		if err != nil {
			results <- &result{duration: duration, err: err}
		} else {
			if duration > threshold {
				err := fmt.Errorf("Lookup succeeded but took longer than the allowed threshold of %dms", *thresholdMilli)
				results <- &result{duration: duration, err: err}
			} else {
				results <- &result{duration: duration}
			}
		}
	}

	wg.Done()
}
