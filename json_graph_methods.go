package osdat

import (
	"time"
)

// add child to container
func (jgc *JsonGraphContainer) AppendChild(child JsonGraphChild) {
	// log.Println(">>>>>>", jgc)
	jgc.Children = append(jgc.Children, &child)

	// write to file
	jgc.truncateAndUpdateLog()
	// log.Println(">>>>>>", jgc)
}

// add a child to child
func (jgc *JsonGraphContainer) AppendNextChild(child JsonGraphChild, pid uint32) {

	// find first level
	if found, i := jgc.findByPID(pid); found {
		jgc.Children[i].Children = append(jgc.Children[i].Children, &child)

		// write to file
		jgc.truncateAndUpdateLog()
		return
	}

	// look for other children
	for _, ch := range jgc.Children {
		if found, _, _ := ch.findAndAppendByPID(pid, child); found {
			// log.Println("FOUND")
			// c = append(c, &child)
			// write to file
			jgc.truncateAndUpdateLog()
			return
		}
	}

}

// grand parent | First process exited
func (jgc *JsonGraphContainer) ExitStatusUpdates(pid, exitCode uint32) {
	// find first level
	if found, i := jgc.findByPID(pid); found {
		jgc.Children[i].TimeEnds = time.Now()
		jgc.Children[i].ExitCode = exitCode
		// write to file
		jgc.truncateAndUpdateLog()
		return
	}
}

// child exited
func (jgc *JsonGraphContainer) ChildExitStatusUpdates(pid, exitCode uint32) {

	// find first level
	if found, i := jgc.findByPID(pid); found {
		// log.Println(">>>>>>>>>>>>>>>", jgc.Children[i].Name, pid)
		jgc.Children[i].TimeEnds = time.Now()
		jgc.Children[i].ExitCode = exitCode
		// write to file
		jgc.truncateAndUpdateLog()
		return
	}

	// look for other children
	for _, ch := range jgc.Children {
		if found, _ := ch.findAndUpdateExitCodeByPID(pid, exitCode); found {

			// log.Println("<<<<<", ch.Children[j].Children[i].Name, ch.Children[j].Children[i].PID)

			// write to file
			jgc.truncateAndUpdateLog()
			return
		}
	}

}

// close the graph file handler
func (jgc *JsonGraphContainer) Close() {
	jgc.f.Close()
}
