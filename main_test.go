package main

import "testing"

func TestAdd(t *testing.T) {
	testCases := []struct {
		a, b     int
		expected int
		name     string
	}{
		{1, 2, 3, "Positive numbers"},
		{-1, 2, 1, "Negative and positive numbers"},
		{-1, -2, -3, "Negative numbers"},
		{0, 0, 0, "Zero values"},
		{100, 200, 300, "Large positive numbers"},
		{-100, -200, -300, "Large negative numbers"},
		{1, -1, 0, "Opposite values"},
		{0, 1, 1, "Zero and positive"},
		{0, -1, -1, "Zero and negative"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Add(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Add(%d, %d) = %d, expected %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}
