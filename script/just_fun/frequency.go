package main

import (
	"fmt"
	"time"
)

func main() {
	times := 80
	var start_time time.Time
	for i := 0; i < times; i++ {
		fmt.Scanln()
		fmt.Printf("got it %d ~\n", i)
		if start_time.Equal(time.Time{}) {
			start_time = time.Now()
		}
	}
	
	gap_time_nano := time.Now().UnixNano() - start_time.UnixNano()
	one_minute_nano := int64(60 * 1000 * 1000 * 1000)
	times_in_a_minute := float64(one_minute_nano) / float64(gap_time_nano)

	time_of_a_minute := float64(times) * times_in_a_minute

	fmt.Println("frequency: ", time_of_a_minute)
}
