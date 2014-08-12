package main

import (
	"flag"
	"fmt"
	"github.com/stvp/clock"
	"github.com/stvp/gostatsd"
	"log"
	"math/rand"
	"time"
)

var (
	// Settings
	statsdUrl        = flag.String("statsd", "statsd://127.0.0.1:8125/firehose.", "Statsd URL including a prefix for all metrics")
	statsdPacketSize = flag.Int("packetsize", 512, "UDP packet size for metrics sent to statsd")

	gaugeCount    = flag.Int("gaugecount", 50000, "Number of individual gauges to run")
	gaugeInterval = flag.Int("gaugeinterval", 60, "Gauge update interval, in seconds")

	counterCount    = flag.Int("countcount", 50000, "Number of individual counters to run")
	counterInterval = flag.Int("countinterval", 60, "Gauge update interval, in seconds")

	// Statistics
	gaugesUpdated   = 0
	countersUpdated = 0

	// Globals
	client statsd.Client
)

func setup() {
	flag.Parse()
	statsd.Setup(*statsdUrl, *statsdPacketSize)
}

func runGauges(count int, interval time.Duration) {
	c, err := clock.New(100*time.Millisecond, interval, 100)
	if err != nil {
		panic(err)
	}
	for key := range keys("g", count) {
		c.Add(key)
	}
	c.Start()

	for key := range c.Channel {
		statsd.Gauge(key, rand.NormFloat64())
		gaugesUpdated++
	}
}

func runCounters(count int, interval time.Duration) {
	c, err := clock.New(100*time.Millisecond, interval, 100)
	if err != nil {
		panic(err)
	}
	for key := range keys("c", count) {
		c.Add(key)
	}
	c.Start()

	for key := range c.Channel {
		statsd.Count(key, 1.0, 1.0)
		countersUpdated++
	}
}

func keys(prefix string, count int) chan string {
	c := make(chan string)
	go func() {
		for i := 0; i < count; i++ {
			c <- fmt.Sprintf("%s.%X", i)
		}
		close(c)
	}()
	return c
}

func main() {
	setup()

	runGauges(*gaugeCount, time.Duration(*gaugeInterval)*time.Second)
	runCounters(*counterCount, time.Duration(*counterInterval)*time.Second)

	go func() {
		for _ = range time.Tick(time.Second) {
			log.Printf("gauges updated: %d", gaugesUpdated)
			log.Printf("counters updated: %d", countersUpdated)
		}
	}()

	// Wait for Ctrl-C
	<-make(chan bool)
}
