package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Routine : task to run over a period
type Routine struct {
	Interval int      `yaml:"Interval"`
	Name     string   `yaml:"Name"`
	Args     []string `yaml:"Args"`
}

func (r *Routine) execute() {
	fmt.Printf("Execute routine : %v %v\n", r.Name, r.Args)
}

func (r *Routine) hash() string {
	hash := sha256.New()
	hash.Write([]byte(r.Name))
	for _, arg := range r.Args {
		hash.Write([]byte(arg))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func createRoutine(interval int, name string, args []string) *Routine {
	return &Routine{Interval: interval, Name: name, Args: args}
}

// DefaultRoutines get exmaple routines to use
func DefaultRoutines() []Routine {
	var ret []Routine
	ret = append(ret, *createRoutine(60, "echo", []string{"good"}))
	ret = append(ret, *createRoutine(60, "ls", []string{"-a", "-l"}))
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
func LoadRoutines() []Routine {
	rawData, err := ioutil.ReadFile(routineFileName)
	check(err)
	var routines []Routine
	yaml.Unmarshal(rawData, &routines)

	return routines
}
