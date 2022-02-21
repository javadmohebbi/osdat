//go:build windows
// +build windows

package osdat

import (
	"os"
	"sync"
	"time"

	"github.com/bi-zone/wmi"
)

// WmiMonitorProcessEvents is helping the caller
// to monitor a specific process children
// whether they are succeeded or not
type WmiMonitorProcessEvents struct {

	// this will keep the last child process
	// exit status in case this app is not
	lastExitStatus uint32

	// normally seconds
	// like 1, 2 , 5, ...
	// could be a miliseconds like 0.001, 0.5, 0.1
	within string

	// when all of the processes are done, this will be called
	alldone chan os.Signal

	// os signals like CTRL + D
	sigs chan os.Signal

	// if child exits with code other than zero,
	// this channel will be filled with the error message
	processExitErr chan error

	// err channel
	errCh chan error

	// GrandParnetID is a ID of the
	// Bootstrap application which griAgent run at first
	grandParentID uint32

	// Since some apps will release themselves from
	// their parent, we will keep this IDs to prevent our
	// service from unwanted stop which cause by the menioned release
	monitoredchild []Monitoredchild

	// for lock and unlock
	// when an item appends to or deletes from monitoredchild
	mu sync.Mutex

	// this is designed for createProcess events for
	// wmi notifications. Since it needs to be stopped at the end
	// we need to have variable for that
	qWmiProcessCreate *wmi.NotificationQuery

	// this is designed for process stop trace events for
	// wmi notifications. Since it needs to be stopped at the end
	// we need to have variable for that
	qWmiProcessStopTrace *wmi.NotificationQuery

	// this is a struct field that
	// provide json report about this process
	Graph *JsonGraphContainer `json:"graph"`
}

// this type will monitor the child, grandchild, ... processes
type Monitoredchild struct {

	// first process that starts all childs
	// its a grand grand... parent
	TheFirstParentId uint32

	// child process id
	ProcessID uint32

	// name of the process
	Name string

	// caption of the process
	Caption string

	// process command line options
	CommandLine string

	// exit status
	ExitStatus uint32
}

// this event is uses for win32
type win32ProcessEvent struct {
	Created uint64 `wmi:"TIME_CREATED"`

	// win32process event instance
	Instance win32Process `wmi:"TargetInstance"`
}

// win32_processstoptrace wmi event output
// see: https://docs.microsoft.com/en-us/previous-versions/windows/desktop/krnlprov/win32-processstoptrace
type win32ProcessStopEvent struct {
	Created uint64 `wmi:"TIME_CREATED"`

	System struct {
		Class string
	} `wmi:"Path_"`

	// process id
	ProcessID uint32 `wmi:"ProcessID"`

	// parent process id
	ParentProcessID uint32 `wmi:"ParentProcessID"`

	// exit code
	ExitStatus uint32 `wmi:"ExitStatus"`

	// process name
	ProcessName string `wmi:"ProcessName"`

	// session id
	SessionID uint32 `wmi:"SessionID"`
}

