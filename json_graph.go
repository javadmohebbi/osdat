package osdat

import (
	"log"
	"time"
)

type JsonGraphContainer struct {
	Name        string            `json:"name"`
	PID         uint32            `json:"pid"`
	Children    []*JsonGraphChild `json:"children,omitempty"`
	TimeStarted time.Time         `json:"timeStarted"`
	TimeEnds    time.Time         `json:"timeEnds"`
}

func NewJsonGraphContainer(name string, pid uint32) *JsonGraphContainer {
	return &JsonGraphContainer{
		Name:     name,
		PID:      pid,
		Children: []*JsonGraphChild{},
	}
}

type JsonGraphChild struct {
	Children    []*JsonGraphChild `json:"children,omitempty"`
	Name        string            `json:"name"`
	PID         uint32            `json:"pid"`
	TimeStarted time.Time         `json:"timeStarted"`
	TimeEnds    time.Time         `json:"timeEnds"`
}

// Find PID on children Level recursive
func (jgchild *JsonGraphChild) findByPID(pid uint32, child JsonGraphChild) (bool, []*JsonGraphChild, int) {
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
			f, chl, _i := chld.findByPID(pid, child)
			if f {
				return f, chl, _i
			}
		}
	}

	return found, nil, i
}

// add child to container
func (jgc *JsonGraphContainer) AppendChild(child JsonGraphChild) {
	// log.Println(">>>>>>", jgc)
	jgc.Children = append(jgc.Children, &child)
	// log.Println(">>>>>>", jgc)
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

// add child to child
func (jgc *JsonGraphContainer) AppendNextChild(child JsonGraphChild, pid uint32) {

	// find first level
	if found, i := jgc.findByPID(pid); found {
		jgc.Children[i].Children = append(jgc.Children[i].Children, &child)
		return
	}

	// look for other children
	for _, ch := range jgc.Children {
		if found, _, _ := ch.findByPID(pid, child); found {
			log.Println("FOUND")
			// c = append(c, &child)
			return
		}
	}

}

// return string
// func (jgc *JsonGraphContainer) String() string {
// 	b, err := json.Marshal(jgc)
// 	if err != nil {
// 		return string(b)
// 	}

// 	return "ERR:" + err.Error()
// }
