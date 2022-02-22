package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/javadmohebbi/osdat"
)

var jsonstr = `{"pid":3680,"name":"SELF","timeStarted":"0001-01-01T00:00:00Z","timeEnds":"0001-01-01T00:00:00Z","children":[{"pid":1584,"name":"powershell.exe","description":"powershell.exe","commandLine":"powershell","priority":8,"timeStarted":"2022-02-22T16:51:29.0557047+03:30","timeEnds":"0001-01-01T00:00:00Z","children":[{"pid":2520,"name":"powershell.exe","description":"powershell.exe","commandLine":"\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\"","execPath":"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe","priority":8,"timeStarted":"2022-02-22T16:51:34.8297949+03:30","timeEnds":"2022-02-22T16:52:59.8454381+03:30","exitCode":10,"children":[{"pid":3524,"name":"powershell.exe","description":"powershell.exe","commandLine":"\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\"","execPath":"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe","priority":8,"timeStarted":"2022-02-22T16:51:38.5387946+03:30","timeEnds":"2022-02-22T16:52:54.9775471+03:30","exitCode":5,"children":[{"pid":2252,"name":"powershell.exe","description":"powershell.exe","commandLine":"\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\"","priority":8,"timeStarted":"2022-02-22T16:51:40.5692984+03:30","timeEnds":"2022-02-22T16:52:52.6285369+03:30","exitCode":3,"children":[{"pid":872,"name":"powershell.exe","description":"powershell.exe","commandLine":"\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\"","execPath":"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe","priority":8,"timeStarted":"2022-02-22T16:51:43.4896112+03:30","timeEnds":"2022-02-22T16:52:34.6329949+03:30","children":[{"pid":1320,"name":"explorer.exe","description":"explorer.exe","commandLine":"\"C:\\Windows\\explorer.exe\"  C:\\Windows\\System32\\taskeng.exe","priority":8,"timeStarted":"2022-02-22T16:52:15.8961963+03:30","timeEnds":"2022-02-22T16:52:16.0487871+03:30","exitCode":1}]},{"pid":3116,"name":"cmd.exe","description":"cmd.exe","commandLine":"\"C:\\Windows\\system32\\cmd.exe\"","execPath":"C:\\Windows\\system32\\cmd.exe","priority":8,"timeStarted":"2022-02-22T16:52:34.6046609+03:30","timeEnds":"2022-02-22T16:52:38.5006773+03:30"},{"pid":3100,"name":"cmd.exe","description":"cmd.exe","commandLine":"\"C:\\Windows\\system32\\cmd.exe\"","execPath":"C:\\Windows\\system32\\cmd.exe","priority":8,"timeStarted":"2022-02-22T16:52:39.6237091+03:30","timeEnds":"2022-02-22T16:52:50.4022864+03:30","exitCode":2,"children":[{"pid":604,"name":"cmd.exe","description":"cmd.exe","commandLine":"cmd  \\admin","execPath":"C:\\Windows\\system32\\cmd.exe","priority":8,"timeStarted":"2022-02-22T16:52:43.9328535+03:30","timeEnds":"2022-02-22T16:52:45.6721955+03:30","exitCode":2}]}]}]}]}]}]}`

func main() {
	var jgc *osdat.JsonGraphContainer
	err := json.Unmarshal([]byte(jsonstr), &jgc)
	if err != nil {
		log.Fatalln(err)
	}

	p := flag.Uint("pid", 0, "Find PID")
	flag.Parse()

	found, kind, jgc, chld := jgc.FindByPID(uint32(*p))
	log.Println(found, kind)
	switch kind {
	case 1:
		log.Println(" jgc: ", jgc.Name, jgc.PID)
	case 2:
		log.Println("chld: ", chld.Name, chld.PID)
	}
}

// func findByPID(pid uint32, jgc *osdat.JsonGraphContainer) {
// 	if jgc.PID == pid {
// 		log.Println("GRAND PARENT")
// 		return
// 	}

// 	for _, ch := range jgc.Children {
// 		_findByPID(pid, ch)
// 	}

// }

// func _findByPID(pid uint32, jgc *osdat.JsonGraphChild) {
// 	if jgc.PID == pid {
// 		fmt.Println(jgc.Name, jgc.PID)
// 		return
// 	}
// 	if len(jgc.Children) == 0 {
// 		return
// 	}
// 	for _, ch := range jgc.Children {
// 		_findByPID(pid, ch)
// 	}
// }

// func printRecursive(level int, sublevel int, jgc osdat.JsonGraphChild) {
// 	for t := 0; t < sublevel; t++ {
// 		fmt.Printf("\t")
// 	}
// 	if len(jgc.Children) == 0 {
// 		return
// 	}
// 	fmt.Printf("\t[%d] %s (%d)\n", level, jgc.Name, jgc.PID)
// 	for i, ch := range jgc.Children {
// 		fmt.Printf("\t\t %s (%d):\n", ch.Name, ch.PID)
// 		printRecursive(level+1, i, *ch)
// 	}
// 	fmt.Println()
// }
