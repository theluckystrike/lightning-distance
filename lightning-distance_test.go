package lightningdistance

import (
	"math"
	"testing"
)

func approx(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func TestSpeedOfSound(t *testing.T) {
	// At 20 C -> 331.3 + 0.606*20 = 343.42 m/s
	v := SpeedOfSound(20.0)
	if !approx(v, 343.42, 1e-6) {
		t.Fatalf("SpeedOfSound(20C) = %v, want 343.42", v)
	}
	// At 0 C -> 331.3
	if !approx(SpeedOfSound(0.0), 331.3, 1e-9) {
		t.Fatalf("SpeedOfSound(0C) = %v, want 331.3", SpeedOfSound(0.0))
	}
}

func TestDistanceMeters(t *testing.T) {
	// 5 s at 20 C -> 343.42 * 5 = 1717.1 m
	m := DistanceMeters(5.0, 20.0)
	if !approx(m, 1717.1, 1e-3) {
		t.Fatalf("DistanceMeters = %v, want 1717.1", m)
	}
}

func TestDistanceKm(t *testing.T) {
	km := DistanceKm(5.0, 20.0)
	if !approx(km, 1.7171, 1e-5) {
		t.Fatalf("DistanceKm = %v, want 1.7171", km)
	}
}

func TestDistanceMiles(t *testing.T) {
	// 1717.1 m -> 1.0669 mi
	mi := DistanceMiles(5.0, 20.0)
	if !approx(mi, 1.06694, 1e-4) {
		t.Fatalf("DistanceMiles = %v, want ~1.0669", mi)
	}
}

func TestNegativeDelayClamps(t *testing.T) {
	if DistanceMeters(-1.0, 20.0) != 0.0 {
		t.Fatal("negative delay should clamp to 0")
	}
	if MilesFromSeconds(-2.0) != 0.0 {
		t.Fatal("negative delay should clamp to 0")
	}
}

func TestMilesFromSecondsRule(t *testing.T) {
	// 5s -> 1 mile (the rule of thumb)
	if !approx(MilesFromSeconds(5.0), 1.0, 1e-9) {
		t.Fatalf("MilesFromSeconds(5) = %v, want 1.0", MilesFromSeconds(5.0))
	}
	if !approx(MilesFromSeconds(10.0), 2.0, 1e-9) {
		t.Fatalf("MilesFromSeconds(10) = %v, want 2.0", MilesFromSeconds(10.0))
	}
}

func TestKmFromSecondsRule(t *testing.T) {
	// 3s -> 1 km
	if !approx(KmFromSeconds(3.0), 1.0, 1e-9) {
		t.Fatalf("KmFromSeconds(3) = %v, want 1.0", KmFromSeconds(3.0))
	}
}

func TestMetersToMiles(t *testing.T) {
	// 1609.344 m == exactly 1 mile
	if !approx(MetersToMiles(1609.344), 1.0, 1e-9) {
		t.Fatalf("MetersToMiles(1mi) = %v, want 1.0", MetersToMiles(1609.344))
	}
}

func TestMetersToKm(t *testing.T) {
	if !approx(MetersToKm(2500.0), 2.5, 1e-9) {
		t.Fatalf("MetersToKm(2500) = %v, want 2.5", MetersToKm(2500.0))
	}
}

// Round-trip: speed -> temperature -> back to speed.
func TestRoundTripTempSpeed(t *testing.T) {
	for _, tc := range []float64{0.0, 15.0, 20.0, 30.0, -10.0} {
		v := SpeedOfSound(tc)
		back := EstimateTemperatureFromSpeed(v)
		if !approx(back, tc, 1e-6) {
			t.Fatalf("round-trip temp = %v, want %v (from speed %v)", back, tc, v)
		}
	}
}

func TestDangerVerdictBuckets(t *testing.T) {
	cases := map[float64]string{
		1.0:  "imminent",
		5.0:  "very-close",
		20.0: "seek-shelter", // ~6.9 km at 20 C
		40.0: "near",
		60.0: "far",
	}
	for delay, want := range cases {
		if got := DangerVerdict(delay); got != want {
			t.Fatalf("DangerVerdict(%vs) = %q, want %q", delay, got, want)
		}
	}
}

func TestClosestApproachNeverNegative(t *testing.T) {
	if ClosestApproach(-5.0, 20.0) != 0.0 {
		t.Fatal("ClosestApproach should clamp negatives to 0")
	}
}
