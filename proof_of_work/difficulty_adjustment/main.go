package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Adjust Difficulty
// Compute new difficulty given old + actual time.
//
// Formula: new_target = old_target * actual_time / expected_time. Cap at 4× change.
//
// Input: <old_target>|<actual_time_seconds>|<expected_time_seconds>.
//
// Output: new target (as integer).
//
// Examples:
// 100|600|600
//   -> 100   (no change; on schedule)
//
// 100|300|600
//   -> 50    (blocks fast → halve target → harder)
//
// 100|2400|600
//   -> 400   (blocks slow → quadruple target → easier; capped at 4×)
//
// 100|9999999|600
//   -> 400   (capped at 4×)
//
// 100|1|600
//   -> 25    (capped at 1/4)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")

		old_target, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		actual_time_seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		expected_time_seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		var new_target float64

		if actual_time_seconds >= expected_time_seconds {
			ratio := getRatio(actual_time_seconds, expected_time_seconds)
			new_target = float64(old_target) * ratio
		} else {
			ratio := getRatio(expected_time_seconds, actual_time_seconds)
			new_target = float64(old_target) / ratio
		}

		fmt.Println(int(new_target))
	}
}

func getRatio(dividend, divisor int) float64 {
	ratio := float64(dividend) / float64(divisor)
	if ratio > 4 {
		ratio = 4
	}

	return ratio
}
