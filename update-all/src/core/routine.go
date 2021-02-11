package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Routine : task to run over a period
type Routine struct {
	Args     []string `yaml:"Args,flow"`
	Interval Interval `yaml:"Interval"`
}

// Interval : minimum time interval to rerun a routine
type Interval struct {
	Hour   int `yaml:"Hour,omitempty"`
	Minute int `yaml:"Minute,omitempty"`
	Second int `yaml:"Second,omitempty"`
}

// ToDuration Convert Interval to `time.Duration`
func (i *Interval) ToDuration() time.Duration {
	return time.Duration(i.Hour)*time.Hour + time.Duration(i.Minute)*time.Minute + time.Duration(i.Second)*time.Second
}

func (r Routine) String() string {
	argStr := strings.Join(r.Args, " ")
	return fmt.Sprintf("\"%v\"", argStr)
}

// Execute : Execute routine
func (r *Routine) Execute() {
	Run(r.Args...)
	log.Debug("Routine executed: ", r)
}

func (r *Routine) hash() string {
	hash := sha256.New()
	for _, arg := range r.Args {
		hash.Write([]byte(arg))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func createRoutine(interval Interval, args ...string) *Routine {
	return &Routine{Interval: interval, Args: args}
}

// DefaultRoutines get exmaple routines to use
func DefaultRoutines() []Routine {
	ret := []Routine{
		*createRoutine(Interval{Minute: 60}, "pyenv", "update"),
		*createRoutine(Interval{}, "pyenv", "rehash"),
		*createRoutine(Interval{Hour: 24}, "poetry", "self", "update"),
		*createRoutine(Interval{Hour: 24}, "tldr", "--update"),
		*createRoutine(Interval{Hour: 24}, "brew", "update"),
		*createRoutine(Interval{Hour: 24}, "rustup", "update"),
	}
	return ret
}

// FlushRoutines write routines to config file
func FlushRoutines(routines []Routine) {
	rawStr, err := yaml.Marshal(&routines)
	check(err)
	fpath := GetRoutineFile()
	ensureDirExists(fpath)
	err = flushToFile(fpath, rawStr)
	check(err)
}

// GetRoutineFile Get the location of the routine file
func GetRoutineFile() string {
	return getActualFileLoc(routineFileName)
}

// IfRoutineFileExists check if config file already exists
func IfRoutineFileExists() bool {
	return ifFileExists(GetRoutineFile())
}

func parseRoutines(data []byte) ([]Routine, error) {
	var routines []Routine
	err := yaml.Unmarshal(data, &routines)
	if err != nil {
		return nil, err
	}
	return routines, nil
}

//LoadRoutines load target routines from config file
func LoadRoutines() ([]Routine, error) {
	rawData, err := ioutil.ReadFile(GetRoutineFile())
	if err != nil {
		return nil, err
	}
	routines, err := parseRoutines(rawData)
	return routines, err
}
