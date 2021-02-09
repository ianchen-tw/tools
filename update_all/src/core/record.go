package core

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// RunRecord Store the time that a routine gets run
type RunRecord struct {
	Routine Routine   `json:"Routine"`
	LastRun time.Time `json:"Time"`
}

func (r *RunRecord) update() {
	r.LastRun = time.Now()
}

// RecordMap a map strcuture to store RunRecord
type RecordMap struct {
	Map map[string]RunRecord `json:"Records"`
}

// CreateRecordMap create a clean RecordMap
func CreateRecordMap() RecordMap {
	return RecordMap{Map: make(map[string]RunRecord)}
}

// RunRoutineIfOutdated run routines that are not runned in the given period
// dry: dry run, do not execute
func (m *RecordMap) RunRoutineIfOutdated(routine Routine, dry bool) {
	var execute bool = false

	if record, ok := m.Map[routine.hash()]; ok {
		timeToLastRun := time.Now().Sub(record.LastRun).Round(100 * time.Millisecond)
		minInterval := record.Routine.Interval.ToDuration()
		log.WithFields(log.Fields{"definedInterval": minInterval, "timeSinceLastRun": timeToLastRun}).Debug("Routine: ", record.Routine.String())
		if timeToLastRun > minInterval {
			execute = true
			record.update()
			m.Map[record.Routine.hash()] = record
		} else {
			log.Info("Skip: ", routine.String(), " excuted ", timeToLastRun, " ago")
		}
	} else {
		execute = true
		m.Map[routine.hash()] = RunRecord{Routine: routine, LastRun: time.Now()}
	}

	if execute {
		routine.Execute()
	}
}

// export RecordMap to byte string
func (m *RecordMap) export() []byte {
	rawStr, _ := json.MarshalIndent(*m, "", strings.Repeat(" ", 4))
	return rawStr
}

// Flush RecordMap to file
func (m *RecordMap) Flush() {
	flushToFile(runRecordFilename, m.export())
}

// TryLoad try to RecordMap from file
func (m *RecordMap) TryLoad() error {
	rawData, err := ioutil.ReadFile(runRecordFilename)
	if err != nil {
		return err
	}
	json.Unmarshal(rawData, m)
	return nil
}
