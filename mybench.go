package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

type responseInfo struct {
	number   int64
	status   int
	bytes    int64
	duration time.Duration
}

type summaryInfo struct {
	host             string
	port             string
	documentPath     string
	documentLength   int64
	requested        int64
	responded        int64
	failed           int64
	totalTransferred int64
	totalDuration    time.Duration
}

func main() {
	fmt.Println("")
	fmt.Println("This is HTTPBenchmarker, a simple HTTP server benchmarking tool")
	fmt.Println("Written by Arturo Locsin in Go as part of pre-work for a CoderSchool course in 2019")
	fmt.Println("Repository: https://github.com/arturolocsin/HTTPBenchmarker")
	fmt.Println("")

	requests := flag.Int64("n", 1, "Number of requests to perform")
	concurrency := flag.Int64("c", 1, "Number of multiple requests to make at a time")
	timeout := flag.Float64("s", 30, "Maximum number of seconds to wait before the socket times out. Default is 30 seconds.")
	timelimit := flag.Float64("t", 0, "Maximum number of seconds to spend for benchmarking. This implies -n 50000.")

	// Read the flags and options
	flag.Parse()
	if flag.NArg() == 0 || *requests == 0 || *requests < *concurrency {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	link := flag.Arg(0)

	fmt.Print("Benchmarking ", link)

	if *timelimit > 0 {
		*requests = 50000
		fmt.Print(" for ", *timelimit, " seconds ")
	}

	fmt.Print("(please wait)...")

	// Start benchmarking
	fmt.Println("")
	starttime := time.Now()
	responseChannel := make(chan responseInfo)
	summary := summaryInfo{}

	for i := int64(0); i < *concurrency; i++ {
		if *timelimit > 0 && time.Since(starttime).Seconds() >= *timelimit {
			break
		}
		summary.requested++
		go checkLink(link, responseChannel, summary.requested, *timeout)
	}

	for response := range responseChannel {
		if *timelimit > 0 && time.Since(starttime).Seconds() >= *timelimit {
			break
		}
		if summary.requested < *requests {
			summary.requested++
			go checkLink(link, responseChannel, summary.requested, *timeout)
		}

		if response.status < 200 || response.status >= 300 || response.status == 0 {
			summary.failed++
		}
		summary.responded++
		fmt.Println(" Request#", response.number, " returned HTTP ", response.status, " with ", response.bytes, " bytes in ", response.duration.Seconds(), "secs")
		if summary.documentLength == 0 {
			summary.documentLength = response.bytes
		}
		summary.totalTransferred += response.bytes
		if summary.responded == summary.requested {
			break
		}
	}
	summary.totalDuration = time.Since(starttime)
	fmt.Printf("\nCompleted!\n\n")

	// Print Info
	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}

	host, port, _ := net.SplitHostPort(u.Host)
	if port == "" {
		summary.host = u.Host
		summary.port = "80"
	} else {
		summary.host = host
		summary.port = port
	}
	summary.documentPath = u.Path

	fmt.Println("SUMMARY")
	fmt.Println(" Server Hostname: ", summary.host)
	fmt.Println(" Server Port: ", summary.port)
	fmt.Println("")
	fmt.Println(" Document Path: ", summary.documentPath)
	fmt.Println(" Document Length: ", summary.documentLength, " bytes")
	fmt.Println("")
	fmt.Println(" Concurrency Level: ", *concurrency)
	fmt.Println(" Time taken for tests: ", summary.totalDuration.Seconds(), "[sec]")
	fmt.Println(" Completed requests: ", summary.responded)
	fmt.Println(" Failed requests: ", summary.failed)
	fmt.Println(" Total transferred: ", summary.totalTransferred, " bytes")
	fmt.Println(" Requests per second: ", float64(summary.requested)/summary.totalDuration.Seconds(), "[#/sec] (mean)")
	fmt.Println(" Time per request: ", summary.totalDuration.Seconds()/float64(summary.requested), "[sec] (mean, across all concurrent requests)")
	fmt.Println(" Transfer rate: ", float64(summary.totalTransferred)/summary.totalDuration.Seconds(), "[bytes/sec] received")
	fmt.Print("\n\n")
}

func checkLink(link string, responseChannel chan responseInfo, requestNumber int64, timeout float64) {

	start := time.Now()

	res, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	read, _ := io.Copy(ioutil.Discard, res.Body)
	elapsed := time.Now().Sub(start)
	if elapsed.Seconds() < float64(timeout) {
		responseChannel <- responseInfo{
			number:   requestNumber,
			status:   res.StatusCode,
			bytes:    read,
			duration: elapsed,
		}
	} else {
		responseChannel <- responseInfo{
			number:   requestNumber,
			status:   408,
			bytes:    0,
			duration: elapsed,
		}
	}
}
