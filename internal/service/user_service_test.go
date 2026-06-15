package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		today    time.Time
		expected int
	}{
		{
			name:     "birthday already passed",
			dob:      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			today:    time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 26,
		},
		{
			name:     "birthday today",
			dob:      time.Date(2000, 6, 15, 0, 0, 0, 0, time.UTC),
			today:    time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 26,
		},
		{
			name:     "birthday not yet reached",
			dob:      time.Date(2000, 12, 20, 0, 0, 0, 0, time.UTC),
			today:    time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 25,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateAge(tc.dob, tc.today)

			if result != tc.expected {
				t.Errorf(
					"expected %d, got %d",
					tc.expected,
					result,
				)
			}
		})
	}
}
