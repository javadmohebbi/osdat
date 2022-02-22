//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/javadmohebbi/osdat"
)

// file to analyze
var f2aPath string

// report file path
var rptDirPath string

// array of args
var args []string

// report file
var rpt *os.File

// working dir
var workingDir string

// osdat pid, and ppid
var pid, ppid int

// SIGINT channel
var sigs chan os.Signal

// WMIonitorEvent
var wev *osdat.WmiMonitorProcessEvents

func main() {

	fmt.Println("Start running the analysis...!")

	// monitor the events
	err, children := wev.Do()

	if err != nil {
		log.Println(err)
		// we could decide to kill a process
		// or leave them. But usually we should kill
		// them because their parents are killed
		log.Println(children)

	}

	// cmd := exec.Command(f2aPath, args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stderr

	// log.Fatal(cmd.Run())
}

// initialize app
func init() {

	// get PID and PPID od osdat APP
	pid = os.Getpid()
	ppid = os.Getppid()

	// check signals
	sigs = make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT, // CTRL + C
	)

	// get the working directory of osdat
	workingDir = getCurrentDir()

	// filename that we want to analyze it using this app
	fileToAnalyze := flag.String("f", "", "Executable file path to analysis. currently EXE, BAT, PS1 are supported")

	// command line argument that might require for analysis
	commandlineOptions := flag.String("c", "", "Commandline arguments separeted by comma (,). eg: /F, /D")

	// flag for storing report file
	rptFolder := flag.String("l", "osdat.analysis", "Report directory for the analysis of the provided executable")

	// flag for printing usage
	h := flag.Bool("h", false, "Print help")

	// parse flag arguments
	flag.Parse()

	// print usage if -h provided
	if *h {
		usage(0)
	}

	// check if -f executable is provided
	if *fileToAnalyze == "" {
		fmt.Println("-f must be provided!")
		usage(1)
	}
	f2aPath = *fileToAnalyze

	// check -l report file is provided
	if *rptFolder == "" {
		fmt.Println("-l must be provided!")
		usage(1)
	} else {
		if abs := filepath.IsAbs(*rptFolder); abs {
			rptDirPath = *rptFolder
		} else {
			rptDirPath = fmt.Sprintf("%s\\%s", workingDir, *rptFolder)
		}

		if _, err := os.Stat(rptDirPath); os.IsNotExist(err) {
			// not exist and should be created
			err := os.MkdirAll(rptDirPath, 0777)
			if err != nil {
				log.Fatalln(err)
			}
		}

	}

	// parse the command line if provided
	if *commandlineOptions != "" {
		_args := strings.Split(*commandlineOptions, ",")
		for _, _arg := range _args {
			_a := strings.TrimSpace(_arg)
			if _a != "" {
				args = append(args, _a)
			}
		}
	}

	// check if executable exists
	if e := checkExecExistance(); !e {
		fmt.Printf("It seems that '%s' is not exist!\n", f2aPath)
		os.Exit(4)
	}

	test_pid := 3680

	jgc := osdat.NewJsonGraphContainer("SELF", uint32(test_pid), rptDirPath)

	// create new WMIMonitor
	wev = osdat.NewWMIWmiMonitorProcessEvents(
		// uint32(pid),
		uint32(test_pid),
		"0.01",
		sigs,
		jgc,
	)

}

// check if -f file is exist or not
func checkExecExistance() (exist bool) {

	// default considred as false
	exist = false

	// check releative path
	relPath := fmt.Sprintf("%s\\%s", workingDir, f2aPath)
	if _, err := os.Stat(relPath); err == nil {
		f2aPath = relPath
		exist = true
	}

	// check absPath
	if _, err := os.Stat(f2aPath); err == nil {
		exist = true
	}

	return exist
}

// get osdat executable dir
func getCurrentDir() string {

	d, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	return d

}

// usage will print the app usage and exit with exitCode argument
func usage(exitCode int) {
	usage := `
Usage: osdat.exe [OPTIONS]...

'OSDAT' dynamically analyzes the provided executable and store the report in the 'osdat.json' (by default)

GitHub Repo: https://github.com/javadmohebbi/osdat
Website: https://openintelligence24.com

***NOTE***: If you are analyzing a 'malware' or 'suspicious' execuatable, please run this app in an isolated virtual machine and take a clean snapshot before running the whole process



Examples:
	osdat.exe -f malware1.exe -c /arg1, /arg2, arg2-value -l osdat-malware1
	osdat.exe -f malware2.exe
	osdat.exe -f c:\samples\malbatch.exe



OPTIONS:
	-f {PATH-TO-FILE}	Executable file path to analysis. currently EXE, BAT, PS1 are supported
		{PATH-TO-FILE} could be releative to 'osdat.exe' or an absolute path.


	-c {ARGUMENTS}		Commandline arguments separeted by comma (,). eg: /F, /D
		{ARGUMENTS} must be separeted with comma.


	-l {REPORT-FILE}	Report directory for the analysis of the provided executable


	-h			Print this help
	`

	fmt.Println(usage)

	os.Exit(exitCode)
}
