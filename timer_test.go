package evil_proxy

import (
	"testing"
	"time"
)

func TestFuzzyEqualsMatchesEqualValues(t *testing.T) {
	if !FuzzyEquals(0, 0, 0) {
		t.Fatalf("Expected match with exact values\n")
	}
}

func TestFuzzyEqualsMatchesCloseValues(t *testing.T) {
	if !FuzzyEquals(0, 1, 1) {
		t.Fatalf("Expected match with close values\n")
	}
}

func TestFuzzyEqualsDoesntMatchFarValues(t *testing.T) {
	if FuzzyEquals(0, 10, 1) {
		t.Fatalf("Expected mismatch with far values\n")
	}
}

func TestElapsedMillisecondsReportsElapsedTime(t *testing.T) {
	const maxdelay = time.Millisecond * 300
	duration := StartTimer().ElapsedMilliseconds()
	if !FuzzyEquals(0, duration, maxdelay) {
		t.Fatalf("Expected a small time had elapsed %v\n", duration)
	}
}
