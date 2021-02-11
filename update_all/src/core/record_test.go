package core

import (
	"reflect"
	"testing"
	"time"
)

func TestRecordMapFetch(t *testing.T) {
	m := CreateRecordMap()
	routine := *createRoutine(Interval{Minute: 8}, "echo", "good")

	record := m.fetchRecord(routine)
	if reflect.DeepEqual(record.Routine, routine) {
		t.Errorf("Get different routine expect:%+v, get:%+v", routine, record.Routine)
	}
}

func TestRecordMapUpdate(t *testing.T) {
	date := func(year int, mon time.Month, day int) time.Time {
		return time.Date(year, mon, day, 0, 0, 0, 0, time.UTC)
	}

	// Patch the underlying function to fix result
	GetCurrentTime = func() time.Time {
		return date(1996, 11, 15)
	}

	tests := []time.Time{
		date(1987, 2, 2),
		date(1911, 19, 32),
		date(2023, 9, 8),
	}
	routine := *createRoutine(Interval{Minute: 8}, "echo", "good")
	for _, tt := range tests {
		m := CreateRecordMap()
		record := RunRecord{Routine: routine, LastRun: tt}
		m.Map[record.Routine.hash()] = record
		m.update(record)
		if r, ok := m.Map[record.Routine.hash()]; ok {
			if r.LastRun != GetCurrentTime() {
				t.Errorf("RecordMapUpdate does not update lastRun time to current time.\nexpect=%v, got=%v testcase=%+v", GetCurrentTime(), r.LastRun, tt)
			}
		} else {
			t.Errorf("RecordMapUpdate, could not find stored record, testcase=%+v", tt)
		}
	}

}
func TestRunRecordShouldUpdate(t *testing.T) {
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
