package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type Result struct {
	statusCode int
	duration   time.Duration
}

var (
	url         string
	requests    int
	concurrency int
)

var rootCmd = &cobra.Command{
	Use:   "load-tester",
	Short: "CLI tool for web service load testing",
	Run: func(cmd *cobra.Command, args []string) {
		runLoadTest(url, requests, concurrency)
	},
}

func Execute() {
	rootCmd.Flags().StringVar(&url, "url", "", "Target service URL to test")
	rootCmd.Flags().IntVar(&requests, "requests", 1, "Total number of requests")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Number of concurrent requests")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.MarkFlagRequired("concurrency")

	rootCmd.Execute()
}

func worker(url string, reqs int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < reqs; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{statusCode: 0, duration: duration}
			continue
		}

		results <- Result{statusCode: resp.StatusCode, duration: duration}
		resp.Body.Close()
	}
}

func runLoadTest(url string, totalRequests int, concurrency int) {
	start := time.Now()
	var wg sync.WaitGroup
	results := make(chan Result, totalRequests)

	reqsPerWorker := totalRequests / concurrency
	extra := totalRequests % concurrency

	for i := 0; i < concurrency; i++ {
		count := reqsPerWorker
		if i < extra {
			count++
		}
		wg.Add(1)
		go worker(url, count, results, &wg)
	}

	wg.Wait()
	close(results)

	total := 0
	okCount := 0
	statusMap := make(map[int]int)
	var totalDuration time.Duration

	for r := range results {
		total++
		statusMap[r.statusCode]++
		if r.statusCode == 200 {
			okCount++
		}
		totalDuration += r.duration
	}

	fmt.Println("\nLoad Test Report:")
	fmt.Printf("Total time: %v\n", time.Since(start))
	fmt.Printf("Total requests: %d\n", total)
	fmt.Printf("HTTP 200 OK: %d\n", okCount)
	fmt.Println("HTTP Status Distribution:")
	for code, count := range statusMap {
		fmt.Printf("  %d: %d\n", code, count)
	}
}
