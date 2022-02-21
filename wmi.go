//go:build windows
// +build windows

package osdat

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/bi-zone/wmi"
)

// Start monitoring the instalaltion process
func (e *WmiMonitorProcessEvents) Do() (error, []Monitoredchild) {

	// monitor process create event
	go e.createProcessNotify()

	// monitor process stop trace event
	go e.deleteProcessTraceNotify()

	for {

		select {
		case <-e.alldone:
			e.Stop()
			return nil, []Monitoredchild{}

		case err := <-e.processExitErr:
			e.Stop()
			return err, e.monitoredchild

		case <-e.sigs:
			e.Stop()
			return errors.New("signal recieved from OS or user"),
				e.monitoredchild

		}

	}

}

// this function will check if new process
// is created in the OS and will check if needs to be
// monitored by our app, if needed, it will added the
// in to monitoredChild array
func (e *WmiMonitorProcessEvents) createProcessNotify() error {

	// WMI Query to get the creation events
	// every WITHIN seconds
	// this will help our app to get a callback event
	// from WMI to monitor new processes
	query := `
SELECT * FROM  __InstanceCreationEvent
WITHIN ` + e.within + `
WHERE TargetInstance ISA 'Win32_Process'
`

	// create a new chan based on Win32_Process struct which is
	// win32ProcessEvent struct
	c_events := make(chan win32ProcessEvent)

	// create new NotifyQuery
	// to listen to new events
	// based on our provided query
	q, err := wmi.NewNotificationQuery(c_events, query)

	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to create NotificationQuery; %s", err),
		)
	}

	// start listening to events
	go func() {
		e.errCh <- q.StartNotifications()
	}()

	log.Println("Listening for createprocess events")

	// waiting for events
	for {

		select {

		// check if event recieved from WMI
		case ev := <-c_events:
			// append child process to the list if needed
			// and check if needs to monitor the process
			check := e.appendMonitoredChiled(ev)

			if check {
				log.Printf("[%v] Name: %v, Pid: %v, PPid: %v (MON: %v)\n",
					"Process Created",
					ev.Instance.Caption, ev.Instance.ProcessId, ev.Instance.ParentProcessId,
					check,
				)
			}

		// check if signals are comming from
		// os or user input to terminate the process
		case <-e.sigs:
			// log.Printf("Got system signal %s; stopping", sig)
			// e.Stop()
			return nil

		// Query will never stop here w/o error.
		// if IO error happend, just notify users using log
		case err := <-e.errCh:
			log.Printf("[ERR] Got StartNotifications error; %s", err)
			return nil

		// if all child processes are exited with
		// exit code = 0, this will be run
		case <-e.alldone:
			log.Println("Job done!")
			// e.Stop()
			return nil

		// Process exit with status other than 0
		case err := <-e.processExitErr:
			log.Printf("[ERR] %s\n", err)
			// e.Stop()
			return err

		}

	}

}

// this function will check if a process
// is stopped in the OS and will check if needs be removed
// from monitored process by our app, if needed, it will be removed
// from monitoredChild array
// if process exits with code other than 0, an error channel will be
// filled with the error message an all instances of app
// will be stopped
func (e *WmiMonitorProcessEvents) deleteProcessTraceNotify() error {

	// WMI Query to get the creation events
	// this will help our app to get a callback event
	// from WMI to monitor stopped processes
	query := `SELECT * FROM Win32_ProcessStopTrace`

	// create a new chan based on Win32_ProcessStopTrace struct which is
	// win32ProcessStopEvent struct
	s_events := make(chan win32ProcessStopEvent)

	q, err := wmi.NewNotificationQuery(s_events, query)
	if err != nil {
		return errors.New(
			fmt.Sprintf("Failed to create NotificationQuery; %s", err),
		)

	}

	// AllowMissingFields specifies that struct fields not present in the
	// query result should not result in an error.
	q.AllowMissingFields = true

	// start listening to events
	go func() {
		e.errCh <- q.StartNotifications()
	}()

	log.Println("Listening for process stoptrace events")
	// waiting for events

	for {

		select {

		// check if event recieved from WMI
		case ev := <-s_events:
			if ev.ProcessID == e.grandParentID {
				e.processExitErr <- errors.New(
					fmt.Sprintf("The very first parent '%s(%d)' is exited with exit code '%d'",
						ev.ProcessName, ev.ProcessID, ev.ExitStatus,
					),
				)
			}

			// it will remove process from monitored child if needed
			// also, if check errors and also will check if all
			// child processes are done without error and alldone var will be
			// set to true
			check, alldone, chErr := e.removeMonitoredChiled(ev)

			// if alldone is true, send a sigquit(3) to
			// e.alldone channel
			if alldone {
				e.alldone <- syscall.SIGQUIT
			}

			// if it is last child process
			// and exited with error,
			// we will exit the app, otherwise,
			// it will continiue
			if len(e.monitoredchild) == 1 {
				// if a child process stopped with exit code other than
				// zero, update e.processExitErr with the erro message
				if chErr != nil {
					e.processExitErr <- chErr
				}
			}

			if check {
				log.Printf("[%v] Name: %v, Pid: %v, PPid: %v(MON: %v)\n\tExitStatus: %v\n",
					"Process Stopped",
					ev.ProcessName, ev.ProcessID, ev.ParentProcessID,
					check,
					ev.ExitStatus,
				)
			}

		// check if signals are comming from
		// os or user input to terminate the process
		case <-e.sigs:
			// log.Printf("Got system signal %s; stopping", sig)
			// e.Stop()
			return nil

		// if all child processes are exited with
		// exit code = 0, this will be run
		case err := <-e.errCh:
			log.Printf("[ERR] Got StartNotifications error; %s", err)
			return nil

		// if all child processes are exited with
		// exit code = 0, this will be run
		case <-e.alldone:
			log.Println("Job done!")
			// e.Stop()
			return nil

		// Process exit with status other than 0
		case err := <-e.processExitErr:
			log.Printf("[ERR] %s\n", err)
			// e.Stop()
			os.Exit(1)
			return err
		}

	}

}
