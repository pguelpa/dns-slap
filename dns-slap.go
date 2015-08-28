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
	errorsFound := report(results)

	if errorsFound {
		os.Exit(1)
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

// report pretty-prints the collected results and returns a boolean that
// indicates whether or not any errors were found.
func report(results chan *result) bool {
	fmt.Printf("Workers finished, calculating results...\n\n")

	var totalDuration float64
	var totalCount int
	var errorCount int
	var errorMap = make(map[string]int)
	var minDuration time.Duration
	var maxDuration time.Duration
	var zeroDuration time.Duration

	for result := range results {
		totalDuration += result.duration.Seconds()
		totalCount++
		if result.err != nil {
			errorMap[result.err.Error()]++
			errorCount++
		}
		// if min/maxDuration are set to the zero value then we unconditionally set them to the
		// time of the current result in the loop. This should only run the first time through.
		if minDuration == zeroDuration {
			minDuration = result.duration
		}
		if maxDuration == zeroDuration {
			maxDuration = result.duration
		}
		if result.duration < minDuration {
			minDuration = result.duration
		}
		if result.duration > maxDuration {
			maxDuration = result.duration
		}
	}

	fmt.Println("Results")
	fmt.Println("=======")
	fmt.Println("")

	fmt.Printf("Total lookups: %d\n", totalCount)
	fmt.Printf("Total errors:  %d\n", errorCount)
	fmt.Println("")

	averageDuration := totalDuration / float64(totalCount)
	fmt.Printf("Mean latency:  %fs\n", averageDuration)
	fmt.Printf("Min latency:   %fs\n", minDuration.Seconds())
	fmt.Printf("Max latency:   %fs\n", maxDuration.Seconds())

	if errorCount > 0 {
		fmt.Println("")
		fmt.Println("Error details")
		fmt.Printf("==============\n\n")
		for err, count := range errorMap {
			fmt.Printf("- %s (returned %d times)\n", err, count)
		}
		return true
	}
	return false
}
