// This program is a CPU stress test used to simulate high load for testing
// a cardiac monitor (Pulsar) visualization. It launches a specified number
// of infinite goroutines (default 12) to generate CPU activity, allowing
// the Pulsar to react to varying CPU usage. The "-n" flag allows the user
// to adjust the number of goroutines to control the intensity of the load.

package main

import (
	"flag"
	"time"
)

func main() {
	n := flag.Int("n", 12, "number of gorutines")
	flag.Parse()

	for i := 0; i < *n; i++ {
		go func() {
			for {
				// infinite loop to generate CPU load
			}
		}()
	}

	// run stress test for 1 minute
	time.Sleep(1 * time.Minute)
}
