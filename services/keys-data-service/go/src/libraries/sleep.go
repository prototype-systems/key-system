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
//
// # Algorithm
//
//  1. The current second of the minute is normalized into the range [0, 1).
//  2. An exponential interpolation is applied so that the delay grows slowly
//     near the beginning of the minute and more rapidly toward the end.
//  3. The interpolated value is mapped into the interval [starter, starter × e].
//  4. A uniformly distributed random jitter of approximately ±18.4%
//     (1 / (2e)) is applied to reduce synchronization between callers.
//  5. The final value is clamped back into the allowed range and rounded
//     to the nearest millisecond.
//
// # Properties
//
//   - Never returns a value less than starter.
//   - Never returns a value greater than starter × e.
//   - Returns deterministic growth based on the current wall-clock second,
//     with a small random variation on each invocation.
//   - Ensures starter is at least 1 ms.
//
// # Example
//
//	Delay(300)
//	// Beginning of the minute: ~300 ms
//	// Middle of the minute:    ~435 ms
//	// End of the minute:       ~815 ms
//
// startingMilliseconds is the minimum delay in milliseconds. Values less
// than 1 are treated as 1.
func Delay(startingMilliseconds int) time.Duration {
	// Ensure the starting milliseconds are at least 1 to avoid degenerate or zero delays.
	if startingMilliseconds < 1 {
		startingMilliseconds = 1
	}

	// Define the minimum delay in milliseconds as the provided starting value.
	minimumDelayMilliseconds := float64(startingMilliseconds)

	// Define the maximum delay in milliseconds as one e-fold above the start (start * e).
	maximumDelayMilliseconds := minimumDelayMilliseconds * math.E

	// Capture the current second-of-minute (0..59) and bound it with modulo to remain in cycle.
	currentSecondOfMinute := time.Now().Second() % 60

	// Normalize the second-of-minute into the [0,1) interval to drive the growth curve smoothly.
	normalizedSecondPosition := float64(currentSecondOfMinute) / 60.0

	// Use Euler's number (e) as the base for the growth curve.
	eulerNumber := math.E

	// Prepare the denominator for the normalized exponential expression; guard against zero.
	denominatorForNormalization := eulerNumber - 1.0
	if denominatorForNormalization == 0.0 {
		denominatorForNormalization = 1.0
	}

	// Compute a normalized exponential growth factor in [0,1]:
	// (e^x - 1) / (e - 1), where x is the normalized second-of-minute.
	normalizedExponentialGrowth := (math.Exp(normalizedSecondPosition) - 1.0) / denominatorForNormalization

	// Interpolate the base delay in milliseconds between minimum and maximum using the growth factor.
	baseDelayMilliseconds := minimumDelayMilliseconds +
		normalizedExponentialGrowth*(maximumDelayMilliseconds-minimumDelayMilliseconds)

	// Choose a small jitter fraction relative to e to de-synchronize concurrent restarts.
	jitterFractionRelativeToEuler := 1.0 / (2.0 * eulerNumber)

	// Generate a uniform random value in [-1, +1] to vary the delay up or down within the jitter band.
	uniformRandomBetweenMinusOneAndOne := 2.0*rand.Float64() - 1.0

	// Compute a multiplicative jitter (1 ± fraction) and apply it to the base delay.
	jitterMultiplier := 1.0 + jitterFractionRelativeToEuler*uniformRandomBetweenMinusOneAndOne

	// Apply jitter to the base delay in milliseconds.
	unclampedDelayMilliseconds := baseDelayMilliseconds * jitterMultiplier

	// Clamp the jittered delay to always remain within [start, start*e].
	clampedDelayMilliseconds := math.Min(
		maximumDelayMilliseconds,
		math.Max(minimumDelayMilliseconds, unclampedDelayMilliseconds),
	)

	// Round to the nearest millisecond and convert to time.Duration for sleeping.
	roundedMilliseconds := math.Round(clampedDelayMilliseconds)
	return time.Duration(roundedMilliseconds) * time.Millisecond
}

// Sleep blocks the calling goroutine for Delay(startingMilliseconds).
func Sleep(startingMilliseconds int) {
	time.Sleep(Delay(startingMilliseconds))
}
