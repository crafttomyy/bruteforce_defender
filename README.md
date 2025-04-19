[![GitHub go.mod Go version of a Go module](https://img.shields.io/badge/go-1.24.0-blue)](https://go.dev/dl/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/badge/integration_test-passing-green)](https://go.dev/dl/)

# bruteforce_defender

Defender is a low-level package to help prevent brute force attacks, built on top of golang.org/x/time/rate.

## How it works?
this package cache client ip on memory when its expire recieve, then remove it. <br>
it has some feature when create new defender:
- max (int):max connection or request
- duration (time.Duration):duration time to count up
- banDuration (time.Duration):duration time of ban ip

## Use case
- http middelware for /login
- dos attacks

## Implement code example
there is full example in defender_test.go but for quick start use this code:
```go
package main

import (
	"time"
	"github.com/crafttomyy/bruteforce_defender"
)

func main() {
	// Ban client for 5 Minute if more than 3 events per seconds are performed
	defender:=bruteforce_defender.New(3, 1*time.Minute, 5*time.Minute)
	ip := strings.Split(r.RemoteAddr, ":")[0]
	client := bruteForceDefender.Client(ip)
    // count up rate limiter
	if err := bruteForceDefender.Inc(ip); err != nil {
		log.Fatal(err.Error())
	// Check if the client is banned
	} else if client.Banned() {
		println("Too many requests , try again later")
	}
}
```