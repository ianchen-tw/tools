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
	Name     string   `yaml:"Name"`
	Args     []string `yaml:"Args"`
	Interval Interval `yaml:"Interval"`
}

// Interval : minimum time interval to rerun a routine
type Interval struct {
	Hour   int `yaml:"Hour"`
	Minute int `yaml:"Minute"`
	Second int `yaml:"Second"`
}

// ToDuration Convert Interval to `time.Duration`
func (i *Interval) ToDuration() time.Duration {
	return time.Duration(i.Hour)*time.Hour + time.Duration(i.Minute)*time.Minute + time.Duration(i.Second)*time.Second
}

func (r Routine) String() string {
	argStr := strings.Join(r.Args, " ")
	return fmt.Sprintf("\"%v %v\"", r.Name, argStr)
}

// Execute : Execute routine
func (r *Routine) Execute() {
	log.Warn("Execute: ", r.String())
}

func (r *Routine) hash() string {
	hash := sha256.New()
	hash.Write([]byte(r.Name))
	for _, arg := range r.Args {
		hash.Write([]byte(arg))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func createRoutine(interval Interval, name string, args []string) *Routine {
	return &Routine{Interval: interval, Name: name, Args: args}
}

// DefaultRoutines get exmaple routines to use
func DefaultRoutines() []Routine {
	var ret []Routine
	ret = append(ret, *createRoutine(Interval{Second: 30}, "echo", []string{"good"}))
	ret = append(ret, *createRoutine(Interval{}, "ls", []string{"-a", "-l"}))
	return ret
}

// FlushRoutines write routines to config file
func FlushRoutines(routines []Routine) {
	rawStr, err := yaml.Marshal(&routines)
	check(err)
	err = flushToFile(routineFileName, rawStr)
	check(err)
}

//LoadRoutines load target routines from config file
func LoadRoutines() ([]Routine, error) {
	rawData, err := ioutil.ReadFile(routineFileName)
	if err != nil {
		return nil, err
	}
	var routines []Routine
	yaml.Unmarshal(rawData, &routines)

	return routines, nil
}
