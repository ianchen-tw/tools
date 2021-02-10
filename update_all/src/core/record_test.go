package core

import (
	"testing"
	"time"
)

func TestRunRecordUpdate(t *testing.T) {
	tests := []struct {
		require Interval
		given   Interval
		expect  bool
	}{
		{require: Interval{Minute: 15},
			given: Interval{Second: 1}, expect: false},
		{require: Interval{Minute: 1},
			given: Interval{Minute: 3}, expect: true},
		{require: Interval{Hour: 1},
			given: Interval{Hour: 99}, expect: true},
	}
	// Patch the underlying function to fix result
	GetCurrentTime = func() time.Time {
		return time.Date(1996, 11, 15, 0, 0, 0, 0, time.UTC)
	}
	for _, tt := range tests {
		args := []string{"ls", "-a", "-l"}
		routine := createRoutine(tt.require, args...)
		lastrun := GetCurrentTime().Add(tt.given.ToDuration() * -1)
		record := RunRecord{Routine: *routine, LastRun: lastrun}
		ans, _ := record.shouldUpdate()
		if ans != tt.expect {
			t.Errorf("Record.shouldUpdate gotResult = %v, expect = %v\ntestcase = %+v", ans, tt.expect, tt)
		}
	}
}
