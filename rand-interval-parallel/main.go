package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Config struct {
	interval          time.Duration
	timeout           time.Duration
	probPassedCeiling float64
	worker            int
}

type JailbreakStatus string

const (
	StatusFailed JailbreakStatus = "failed"
	StatusPassed JailbreakStatus = "passed"
)

type JailbreakResult struct {
	Status      JailbreakStatus
	Probability float64
	Attempt     int
}

func getConfig() *Config {
	return &Config{
		interval:          10 * time.Millisecond,
		timeout:           30 * time.Second,
		probPassedCeiling: 0.0005,
		worker:            5,
	}
}

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

func jailbreak(probability float64, attempt int) JailbreakResult {
	randNumb := rand.Float64()

	var status = StatusFailed

	if randNumb < probability {
		status = StatusPassed
	}

	return JailbreakResult{
		Status:      status,
		Probability: randNumb,
		Attempt:     attempt,
	}
}

func worker(ctx context.Context, prob float64, attempt int, chanResult chan<- JailbreakResult) {
	select {
	case <-ctx.Done():
		return
	default:
		jobResult := jailbreak(prob, attempt)
		select {
		case chanResult <- jobResult:
		case <-ctx.Done():
		}
	}
}

func jailbreakAttempt(ctx context.Context, conf *Config) error {
	ticker := time.NewTicker(conf.interval)
	defer ticker.Stop()

	chanResult := make(chan JailbreakResult)
	defer close(chanResult)

	attempt := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("pooling timeout!!. after %s with %d attempts", conf.timeout, attempt)

		case <-ticker.C:
			attempt++
			go worker(ctx, conf.probPassedCeiling, attempt, chanResult)
		case jobResult := <-chanResult:
			fmt.Printf("attempt = %d, status = %s (prob: %.4f)\n", jobResult.Attempt, jobResult.Status, jobResult.Probability)
			if jobResult.Status == StatusPassed {
				fmt.Println("complete!! exiting pooling")
				return nil
			}
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixMicro()))

	conf := getConfig()

	ctx, cancel := context.WithTimeout(context.Background(), conf.timeout)
	defer cancel()

	fmt.Println("start pooling...")
	if err := jailbreakAttempt(ctx, conf); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Println("pooling complete. yea")
}
