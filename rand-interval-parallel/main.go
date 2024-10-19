package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Config struct {
	interval          time.Duration
	timeout           time.Duration
	probPassedCeiling float64
	worker            int
}

func getConfig() *Config {
	return &Config{
		interval:          5 * time.Millisecond,
		timeout:           30 * time.Second,
		probPassedCeiling: 0.00005,
		worker:            6,
	}
}

type JailbreakStatus string

const (
	StatusFailed JailbreakStatus = "failed"
	StatusPassed JailbreakStatus = "passed"
)

type JailbreakResult struct {
	WorkerID    int
	Status      JailbreakStatus
	Probability float64
	Attempt     int
}

func (s JailbreakStatus) String() string {
	switch s {
	case StatusPassed:
		return "GGEZ bro!!."
	case StatusFailed:
		return "F**k"
	default:
		return "Ummm"
	}
}

func jailbreak(r *rand.Rand, probability float64) (JailbreakStatus, float64) {
	randNumb := r.Float64()

	var status = StatusFailed

	if randNumb < probability {
		status = StatusPassed
	}

	return status, randNumb
}

func worker(ctx context.Context, workerID int, conf *Config, wg *sync.WaitGroup, results chan<- JailbreakResult) {
	defer wg.Done()

	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

	ticker := time.NewTicker(conf.interval)
	defer ticker.Stop()

	attempt := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			attempt++
			status, prob := jailbreak(r, conf.probPassedCeiling)
			job := JailbreakResult{
				WorkerID:    workerID,
				Status:      status,
				Probability: prob,
				Attempt:     attempt,
			}

			select {
			case results <- job:
			case <-ctx.Done():
				return
			}
		}
	}
}

func jailbreakAttempt(ctx context.Context, conf *Config) error {
	results := make(chan JailbreakResult)
	defer close(results)

	var wg sync.WaitGroup

	for i := 1; i <= conf.worker; i++ {
		wg.Add(1)
		go worker(ctx, i, conf, &wg, results)

	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	attempts := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("pooling timeout!!. after %s with %d attempts", conf.timeout, attempts)
		case job := <-results:
			attempts++
			fmt.Printf("Worker %d | Attempt = %d | Status = %s (prob: %.5f)\n", job.WorkerID, job.Attempt, job.Status, job.Probability)
			if job.Status == StatusPassed {
				fmt.Printf("\nJailbroken. Worker %d win with %d attemp\n", job.WorkerID, job.Attempt)
				return nil
			}
		case <-done:
			return fmt.Errorf("all workers complete but no one can!!")
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	conf := getConfig()

	ctx, cancel := context.WithTimeout(context.Background(), conf.timeout)
	defer cancel()

	fmt.Println("Start pooling..")
	if err := jailbreakAttempt(ctx, conf); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Println("Pooling complete. Yea!")
}
