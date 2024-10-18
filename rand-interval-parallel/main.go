package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Config struct {
	interval     time.Duration
	timeout      time.Duration
	probComplete float64
}

type JailbreakStatus string

const (
	StatusFailed JailbreakStatus = "failed"
	StatusPassed JailbreakStatus = "passed"
)

func (s JailbreakStatus) String() string {
	switch s {
	case StatusPassed:
		return "GGEZ bro!!."
	case StatusFailed:
		return "Oh..s**t, f**k"
	default:
		return "Ummm"
	}
}

func getConfig() *Config {
	return &Config{
		interval:     10 * time.Millisecond,
		timeout:      30 * time.Second,
		probComplete: 0.005,
	}
}

func jailbreak(prob float64) (JailbreakStatus, float64) {
	randNumb := rand.Float64()

	var status = StatusFailed

	if randNumb < prob {
		status = StatusPassed
	}

	return status, randNumb
}

func attempJailbreak(conf *Config) error {
	ticker := time.NewTicker(conf.interval)
	defer ticker.Stop()

	timeoutTimer := time.NewTimer(conf.timeout)
	defer timeoutTimer.Stop()

	attemp := 0

	for {
		select {
		case <-ticker.C:
			attemp++
			status, randVal := jailbreak(conf.probComplete)
			fmt.Printf("attemp = %d, status = %s (prob: %.4f)\n", attemp, status, randVal)
			if status == "passed" {
				fmt.Println("complete!!. exist pooling")
				return nil
			}
		case <-timeoutTimer.C:
			return fmt.Errorf("pooling timeout!!. after %s with %d attemps", conf.timeout, attemp)
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixMicro()))

	conf := getConfig()

	_, cancel := context.WithTimeout(context.Background(), conf.timeout)
	defer cancel()

	fmt.Println("start pooling...")
	if err := attempJailbreak(conf); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Println("pooling complete. yea")
}
