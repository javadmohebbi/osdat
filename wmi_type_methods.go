//go:build windows
// +build windows

package osdat

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// create new instance of WmiMonitorProcessEvents for monitoring
// installation process and child processes
func NewWMIWmiMonitorProcessEvents(parentId uint32, within string, sig chan os.Signal, graphContainer *JsonGraphContainer) *WmiMonitorProcessEvents {

	ad := make(chan os.Signal, 1)

	signal.Notify(ad, syscall.SIGQUIT)

	// return a new instance of WmiMonitorProcessEvents
	return &WmiMonitorProcessEvents{

		within: within,

		grandParentID: parentId,

		alldone: ad,

		sigs: sig,

		processExitErr: make(chan error, 1),

		errCh: make(chan error, 1),

		monitoredchild: []Monitoredchild{},

		mu: sync.Mutex{},

		Graph: graphContainer,
	}

}

// this is used to provid a pointer array of monitored children
// for the time that parent caller like griAgent needs to stop the child processess
func (e *WmiMonitorProcessEvents) SetMonitoredChildren(mc *[]Monitoredchild) {
	mc = &e.monitoredchild
}

// append process to the list of monitored child
// if does not exists (Actually new items are never exists in the array :-D)
func (e *WmiMonitorProcessEvents) appendMonitoredChiled(ev win32ProcessEvent) (check bool) {

	check = false

	if ev.Instance.ParentProcessId == e.grandParentID {
		e.mu.Lock()
		check = true
		e.monitoredchild = append(e.monitoredchild, Monitoredchild{
			TheFirstParentId: e.grandParentID,
			ProcessID:        ev.Instance.ProcessId,
			Name:             ev.Instance.Name,
			CommandLine:      ev.Instance.CommandLine,
			Caption:          ev.Instance.Caption,
		})

		// // append to graph as it's direct child
		jgchild := JsonGraphChild{}
		jgchild.Name = ev.Instance.Name
		jgchild.Description = ev.Instance.Description
		jgchild.ExecutablePath = ev.Instance.ExecutablePath
		jgchild.CommandLine = ev.Instance.CommandLine
		jgchild.Priority = ev.Instance.Priority
		jgchild.PID = ev.Instance.ProcessId
		jgchild.TimeStarted = time.Now()
		jgchild.Children = []*JsonGraphChild{}
		e.Graph.AppendChild(jgchild)

		e.mu.Unlock()

	} else {
		for _, mc := range e.monitoredchild {
			if ev.Instance.ParentProcessId == mc.ProcessID {
				e.mu.Lock()
				check = true
				e.monitoredchild = append(e.monitoredchild, Monitoredchild{
					TheFirstParentId: e.grandParentID,
					ProcessID:        ev.Instance.ProcessId,
					Name:             ev.Instance.Name,
					CommandLine:      ev.Instance.CommandLine,
					Caption:          ev.Instance.Caption,
				})

				// // append to graph as it's in-direct child
				jgchild := JsonGraphChild{}
				jgchild.Name = ev.Instance.Name
				jgchild.Description = ev.Instance.Description
				jgchild.ExecutablePath = ev.Instance.ExecutablePath
				jgchild.CommandLine = ev.Instance.CommandLine
				jgchild.Priority = ev.Instance.Priority
				jgchild.PID = ev.Instance.ProcessId
				jgchild.TimeStarted = time.Now()
				jgchild.Children = []*JsonGraphChild{}
				e.Graph.AppendNextChild(jgchild,
					ev.Instance.ParentProcessId,
				)

				e.mu.Unlock()
				break

			}

		}

	}

	return check

}

// remove items based on their indexes
func (e *WmiMonitorProcessEvents) _removeMonitoredChildByIndex(s []Monitoredchild, index int) []Monitoredchild {
	return append(s[:index], s[index+1:]...)
}

// remote monitored process from the list
// when a process is terminated, it will be removed from
// monitored child. this will help us to know the steps of instalaltion
func (e *WmiMonitorProcessEvents) removeMonitoredChiled(ev win32ProcessStopEvent) (check, alldone bool, err error) {

	check = false
	alldone = false

	// if no error
	for i, mc := range e.monitoredchild {
		if mc.ProcessID == ev.ProcessID {
			check = true
			e.mu.Lock()
			e.monitoredchild = e._removeMonitoredChildByIndex(e.monitoredchild, i)
			e.mu.Unlock()
			if ev.ExitStatus == 0 {
				if len(e.monitoredchild) == 0 {
					alldone = true
					err = nil
					break
				}
			} else {
				err = errors.New(
					fmt.Sprintf("process '%s (%d) has exited with code '%d'",
						ev.ProcessName, ev.ProcessID, ev.ExitStatus,
					),
				)
			}
		}
	}

	return check, alldone, err

}

// stop all wmi query notification instances
func (e *WmiMonitorProcessEvents) Stop() {

	// stop listening to ProcessCreate events
	if e.qWmiProcessCreate != nil {
		e.qWmiProcessCreate.Stop()
	}

	// stop listening to ProcessStopTrace events
	if e.qWmiProcessStopTrace != nil {
		e.qWmiProcessStopTrace.Stop()
	}

	// write updates before ends
	e.Graph.truncateAndUpdateLog()

	// close json graph
	e.Graph.Close()

}
