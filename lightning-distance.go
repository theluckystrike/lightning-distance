// Package lightningdistance computes the distance to a lightning strike from
// the delay between seeing the flash and hearing the thunder. The light from
// the strike arrives essentially instantly (light travels ~300,000 km/s, so
// even a 50 km strike is seen in well under a millisecond), so the entire
// measurable delay is sound travel time. This is the same physics behind the
// EarlyThunder lightning distance calculator at https://earlythunder.com/.
//
// Speed of sound in air depends on temperature. At 20 °C (68 °F) dry air it is
// ~343 m/s. DistanceKm uses a temperature-corrected speed of sound:
//
//	v = 331.3 + 0.606 * T   (m/s, T in Celsius)
//
// which is the standard linear approximation for the speed of sound in air.
//
// Quick example:
//
//	km := lightningdistance.DistanceKm(5.0, 20.0) // 5 s delay, 20 C
//	// km ~= 1.715  (~331.3 + 0.606*20) * 5 / 1000
package lightningdistance

import "math"

// SpeedOfSound returns the speed of sound in air in m/s for a given air
// temperature in degrees Celsius, using the standard linear model
// v = 331.3 + 0.606 * T.
func SpeedOfSound(tempC float64) float64 {
	return 331.3 + 0.606*tempC
}

// DistanceMeters returns the distance to a lightning strike in meters from the
// flash-to-thunder delay (seconds) at a given air temperature (Celsius). It
// returns 0.0 if the delay is negative.
func DistanceMeters(delaySeconds, tempC float64) float64 {
	if delaySeconds < 0.0 {
		return 0.0
	}
	return SpeedOfSound(tempC) * delaySeconds
}

// DistanceKm returns the strike distance in kilometers.
func DistanceKm(delaySeconds, tempC float64) float64 {
	return DistanceMeters(delaySeconds, tempC) / 1000.0
}

// DistanceMiles returns the strike distance in miles (1 mile = 1609.344 m).
func DistanceMiles(delaySeconds, tempC float64) float64 {
	return DistanceMeters(delaySeconds, tempC) / 1609.344
}

// MilesFromSeconds is the popular 5-seconds-per-mile rule (assuming ~20 °C),
// returning distance in miles. Roughly: count seconds, divide by 5.
func MilesFromSeconds(delaySeconds float64) float64 {
	if delaySeconds < 0.0 {
		return 0.0
	}
	return delaySeconds / 5.0
}

// KmFromSeconds is the popular 3-seconds-per-kilometre rule (assuming ~20 °C),
// returning distance in km. Roughly: count seconds, divide by 3.
func KmFromSeconds(delaySeconds float64) float64 {
	if delaySeconds < 0.0 {
		return 0.0
	}
	return delaySeconds / 3.0
}

// MetersToMiles converts meters to miles using the international mile.
func MetersToMiles(meters float64) float64 {
	return meters / 1609.344
}

// MetersToKm converts meters to kilometers.
func MetersToKm(meters float64) float64 {
	return meters / 1000.0
}

// DangerVerdict returns a plain-language safety read on how far away a strike
// is. The "30/30 rule" says seek shelter if thunder is within 30 seconds
// (~10 km) and wait 30 minutes after the last thunder.
func DangerVerdict(delaySeconds float64) string {
	if delaySeconds < 0.0 {
		return "invalid"
	}
	// Use 20 C as the default reference for the verdict.
	km := DistanceKm(delaySeconds, 20.0)
	switch {
	case km < 1.0:
		return "imminent"
	case km < 5.0:
		return "very-close"
	case km < 10.0:
		return "seek-shelter" // within 30/30-rule range
	case km < 16.0:
		return "near"
	default:
		return "far"
	}
}

// EstimateTemperatureFromSpeed returns the air temperature (Celsius) implied by
// an observed speed of sound (m/s), the inverse of SpeedOfSound. Useful when
// calibrating against a measured sound speed.
func EstimateTemperatureFromSpeed(speedMps float64) float64 {
	return (speedMps - 331.3) / 0.606
}

// ClosestApproach returns the minimum distance (meters) from a strike given the
// observed delay, clamped to >= 0 (you cannot be inside the strike point).
func ClosestApproach(delaySeconds, tempC float64) float64 {
	return math.Max(0.0, DistanceMeters(delaySeconds, tempC))
}
