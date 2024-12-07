package db

const MAX_N_TASKS = 100

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

// Files are stored on the C2 server
type File_info struct {
	ImplantID string
	FileName  string
	FileSize  int64 //Store the file size in bytes
	FileType  string
	FilePath  string //Store the file location in the filesystem
	CreatedAt int64
}
