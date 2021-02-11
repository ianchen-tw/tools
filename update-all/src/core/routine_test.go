package core

import (
	"testing"
)

func TestParseRoutines(t *testing.T) {
	data := `
- Args: [pyenv, update]
  Interval:
    Minute: 60
- Args: [pyenv, rehash]
  Interval: {}
- Args: [tldr, --update]
  Interval:
    Hour: 24
`
	_, err := parseRoutines([]byte(data))
	if err != nil {
		t.Error("Cannot parse data", data, err)
	}
}
func TestRoutineHash(t *testing.T) {
	r := []Routine{
		*createRoutine(Interval{Minute: 2}, "ls", "-a"),
		*createRoutine(Interval{Minute: 3}, "ls", "-a"),
	}
	if r[0].hash() != r[1].hash() {
		t.Error("Routine's hash value should not involve Interval information")
	}
}

func TestIntervalToDuration(t *testing.T) {

	tests := []struct {
		a Interval
		b Interval
	}{
		{
			a: Interval{Minute: 2},
			b: Interval{Second: 120},
		},
		{
			a: Interval{Minute: 2},
			b: Interval{Minute: 1, Second: 60},
		},
		{
			a: Interval{Hour: 2, Minute: 3, Second: 9},
			b: Interval{Minute: 120, Second: 189},
		},
	}
	for _, tt := range tests {
		if tt.a.ToDuration() != tt.b.ToDuration() {
			t.Error("Inteval convert error", tt.a, "!=", tt.b)
		}
	}

}
