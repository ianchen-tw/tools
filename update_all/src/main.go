package main

import (
	"update_all/src/cmd"
)

func main() {
	// routines := core.DefaultRoutines()
	// fmt.Printf("tasks: %+v\n", routines)
	cmd.Execute()
	// m := core.CreateRecordMap()
	// err := m.Load()
	// if err != nil {
	// 	for _, routine := range routines {
	// 		m.UpdateRecord(routine)
	// 	}
	// 	m.Flush()
	// }
	// // fmt.Println(m)

	// core.FlushRoutines(routines)
	// core.LoadRoutines()
}
