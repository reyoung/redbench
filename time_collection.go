package main

import (
	"fmt"
	"io"
	"time"
)

type timeCollection struct {
	times []time.Duration
}

func (t *timeCollection) add(d time.Duration) {
	t.times = append(t.times, d)
}

func (t *timeCollection) average() time.Duration {
	var sum time.Duration
	for _, d := range t.times {
		sum += d
	}
	return sum / time.Duration(len(t.times))
}

func (t *timeCollection) max() time.Duration {
	var max time.Duration
	for _, d := range t.times {
		if d > max {
			max = d
		}
	}
	return max
}

func (t *timeCollection) merge(other *timeCollection) {
	t.times = append(t.times, other.times...)
}

func (t *timeCollection) size() int {
	return len(t.times)
}

func (t *timeCollection) prettyPrint(writer io.Writer) {
	fmt.Fprintf(writer, "Average: %f ms, Max: %f ms, N: %d\n",
		t.average().Seconds()*1000, t.max().Seconds()*1000, t.size())
}

func (t *timeCollection) do(callback func()) {
	start := time.Now()
	callback()
	t.add(time.Since(start))
}
