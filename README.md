# lightning-distance

[![Go Reference](https://pkg.go.dev/badge/github.com/theluckystrike/lightning-distance.svg)](https://pkg.go.dev/github.com/theluckystrike/lightning-distance)

Compute the distance to a lightning strike from the delay between seeing
the flash and hearing the thunder. Pure, dependency-free Go physics, with
no external data source and no network access.

The light from a strike arrives essentially instantly, so the whole
measurable delay is sound travel time. The speed of sound in air is
modelled with the standard temperature-corrected approximation
**v = 331.3 + 0.606 × T** (m/s, T in °C) — about 343 m/s at 20 °C.

## Install

```sh
go get github.com/theluckystrike/lightning-distance
```

## Quick example

```go
package main

import (
	"fmt"

	lightningdistance "github.com/theluckystrike/lightning-distance"
)

func main() {
	// 5-second flash-to-thunder delay at 20 C
	fmt.Println(lightningdistance.DistanceKm(5.0, 20.0))    // ~1.717 km
	fmt.Println(lightningdistance.DistanceMiles(5.0, 20.0)) // ~1.067 mi
	fmt.Println(lightningdistance.DangerVerdict(5.0))        // "very-close"

	// Rule of thumb
	fmt.Println(lightningdistance.MilesFromSeconds(10.0)) // 2.0 mi
}
```

## API

| Function | Description |
| --- | --- |
| `SpeedOfSound(tempC float64) float64` | Speed of sound in m/s (v = 331.3 + 0.606·T). |
| `DistanceMeters(delaySeconds, tempC float64) float64` | Strike distance in meters (negatives clamped to 0). |
| `DistanceKm(delaySeconds, tempC float64) float64` | Strike distance in kilometers. |
| `DistanceMiles(delaySeconds, tempC float64) float64` | Strike distance in miles (international mile = 1609.344 m). |
| `MilesFromSeconds(delaySeconds float64) float64` | The ~5-seconds-per-mile rule of thumb. |
| `KmFromSeconds(delaySeconds float64) float64` | The ~3-seconds-per-kilometre rule of thumb. |
| `MetersToMiles(meters float64) float64` / `MetersToKm` | Unit helpers. |
| `EstimateTemperatureFromSpeed(speedMps float64) float64` | Inverse of `SpeedOfSound`. |
| `DangerVerdict(delaySeconds float64) string` | `imminent` / `very-close` / `seek-shelter` (30/30 rule) / `near` / `far`. |

## License

MIT.

## Links

- Author's site: <https://earlythunder.com/>
- Package docs: <https://pkg.go.dev/github.com/theluckystrike/lightning-distance>
