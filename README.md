# HTTPBenchmarker | CoderSchool Golang Course Prework 
An HTTP server benchmarking tool written in Go as part of pre-work for a CoderSchool course.

1. **Submitted by: Arturo Locsin**
2. **Time spent: 2 hours**

## List of Features/User Stories
### Required
The following *required* functionalitiy is complete:
* [x] Command-line argument parsing
* [x] Input params
   * [x] Requests - Number of requests to perform
   * [x] Concurrency - Number of multiple requests to make at a time
   * [x] URL - The URL for testing
* [x] Prints use information if wrong arguments provided
* [x] Implements  HTTP load and summarize it
* [x] Concurrency must be implemented with goroutine
### Bonus
The following *optional* features are implemented:
* [x] Extend input params with: 
   * [x] Timeout - Seconds to max. wait for each response
   * [x] Timelimit - Maximum number of seconds to spend for benchmarking
* [x] Prints key metrics of summary, such:
   * [x] Server Hostname
   * [x] Server Port
   * [x] Document Path
   * [x] Document Length
   * [x] Concurrency Level
   * [x] Time taken for tests
   * [x] Complete requests
   * [x] Failed requests
   * [x] Total transferred
   * [x] Requests per second
   * [x] Time per request
   * [x] Transfer rate