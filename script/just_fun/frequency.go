package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	HIT_TIMES       = 60
	CALCULATE_TIMES = 10
	CALCULATE_STEP  = 10
	ONE_MINUTE_NANO = int64(60 * 1000 * 1000 * 1000)
)

func getAvg(numbers []float64) float64 {
	var sum float64 = 0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func main() {
	timestamps := make([]time.Time, 0, HIT_TIMES)
	fmt.Println("I'm ready, let's start ~")
	for i := 0; i < HIT_TIMES; i++ {
		fmt.Scanln()
		fmt.Printf("got it %d ~\n", i)
		timestamps = append(timestamps, time.Now())
	}

	calculate_ans := make([]float64, 0, CALCULATE_TIMES)
	for i := 0; i < CALCULATE_TIMES; i++ {
		start_ts_index := rand.Int() % (HIT_TIMES - CALCULATE_STEP)
		end_ts_index := start_ts_index + CALCULATE_STEP
		start_time := timestamps[start_ts_index]
		end_time := timestamps[end_ts_index]

		gap_time_nano := end_time.UnixNano() - start_time.UnixNano()
		bpm := float64(ONE_MINUTE_NANO) / float64(gap_time_nano) * CALCULATE_STEP

		calculate_ans = append(calculate_ans, bpm)
	}

	frequency := getAvg(calculate_ans)

	fmt.Println("frequency: ", frequency)
}
