package osdat

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// this is the json graph container that contains all JSON stryucture
// usually root json object is the GrandParent process
type JsonGraphContainer struct {
	jsonGraphModel

	f *os.File `json:"-"`
}

type JsonGraphChild struct {
	jsonGraphModel
	// Children    []*JsonGraphChild `json:"children,omitempty"`
	// PID         uint32            `json:"pid"`
	// Name        string            `json:"name"`
	// Description string            `json:"description"`
	// TimeStarted time.Time         `json:"timeStarted"`
	// TimeEnds    time.Time         `json:"timeEnds"`
}

// this model is used for both container adn child
type jsonGraphModel struct {
	PID         uint32 `json:"pid"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CommandLine string `json:"commandLine,omitempty"`

	ExecutablePath string  `json:"execPath,omitempty"`
	Priority       *uint32 `json:"priority,omitempty"`

	TimeStarted time.Time `json:"timeStarted,omitempty"`

	TimeEnds time.Time `json:"timeEnds,omitempty"`
	ExitCode uint32    `json:"exitCode,omitempty"`

	Children []*JsonGraphChild `json:"children,omitempty"`
}

// create new instance of JSON Graph
func NewJsonGraphContainer(name string, pid uint32, logdir string) *JsonGraphContainer {

	// open json log file
	f, err := os.OpenFile(fmt.Sprintf("%s\\osdat.json", logdir), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	jgc := &JsonGraphContainer{
		f: f,
	}

	jgc.Name = name
	jgc.PID = pid
	jgc.Children = []*JsonGraphChild{}

	return jgc
}

// this method will truncate JSON file and rewrite the entire
// logs to it once updates are done
func (jgc *JsonGraphContainer) truncateAndUpdateLog() error {
	err := jgc.f.Truncate(0)
	if err != nil {
		log.Println("truncateAndUpdateLog err[1]:", err)
		return err
	}

	_, err = jgc.f.Seek(0, 0)
	if err != nil {
		log.Println("truncateAndUpdateLog err[1_1]:", err)
		return err
	}

	// convert struct to json
	b, err := json.Marshal(&jgc)
	if err != nil {
		log.Println("truncateAndUpdateLog err[2]:", err)
		return nil
	}

	// write to file
	_, err = jgc.f.Write(b)
	if err != nil {
		log.Println("truncateAndUpdateLog err[3]:", err)
		return nil
	}

	return nil
}

// Find PID on children Level recursive
func (jgchild *JsonGraphChild) findAndUpdateExitCodeByPID(pid, exitCode uint32) (bool, int) {
	found := false
	i := -1

	// find in children and children of children, ... level
	for i, chld := range jgchild.Children {
		if pid == chld.PID {
			found = true

			chld.TimeEnds = time.Now()
			chld.ExitCode = exitCode

			return found, i
		} else {
			// recursive call
			f, _i := chld.findAndUpdateExitCodeByPID(pid, exitCode)
			if f {
				return f, _i
			}
		}
	}

	return found, i
}

// Find PID on children Level recursive
func (jgchild *JsonGraphChild) findAndAppendByPID(pid uint32, child JsonGraphChild) (bool, []*JsonGraphChild, int) {
	found := false
	i := -1

	// find in children and children of children, ... level
	for i, chld := range jgchild.Children {
		if pid == chld.PID {
			found = true
			chld.Children = append(chld.Children, &child)
			return found, chld.Children, i
		} else {
			// recursive call
			f, chl, _i := chld.findAndAppendByPID(pid, child)
			if f {
				return f, chl, _i
			}
		}
	}

	return found, nil, i
}

// Find PID on Container Level
func (jgc *JsonGraphContainer) findByPID(pid uint32) (found bool, i int) {
	found = false
	i = -1

	// find in first level
	for i, chld := range jgc.Children {
		if pid == chld.PID {
			found = true
			return found, i
		}
	}

	return found, i

}
