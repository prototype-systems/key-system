package libraries

import (
	"math"
	"math/rand"
	"time"
)

// Delay computes a restart delay using only the starting milliseconds.
// It shapes growth with Euler's number (e), phases it by the current
// second-of-minute, adds bounded randomization, and constrains the delay
// within one e-fold (start to start*e).
//
// The returned delay is always constrained to the range:
//
//	[starter, starter × e]
//
// where e is Euler's number (math.E ≈ 2.71828).
func Delay(startingMilliseconds int) time.Duration {
	if startingMilliseconds < 1 {
		startingMilliseconds = 1
	}

	minimumDelayMilliseconds := float64(startingMilliseconds)
	maximumDelayMilliseconds := minimumDelayMilliseconds * math.E

	currentSecondOfMinute := time.Now().Second() % 60
	normalizedSecondPosition := float64(currentSecondOfMinute) / 60.0

	eulerNumber := math.E

	denominatorForNormalization := eulerNumber - 1.0
	if denominatorForNormalization == 0.0 {
		denominatorForNormalization = 1.0
	}

	normalizedExponentialGrowth := (math.Exp(normalizedSecondPosition) - 1.0) / denominatorForNormalization

	baseDelayMilliseconds := minimumDelayMilliseconds +
		normalizedExponentialGrowth*(maximumDelayMilliseconds-minimumDelayMilliseconds)

	jitterFractionRelativeToEuler := 1.0 / (2.0 * eulerNumber)
	uniformRandomBetweenMinusOneAndOne := 2.0*rand.Float64() - 1.0
	jitterMultiplier := 1.0 + jitterFractionRelativeToEuler*uniformRandomBetweenMinusOneAndOne

	unclampedDelayMilliseconds := baseDelayMilliseconds * jitterMultiplier

	clampedDelayMilliseconds := math.Min(
		maximumDelayMilliseconds,
		math.Max(minimumDelayMilliseconds, unclampedDelayMilliseconds),
	)

	roundedMilliseconds := math.Round(clampedDelayMilliseconds)
	return time.Duration(roundedMilliseconds) * time.Millisecond
}

// Sleep blocks the calling goroutine for Delay(startingMilliseconds).
func Sleep(startingMilliseconds int) {
	time.Sleep(Delay(startingMilliseconds))
}
