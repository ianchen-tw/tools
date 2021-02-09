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

func (r *RunRecord) shouldUpdate() (ans bool, timeSinceLastRun time.Duration) {
	timeSinceLastRun = time.Now().Sub(r.LastRun).Round(100 * time.Millisecond)
	minInterval := r.Routine.Interval.ToDuration()
	log.WithFields(log.Fields{"timeSinceLastRun": timeSinceLastRun}).Debug("Routine: ", r.Routine.String())
	return timeSinceLastRun > minInterval, timeSinceLastRun
}

func (m *RecordMap) update(record RunRecord) {
	record.update()
	m.Map[record.Routine.hash()] = record
}

func (m *RecordMap) getRecord(routine Routine) RunRecord {
	if record, ok := m.Map[routine.hash()]; ok {
		return record
	}
	return RunRecord{Routine: routine}
}

// RunRoutineIfOutdated run routines that are not runned in the given period
func (m *RecordMap) RunRoutineIfOutdated(routine Routine, forceUpdate bool, skipExecute bool) {
	record := m.getRecord(routine)
	doUpdate, sinceLastRun := record.shouldUpdate()
	if doUpdate || forceUpdate {
		m.update(record)
	} else {
		log.Info("Skip: ", routine.String(), ", execute ", sinceLastRun, " ago")
	}
	log.Warn("Execute: ", routine.String())
	if !skipExecute {
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
