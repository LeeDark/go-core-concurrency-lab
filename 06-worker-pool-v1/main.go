package main

import (
	"fmt"
	"strings"

	"github.com/LeeDark/go-core-concurrency-lab/06-worker-pool-v1/workerpool"
)

func main() {
	jobs := make(chan workerpool.Job)

	results := workerpool.Run(3, jobs,
		func(job workerpool.Job) workerpool.Result {
			value := strings.ToUpper(job.Payload)

			return workerpool.Result{
				JobID: job.ID,
				Value: value,
				Err:   nil,
			}
		})

	go func() {
		defer close(jobs)

		for i, payload := range []string{"one", "two", "three", "four", "five"} {
			jobs <- workerpool.Job{
				ID:      i + 1,
				Payload: payload,
			}
		}
	}()

	for result := range results {
		fmt.Printf("job=%d value=%s err=%v\n", result.JobID, result.Value, result.Err)
	}
}
