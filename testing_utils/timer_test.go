package testing_utils

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
	const maxdelay = time.Millisecond * 1000
	duration := StartTimer().ElapsedMilliseconds()
	if duration < 0 {
		t.Fatalf("Expected a possitive time change. Got %v\n", duration)
	}
	if !FuzzyEquals(0, duration, maxdelay) {
		t.Fatalf("Expected a small time had elapsed %v\n", duration)
	}
}
