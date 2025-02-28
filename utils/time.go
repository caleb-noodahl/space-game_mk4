package utils

import "time"

const TICKS_PER_SECOND float64 = 60.0

func ConvertTimeToTicks(dur time.Duration) float64 {
	totalSeconds := dur.Seconds()
	totalTicks := totalSeconds * TICKS_PER_SECOND
	return totalTicks
}
