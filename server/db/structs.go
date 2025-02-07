package db

const MAX_N_TASKS = 100

// Task types
const (
	ExecCommand uint8 = iota + 200
	ExecShellcode
	SendFileToImplant
	GetFileFromImplant
	ExecInlineAsm
	KillImplant
)

var TaskTypes = map[uint8]string{
	ExecCommand:        "EXEC_COMMAND",
	ExecShellcode:      "EXEC_SHELLCODE",
	SendFileToImplant:  "SEND_FILE_TO_IMPLANT",
	GetFileFromImplant: "GET_FILE_FROM_IMPLANT",
	ExecInlineAsm:      "EXEC_INLINE_ASM",
	KillImplant:        "KILL_IMPLANT",
}

type ImplantInfo struct {
	ID          string
	Hostname    string
	IntIP       string
	ExtIP       string
	Os          string
	ProcessID   uint64
	ProcessUser string
	ProtName    string // Name of the EDR/AV present on the machine
	LastCheckIn int64
	Active      bool
	KillDate    int64
}

type ImplantTask struct {
	TaskID      string
	ImplantID   string
	FileID      uint64
	TaskType    uint8
	TaskData    []byte
	CreatedAt   int64
	Completed   bool
	CompletedAt int64
	TaskResult  []byte
}

type ListenerInfo struct {
	ListenerID string
	Config     []byte
	Host       string
	Port       uint16
	CreatedAt  int64
	KillDate   int64
}

// FileInfo Files are stored on the C2 server local filesystem
type FileInfo struct {
	ImplantID string
	FileName  string
	FileSize  int64 // Store the file size in bytes
	FileType  string
	FilePath  string // Store the file location in the filesystem
	CreatedAt int64
}