// win32_process wmi event output
// see: https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-process
type win32Process struct {

	// Name of the class or subclass used in the creation of an instance. When used with other key properties of the class, this property allows all instances of the class and its subclasses to be uniquely identified.
	CreationClassName string `json:",omitempty"`

	// Short description of an object—a one-line string.
	Caption string `json:",omitempty"`

	// Command line used to start a specific process, if applicable.
	CommandLine string `json:",omitempty"`

	// Date the process begins executing
	CreationDate time.Time `json:",omitempty"`

	// Creation class name of the scoping computer system
	CSCreationClassName string `json:",omitempty"`

	// Name of the scoping computer system
	CSName string `json:",omitempty"`

	// Description of an object
	Description string `json:",omitempty"`

	// Path to the executable file of the process
	ExecutablePath string `json:",omitempty"`

	// Current operating condition of the process
	// Unknown (0), Other (1), Ready (2), Running (3)
	// Blocked (4), Suspended Blocked (5), Suspended Ready (6)
	// Terminated (7), Stopped (8), Growing (9
	ExecutionState uint16 `json:",omitempty"` //

	// Process identifier. or PID
	Handle string `json:",omitempty"`

	// Total number of open handles owned by the process. HandleCount is the sum of the handles currently open by each thread in this process. A handle is used to examine or modify the system resources. Each handle has an entry in a table that is maintained internally. Entries contain the addresses of the resources and data to identify the resource type
	HandleCount uint32 `json:",omitempty"` //

	// Date an object is installed. The object may be installed without a value being written to this property
	InstallDate time.Time `json:",omitempty"`

	// Time in kernel mode, in milliseconds. If this information is not available, use a value of 0 (zero).
	KernelModeTime uint64 `json:",omitempty"` //

	// Maximum working set size of the process. The working set of a process is the set of memory pages visible to the process in physical RAM. These pages are resident, and available for an application to use without triggering a page fault
	MaximumWorkingSetSize uint32 `json:",omitempty"` //

	// Minimum working set size of the process. The working set of a process is the set of memory pages visible to the process in physical RAM. These pages are resident and available for an application to use without triggering a page fault
	MinimumWorkingSetSize uint32 `json:",omitempty"` //

	// Name of the executable file responsible for the process, equivalent to the Image Name property in Task Manager.
	// When inherited by a subclass, the property can be overridden to be a key property. The name is hard-coded into the application itself and is not affected by changing the file name. For example, even if you rename Calc.exe, the name Calc.exe will still appear in Task Manager and in any WMI scripts that retrieve the process name
	Name string `json:",omitempty"`

	// Creation class name of the scoping operating system.
	OSCreationClassName string `json:",omitempty"`

	// Name of the scoping operating system
	OSName string `json:",omitempty"`

	// Number of I/O operations performed that are not read or write operations
	OtherOperationCount uint64 `json:",omitempty"` //

	// Amount of data transferred during operations that are not read or write operations.
	OtherTransferCount uint64 `json:",omitempty"` //

	// Number of page faults that a process generates
	PageFaults uint32 `json:",omitempty"` //

	// Amount of page file space that a process is using currently. This value is consistent with the VMSize value in TaskMgr.exe
	PageFileUsage uint32 `json:",omitempty"` //

	// Unique identifier of the process that creates a process. Process identifier numbers are reused, so they only identify a process for the lifetime
	// of that process. It is possible that the process identified by ParentProcessId is terminated, so ParentProcessId may not refer to a running process.
	// It is also possible that ParentProcessId incorrectly refers to a process that reuses a process identifier. You can use the CreationDate property to
	// determine whether the specified parent was created after the process represented by this Win32_Process instance was created
	ParentProcessId uint32 `json:",omitempty"` //

	// Maximum amount of page file space used during the life of a process.
	PeakPageFileUsage uint32 `json:",omitempty"` //

	// Maximum virtual address space a process uses at any one time. Using virtual address space does not necessarily imply corresponding use of either disk
	// or main memory pages. However, virtual space is finite, and by using too much the process might not be able to load libraries
	PeakVirtualSize uint64 `json:",omitempty"` //

	// Peak working set size of a process
	PeakWorkingSetSize uint32 `json:",omitempty"` //

	// Scheduling priority of a process within an operating system. The higher the value, the higher priority a process receives. Priority values can range from 0 (zero), which is the lowest priority to 31, which is highest priority
	Priority *uint32 `json:",omitempty"` //

	// Current number of pages allocated that are only accessible to the process represented by this Win32_Process instance
	PrivatePageCount uint64 `json:",omitempty"` //

	// Numeric identifier used to distinguish one process from another. ProcessIDs are valid from process creation time to process termination. Upon termination, that same numeric identifier can be applied to a new process
	// This means that you cannot use ProcessID alone to monitor a particular process. For example, an application could have a ProcessID of 7, and then fail. When a new process is started, the new process could be assigned ProcessID 7.
	// A script that checked only for a specified ProcessID could thus be "fooled" into thinking that the original application was still running.
	ProcessId uint32 `json:",omitempty"` //

	// Quota amount of nonpaged pool usage for a process
	QuotaNonPagedPoolUsage uint32 `json:",omitempty"` //

	// Quota amount of paged pool usage for a process
	QuotaPagedPoolUsage uint32 `json:",omitempty"` //

	// Peak quota amount of nonpaged pool usage for a process
	QuotaPeakNonPagedPoolUsage uint32 `json:",omitempty"` //

	// Peak quota amount of paged pool usage for a process
	QuotaPeakPagedPoolUsage uint32 `json:",omitempty"` //

	// Number of read operations performed
	ReadOperationCount uint64 `json:",omitempty"` //

	// Amount of data read
	ReadTransferCount uint64 `json:",omitempty"` //

	// Unique identifier that an operating system generates when a session is created. A session spans a period of time from logon until logoff from a specific system
	SessionId uint32 `json:",omitempty"` //

	// This property is not implemented and does not get populated for any instance of this class. It is always NULL
	// Values include the following: OK ("OK"), Error ("Error"), Degraded ("Degraded"),
	// Unknown ("Unknown"), Pred Fail ("Pred Fail"), Starting ("Starting"),
	// Stopping ("Stopping"), Service ("Service"), Stressed ("Stressed"),
	// NonRecover ("NonRecover"), No Contact ("No Contact"), Lost Comm ("Lost Comm")
	Status string `json:",omitempty"`

	// Process was stopped or terminated. To get the termination time, a handle to the process must be held open. Otherwise, this property returns NULL.
	TerminationDate time.Time `json:",omitempty"`

	// Number of active threads in a process. An instruction is the basic unit of execution in a processor, and a thread is the object that executes an instruction. Each running process has at least one thread
	ThreadCount uint32 `json:",omitempty"` //

	// Time in user mode, in 100 nanosecond units. If this information is not available, use a value of 0 (zero).
	UserModeTime uint64 `json:",omitempty"` //

	// Current size of the virtual address space that a process is using, not the physical or virtual memory actually used by the process.
	// Using virtual address space does not necessarily imply corresponding use of either disk or main memory pages
	// Virtual space is finite, and by using too much, the process might not be able to load libraries. This value is consistent with what you see in Perfmon.exe.
	VirtualSize uint64 `json:",omitempty"` //

	// Version of Windows in which the process is running
	WindowsVersion string `json:",omitempty"`

	// Amount of memory in bytes that a process needs to execute efficiently—for an operating system that uses page-based memory management. If the system does not have enough memory (less than the working set size), thrashing occurs.
	// If the size of the working set is not known, use NULL or 0 (zero). If working set data is provided, you can monitor the information to understand the changing memory requirements of a process.
	WorkingSetSize uint64 `json:",omitempty"` //

	// Number of write operations performed.
	WriteOperationCount uint64 `json:",omitempty"` //

	// Amount of data writte
	WriteTransferCount uint64 `json:",omitempty"` //

}
