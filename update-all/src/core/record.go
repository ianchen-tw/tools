package core

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
)

// RunRecord Store the time that a routine gets run
type RunRecord struct {
	Routine Routine   `json:"Routine"`
	LastRun time.Time `json:"Time"`
}

// RecordMap a map strcuture to store RunRecord
type RecordMap struct {
	Map map[string]RunRecord `json:"Records"`
}

// GetCurrentTime : function to get the system time (override in testing environment)
//	https://stackoverflow.com/questions/18970265/is-there-an-easy-way-to-stub-out-time-now-globally-during-test
var GetCurrentTime = time.Now

func (r *RunRecord) shouldUpdate() bool {
	timeSinceLastRun := GetCurrentTime().Sub(r.LastRun).Round(100 * time.Millisecond)
	minInterval := r.Routine.Interval.ToDuration()
	log.WithFields(log.Fields{"lastRun": r.LastRun, "timeSinceLastRun": timeSinceLastRun}).Debug("Get Routine info: ", r.Routine.String())
	return timeSinceLastRun > minInterval
}

// CreateRecordMap create a clean RecordMap
func CreateRecordMap() RecordMap {
	return RecordMap{Map: make(map[string]RunRecord)}
}

func (m *RecordMap) update(record RunRecord) {
	record.LastRun = GetCurrentTime()
	m.Map[record.Routine.hash()] = record
}

// fetchRecord get a record from routine to test
func (m *RecordMap) fetchRecord(routine Routine) RunRecord {
	if record, ok := m.Map[routine.hash()]; ok {
		return record
	}
	return RunRecord{Routine: routine}
}

// RunRoutineIfOutdated run routines that are not runned in the given period
func (m *RecordMap) RunRoutineIfOutdated(routine Routine, forceUpdate bool, skipExecute bool) {
	record := m.fetchRecord(routine)
	log.Debug("Get record from map: ", record)
	record.Routine = routine
	doUpdate := record.shouldUpdate()
	if doUpdate || forceUpdate {
		m.update(record)
		log.Info(emoji_execute, "Execute: ", routine.String())
		if !skipExecute {
			routine.Execute()
		}
		return
	}
	log.Info(emoji_skip, "Skip: ", routine.String(), ", execute ", humanize.Time(record.LastRun))
}

// export RecordMap to byte string
func (m *RecordMap) export() []byte {
	rawStr, _ := json.MarshalIndent(*m, "", strings.Repeat(" ", 2))
	return rawStr
}

// GetRecordMapFile get the location of recordMap file
func GetRecordMapFile() string {
	return getActualFileLoc(runRecordFilename)
}

// Flush RecordMap to file
func (m *RecordMap) Flush() {
	fpath := GetRecordMapFile()
	ensureDirExists(fpath)
	flushToFile(fpath, m.export())
}

func (m *RecordMap) parseRecordMap(data []byte) error {
	err := json.Unmarshal(data, m)
	return err
}

// TryLoad try to RecordMap from file
func (m *RecordMap) TryLoad() error {
	rawData, err := ioutil.ReadFile(GetRecordMapFile())
	if err != nil {
		return err
	}
	return m.parseRecordMap(rawData)
}
