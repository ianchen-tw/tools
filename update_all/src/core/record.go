package core

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
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

// UpdateRecord update the timestamp for a single RunRecord
func (m *RecordMap) UpdateRecord(routine Routine) {
	if record, ok := m.Map[routine.hash()]; ok {
		record.update()
	} else {
		m.Map[routine.hash()] = RunRecord{Routine: routine, LastRun: time.Now()}
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

// Load RecordMap from file
func (m *RecordMap) Load() error {
	rawData, err := ioutil.ReadFile(runRecordFilename)
	if err != nil {
		return err
	}
	json.Unmarshal(rawData, m)
	return nil
}
