package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/javadmohebbi/osdat"
)

var jsonstr = `{"pid":2568,"name":"SELF","timeStarted":"0001-01-01T00:00:00Z","timeEnds":"0001-01-01T00:00:00Z","children":[{"pid":1676,"name":"powershell.exe","description":"powershell.exe","commandLine":"powershell","priority":8,"timeStarted":"2022-02-22T14:58:19.8688858+03:30","timeEnds":"0001-01-01T00:00:00Z","children":[{"pid":956,"name":"powershell.exe","description":"powershell.exe","commandLine":"\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\"","execPath":"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe","priority":8,"timeStarted":"2022-02-22T14:58:34.6846497+03:30","timeEnds":"2022-02-22T15:00:10.6211627+03:30","exitCode":1,"children":[{"pid":3364,"name":"cmd.exe","description":"cmd.exe","commandLine":"\"C:\\Windows\\system32\\cmd.exe\"  \\admin","priority":8,"timeStarted":"2022-02-22T14:58:56.6642375+03:30","timeEnds":"2022-02-22T14:59:21.1810729+03:30","children":[{"pid":2932,"name":"calc.exe","description":"calc.exe","commandLine":"calc","priority":8,"timeStarted":"2022-02-22T14:59:07.5491311+03:30","timeEnds":"2022-02-22T14:59:11.0266881+03:30"}]},{"pid":2816,"name":"calc.exe","description":"calc.exe","commandLine":"\"C:\\Windows\\system32\\calc.exe\"","execPath":"C:\\Windows\\system32\\calc.exe","priority":8,"timeStarted":"2022-02-22T14:59:22.2817595+03:30","timeEnds":"2022-02-22T14:59:24.1165666+03:30"},{"pid":3596,"name":"taskmgr.exe","description":"taskmgr.exe","commandLine":"\"C:\\Windows\\system32\\taskmgr.exe\"","priority":8,"timeStarted":"2022-02-22T14:59:27.4138009+03:30","timeEnds":"2022-02-22T14:59:29.0581395+03:30","exitCode":1}]}]}]}`

func main() {
	var jgc *osdat.JsonGraphContainer
	err := json.Unmarshal([]byte(jsonstr), &jgc)
	if err != nil {
		log.Fatalln(err)
	}

	p := flag.Uint("pid", 0, "Find PID")
	flag.Parse()

	findByPID(uint32(*p), jgc)
}

func findByPID(pid uint32, jgc *osdat.JsonGraphContainer) {
	if jgc.PID == pid {
		log.Println("GRAND PARENT")
		return
	}

	for _, ch := range jgc.Children {
		_findByPID(pid, ch)
	}

}

func _findByPID(pid uint32, jgc *osdat.JsonGraphChild) {
	if jgc.PID == pid {
		fmt.Println(jgc.Name, jgc.PID)
		return
	}
	if len(jgc.Children) == 0 {
		return
	}
	for _, ch := range jgc.Children {
		_findByPID(pid, ch)
	}
}

func printRecursive(level int, sublevel int, jgc osdat.JsonGraphChild) {
	for t := 0; t < sublevel; t++ {
		fmt.Printf("\t")
	}
	if len(jgc.Children) == 0 {
		return
	}
	fmt.Printf("\t[%d] %s (%d)\n", level, jgc.Name, jgc.PID)
	for i, ch := range jgc.Children {
		fmt.Printf("\t\t %s (%d):\n", ch.Name, ch.PID)
		printRecursive(level+1, i, *ch)
	}
	fmt.Println()
}
