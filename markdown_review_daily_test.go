package main

import (
	"testing"
	"time"
)

func TestMatchesCronField(t *testing.T) {
	tests := []struct {
		field    string
		current  int
		max      int
		expected bool
		desc     string
	}{
		{"*", 5, 59, true, "wildcard matches any value"},
		{"5", 5, 59, true, "exact match"},
		{"5", 6, 59, false, "no match"},
		{"1,5,10", 5, 59, true, "list contains value"},
		{"1,5,10", 3, 59, false, "list does not contain value"},
		{"1-5", 3, 59, true, "value in range"},
		{"1-5", 6, 59, false, "value not in range"},
		{"*/15", 0, 59, true, "step matches at 0"},
		{"*/15", 15, 59, true, "step matches at 15"},
		{"*/15", 30, 59, true, "step matches at 30"},
		{"*/15", 7, 59, false, "step does not match at 7"},
		{"10-20/2", 10, 59, true, "range step matches start"},
		{"10-20/2", 12, 59, true, "range step matches middle"},
		{"10-20/2", 11, 59, false, "range step does not match odd"},
	}

	for _, tt := range tests {
		result := matchesCronField(tt.field, tt.current, tt.max)
		if result != tt.expected {
			t.Errorf("%s: matchesCronField(%q, %d, %d) = %v, want %v",
				tt.desc, tt.field, tt.current, tt.max, result, tt.expected)
		}
	}
}

func TestParseCronSchedule(t *testing.T) {
	tests := []struct {
		input       string
		shouldError bool
		desc        string
	}{
		{"0 9 * * *", false, "valid daily schedule"},
		{"*/15 * * * *", false, "valid 15-minute schedule"},
		{"0 9 1 * *", false, "valid monthly schedule"},
		{"0 9 * *", true, "invalid: too few fields"},
		{"0 9 * * * *", true, "invalid: too many fields"},
        {"25 12 * * *", true, "valid daily schedule"},
	}

	for _, tt := range tests {
		_, err := ParseCronSchedule(tt.input)
		hasError := err != nil
		if hasError != tt.shouldError {
			t.Errorf("%s: ParseCronSchedule(%q) error = %v, want error = %v",
				tt.desc, tt.input, err, tt.shouldError)
		}
	}
}

func TestMatchesToday(t *testing.T) {
	// Test with a specific date: January 15, 2026 (Thursday) at various times
	tests := []struct {
		cronExpr string
		testTime time.Time
		expected bool
		desc     string
	}{
		{"0 9 * * *", time.Date(2026, 1, 15, 14, 30, 0, 0, time.UTC), true, "daily at 9 AM - should match any time on correct day"},
		{"0 9 * * *", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), true, "daily at 9 AM - matches at exact time"},
		{"0 10 * * *", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), true, "daily at 10 AM - should match any time on correct day"},
		{"0 9 15 * *", time.Date(2026, 1, 15, 14, 0, 0, 0, time.UTC), true, "15th of month at 9 AM"},
		{"0 9 16 * *", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), false, "16th of month (wrong day)"},
		{"0 9 * 1 *", time.Date(2026, 1, 15, 14, 0, 0, 0, time.UTC), true, "January at 9 AM"},
		{"0 9 * 2 *", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), false, "February (wrong month)"},
		{"0 9 * * 4", time.Date(2026, 1, 15, 16, 0, 0, 0, time.UTC), true, "Thursday at 9 AM (Jan 15 is Thursday)"},
		{"0 9 * * 5", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), false, "Friday (wrong day of week)"},
		{"0 9 * * 1-5", time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC), true, "weekdays at 9 AM"},
		{"0 9 * * 6-7", time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC), false, "weekends at 9 AM"},
	}

	for _, tt := range tests {
		schedule, err := ParseCronSchedule(tt.cronExpr)
		if err != nil {
			t.Errorf("%s: ParseCronSchedule(%q) error: %v", tt.desc, tt.cronExpr, err)
			continue
		}

		result := schedule.MatchesToday(tt.testTime)
		if result != tt.expected {
			t.Errorf("%s: MatchesToday(%q, %v) = %v, want %v",
				tt.desc, tt.cronExpr, tt.testTime.Format("2006-01-02 15:04 Mon"), result, tt.expected)
		}
	}
}

func TestMatchesTime(t *testing.T) {
	// Test with a specific time: January 15, 2026 at 9:00 AM (Thursday)
	testTime := time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		cronExpr string
		expected bool
		desc     string
	}{
		{"0 9 * * *", true, "daily at 9 AM - exact match"},
		{"0 10 * * *", false, "daily at 10 AM (wrong hour)"},
		{"30 9 * * *", false, "daily at 9:30 AM (wrong minute)"},
		{"0 9 15 * *", true, "15th of month at 9 AM"},
		{"0 9 16 * *", false, "16th of month (wrong day)"},
		{"0 9 * 1 *", true, "January at 9 AM"},
		{"0 9 * 2 *", false, "February (wrong month)"},
		{"0 9 * * 4", true, "Thursday at 9 AM"},
		{"0 9 * * 5", false, "Friday (wrong day of week)"},
		{"*/15 * * * *", true, "every 15 minutes (0 matches)"},
		{"*/15 9 * * *", true, "every 15 minutes during 9 AM hour"},
		{"0 9 * * 1-5", true, "weekdays at 9 AM"},
		{"0 9 * * 6-7", false, "weekends at 9 AM"},
	}

	for _, tt := range tests {
		schedule, err := ParseCronSchedule(tt.cronExpr)
		if err != nil {
			t.Errorf("%s: ParseCronSchedule(%q) error: %v", tt.desc, tt.cronExpr, err)
			continue
		}

		result := schedule.MatchesTime(testTime)
		if result != tt.expected {
			t.Errorf("%s: MatchesTime(%q, %v) = %v, want %v",
				tt.desc, tt.cronExpr, testTime.Format("2006-01-02 15:04 Mon"), result, tt.expected)
		}
	}
}

func TestSundayMatching(t *testing.T) {
	// Test Sunday handling (both 0 and 7 should work)
	sundayTime := time.Date(2026, 1, 18, 9, 0, 0, 0, time.UTC) // Sunday

	tests := []struct {
		cronExpr string
		expected bool
		desc     string
	}{
		{"0 9 * * 0", true, "Sunday as 0"},
		{"0 9 * * 7", true, "Sunday as 7"},
		{"0 9 * * 1", false, "Monday"},
	}

	for _, tt := range tests {
		schedule, err := ParseCronSchedule(tt.cronExpr)
		if err != nil {
			t.Errorf("%s: ParseCronSchedule(%q) error: %v", tt.desc, tt.cronExpr, err)
			continue
		}

		result := schedule.MatchesTime(sundayTime)
		if result != tt.expected {
			t.Errorf("%s: MatchesTime(%q, Sunday) = %v, want %v",
				tt.desc, tt.cronExpr, result, tt.expected)
		}
	}
}
