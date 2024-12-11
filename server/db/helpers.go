package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dovelus/corb4n-c2/server/comunication"
)

// Get all active implants
func GetAllImplants() ([]*Implant_info, error) {
	comunication.Logger.Info("Getting all implants: SELECT * FROM implants")
	statement, err := dbConn.Prepare("SELECT * FROM implants")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Error("Error querying database: ", err)
		return nil, err
	}

	var implants []*Implant_info
	for rows.Next() {
		Info := new(Implant_info)
		err := rows.Scan(
			&Info.ID,
			&Info.Hostname,
			&Info.IntIP,
			&Info.ExtIP,
			&Info.Os,
			&Info.ProcessID,
			&Info.ProcessUser,
			&Info.ProtName,
			&Info.LastCheckIn,
			&Info.Active,
			&Info.KillDate)
		if err != nil {
			comunication.Logger.Error("Error scanning row: ", err)
			return nil, err
		}
		implants = append(implants, Info)
	}
	if len(implants) == 0 {
		return nil, comunication.ErrNoResults
	}

	return implants, nil
}

// Get specific implant informatio for a given ID
func GetImplantByID(ID string) (*Implant_info, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM implants WHERE implant_id = '%s'", ID))
	statement, err := dbConn.Prepare("SELECT * FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(ID)
	Info := new(Implant_info)
	err = row.Scan(
		&Info.ID,
		&Info.Hostname,
		&Info.IntIP,
		&Info.ExtIP,
		&Info.Os,
		&Info.ProcessID,
		&Info.ProcessUser,
		&Info.ProtName,
		&Info.LastCheckIn,
		&Info.Active,
		&Info.KillDate)
	if err != nil {
		comunication.Logger.Error("Error scanning row: ", err)
		return nil, err
	}

	return Info, nil
}

func AddImplant(info *Implant_info) error {
	// Check if the implant already exists
	comunication.Logger.Info(fmt.Sprintf("Checking if implant with ID '%s' already exists", info.ID))
	statement, err := dbConn.Prepare("SELECT COUNT(*) FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	var count int
	err = statement.QueryRow(info.ID).Scan(&count)
	if err != nil {
		comunication.Logger.Error("Error querying database: ", err)
		return err
	}

	if count > 0 {
		comunication.Logger.Warn(fmt.Sprintf("Implant with ID '%s' already exists", info.ID))
		return comunication.ErrImplantExists
	}

	// Add the implant to the database
	comunication.Logger.Info(fmt.Sprintf("INSERT INTO implants VALUES ('%s', '%s', '%s', '%s', '%s', %d, '%s', '%s', %d, %t, %d)",
		info.ID,
		info.Hostname,
		info.IntIP,
		info.ExtIP,
		info.Os,
		info.ProcessID,
		info.ProcessUser,
		info.ProtName,
		info.LastCheckIn,
		info.Active,
		info.KillDate))

	statement, err = dbConn.Prepare("INSERT INTO implants VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		info.ID,
		info.Hostname,
		info.IntIP,
		info.ExtIP,
		info.Os,
		info.ProcessID,
		info.ProcessUser,
		info.ProtName,
		info.LastCheckIn,
		info.Active,
		info.KillDate)

	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	return nil
}

// Given an implant ID, removes all information from the database (implant, task)
func RemoveImplant(ID string) error {
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM implants WHERE implant_id = '%s'", ID))
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM tasks WHERE implant_id = '%s'", ID))

	statement_implants, err := dbConn.Prepare("DELETE FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for implants: ", err)
		return err
	}
	defer statement_implants.Close()

	statement_tasks, err := dbConn.Prepare("DELETE FROM tasks WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for tasks: ", err)
		return err
	}
	defer statement_tasks.Close()

	res_implants, err := statement_implants.Exec(ID)
	if err != nil {
		comunication.Logger.Error("Error executing query for implants: ", err)
		return err
	}

	res_tasks, err := statement_tasks.Exec(ID)
	if err != nil {
		comunication.Logger.Error("Error executing query for tasks: ", err)
		return err
	}

	implants_affected, err := res_implants.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
		return err
	}
	tasks_affected, err := res_tasks.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
		return err
	}

	if implants_affected == 0 && tasks_affected == 0 {
		comunication.Logger.Warn("No rows affected")
		return comunication.ErrNoResults
	}

	return nil
}

// Given an implant ID, sets the active field to false
func SetImplantStatus(ID string, status bool) error {
	comunication.Logger.Info(fmt.Sprintf("UPDATE implants SET active = %t WHERE implant_id = '%s'", status, ID))
	statement, err := dbConn.Prepare("UPDATE implants SET active = ? WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for implants: ", err)
	}
	defer statement.Close()

	res, err := statement.Exec(status, ID)
	if err != nil {
		comunication.Logger.Error("Error executing query for implants: ", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
		return err
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return comunication.ErrNoResults
	}

	return nil
}

// Given an implant ID, returns the status of the implant
func GetImplantStatus(ID string) (bool, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT active FROM implants WHERE implant_id= '%s'", ID))
	statement, err := dbConn.Prepare("SELECT active FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for implants: ", err)
		return false, err
	}
	defer statement.Close()

	row := statement.QueryRow(ID)
	var status bool
	err = row.Scan(&status)
	if err != nil {
		comunication.Logger.Error("Error scanning row: ", err)
		return false, err
	}

	return status, nil
}

// Given an implant ID, sets the killDate field to the current time
func UpdateImplantKillDate(ID string) error {
	var kill_time int64 = time.Now().Unix()
	comunication.Logger.Info(fmt.Sprintf("UPDATE implants SET kill_date = %d WHERE implant_id = '%s'", kill_time, ID))
	statement, err := dbConn.Prepare("UPDATE implants SET kill_date = ? WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for implants: ", err)
	}
	defer statement.Close()

	res, err := statement.Exec(kill_time, ID)
	if err != nil {
		comunication.Logger.Error("Error executing query for implants: ", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
		return err
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return comunication.ErrNoResults
	}

	return nil
}

// Given an implant ID, returns all its tasks
func GetImplantTasks(ID string, completed bool) ([]*Implant_Task, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM tasks WHERE implant_id='%s' AND completed=%t", ID, completed))
	var statement *sql.Stmt
	var err error
	if completed {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE implant_id = ? AND completed = TRUE")
	} else {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE implant_id = ? AND completed = FALSE")
	}
	if err != nil {
		comunication.Logger.Error("Error preparing statement for tasks: ", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(ID)
	if err != nil {
		comunication.Logger.Error("Error executing query for tasks: ", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*Implant_Task = make([]*Implant_Task, MAX_N_TASKS)
	for rows.Next() {
		task := new(Implant_Task)
		err := rows.Scan(
			&task.TaskID,
			&task.ImplantID,
			&task.TaskType,
			&task.TaskData,
			&task.CreatedAt,
			&task.Completed,
			&task.CompletedAt,
			&task.TaskResult)
		if err != nil {
			comunication.Logger.Error("Error scanning row: ", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, comunication.ErrNoResults
	}

	return tasks, nil
}

// Returns all tasks for all implants
func GetAllTasks(completed bool) ([]*Implant_Task, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM tasks WHERE completed=%t", completed))
	var statement *sql.Stmt
	var err error
	if completed {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE completed = TRUE")
	} else {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE completed = FALSE")
	}
	if err != nil {
		comunication.Logger.Error("Error preparing statement for tasks: ", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Error("Error executing query for tasks: ", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []*Implant_Task = make([]*Implant_Task, MAX_N_TASKS)
	for rows.Next() {
		task := new(Implant_Task)
		err := rows.Scan(
			&task.TaskID,
			&task.ImplantID,
			&task.TaskType,
			&task.TaskData,
			&task.CreatedAt,
			&task.Completed,
			&task.CompletedAt,
			&task.TaskResult)
		if err != nil {
			comunication.Logger.Error("Error scanning row: ", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, comunication.ErrNoResults
	}

	return tasks, nil
}

// Adds a task to the database
func AddTask(task *Implant_Task) error {
	comunication.Logger.Info(fmt.Sprintf("INSERT INTO tasks VALUES ('%s', '%s', %d, '%s', %d, %t, %d, '%s')",
		task.TaskID,
		task.ImplantID,
		task.TaskType,
		task.TaskData,
		task.CreatedAt,
		task.Completed,
		task.CompletedAt,
		task.TaskResult))

	statement, err := dbConn.Prepare("INSERT INTO tasks VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		task.TaskID,
		task.ImplantID,
		task.TaskType,
		task.TaskData,
		task.CreatedAt,
		task.Completed,
		task.CompletedAt,
		task.TaskResult)

	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	return nil
}

// Removes non completed tasks for a given implantid and taskid
func RemovePendingTasksImplant(ID string, taskID string) error {
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM tasks WHERE implant_id = '%s' AND task_id = '%s' AND completed = FALSE", ID, taskID))
	statement, err := dbConn.Prepare("DELETE FROM tasks WHERE implant_id = ? AND task_id = ? AND completed = FALSE")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(ID, taskID)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected for tasks: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No task with ID '%s' for implant '%s'", taskID, ID))
	}

	return nil
}

func GetTask(taskID string) (*Implant_Task, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM tasks WHERE task_id = '%s'", taskID))
	statement, err := dbConn.Prepare("SELECT * FROM tasks WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for tasks: ", err)
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(taskID)
	task := new(Implant_Task)
	err = row.Scan(
		&task.TaskID,
		&task.ImplantID,
		&task.TaskType,
		&task.TaskData,
		&task.CreatedAt,
		&task.Completed,
		&task.CompletedAt,
		&task.TaskResult)

	if err != nil {
		if err == sql.ErrNoRows {
			comunication.Logger.Error("No task found with the given taskID: ", taskID)
			return nil, comunication.ErrNoResults
		}
		comunication.Logger.Error("Error scanning row: ", err)
		return nil, err
	}

	return task, nil
}

// Remove task based on taskID
func RemoveTask(taskID string) error {
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM tasks WHERE task_id = '%s'", taskID))
	statement, err := dbConn.Prepare("DELETE FROM tasks WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for tasks: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(taskID)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected for tasks: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No database entry with ID '%s'", taskID))
	}

	return nil
}

// Change status of tast to completed
func CompleteTask(taskID string) error {
	var complete_time int64 = time.Now().Unix()
	comunication.Logger.Info(fmt.Sprintf("UPDATE tasks SET completed = TRUE, completed_at = %d WHERE task_id = '%s'", complete_time, taskID))
	statement, err := dbConn.Prepare("UPDATE tasks SET completed = TRUE, completed_at = ? WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for tasks: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(complete_time, taskID)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected for tasks: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No task with ID '%s'", taskID))
	}

	return nil
}

// Updates last checkin time for an implant
func UpdateImplantCheckin(ID string) error {
	var checkin_time int64 = time.Now().Unix()
	comunication.Logger.Info(fmt.Sprintf("UPDATE implants SET last_checkin = %d WHERE implant_id = '%s'", checkin_time, ID))
	statement, err := dbConn.Prepare("UPDATE implants SET last_check_in = ? WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement for implants: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(checkin_time, ID)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected for implants: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No implant with ID '%s'", ID))
	}

	return nil
}

// Get all listeners
func GetAllListeners() ([]*Listener_info, error) {
	comunication.Logger.Info("SELECT * FROM listeners")
	statement, err := dbConn.Prepare("SELECT * FROM listeners")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Error("Error querying database: ", err)
		return nil, err
	}

	var listeners []*Listener_info
	for rows.Next() {
		Info := new(Listener_info)
		err := rows.Scan(
			&Info.ListenerID,
			&Info.Config,
			&Info.Host,
			&Info.Port,
			&Info.CreatedAt,
			&Info.KillDate)
		if err != nil {
			comunication.Logger.Error("Error scanning row: ", err)
			return nil, err
		}
		listeners = append(listeners, Info)
	}
	if len(listeners) == 0 {
		return nil, comunication.ErrNoResults
	}

	return listeners, nil
}

// Get listener by ID
func GetListenerByID(ID string) (*Listener_info, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM listeners WHERE listener_id = '%s'", ID))
	statement, err := dbConn.Prepare("SELECT * FROM listeners WHERE listener_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(ID)
	Info := new(Listener_info)
	err = row.Scan(
		&Info.ListenerID,
		&Info.Config,
		&Info.Host,
		&Info.Port,
		&Info.CreatedAt,
		&Info.KillDate)
	if err != nil {
		comunication.Logger.Error("Error scanning row: ", err)
		return nil, err
	}

	return Info, nil
}

// Add a listener to the database
func AddListener(info *Listener_info) error {
	comunication.Logger.Info(fmt.Sprintf("INSERT INTO listeners VALUES ('%s', '%s', '%s', %d, %d, %d)",
		info.ListenerID,
		info.Config,
		info.Host,
		info.Port,
		info.CreatedAt,
		info.KillDate))

	statement, err := dbConn.Prepare("INSERT INTO listeners VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		info.ListenerID,
		info.Config,
		info.Host,
		info.Port,
		info.CreatedAt,
		info.KillDate)

	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	return nil
}

// Remove listener by ID
func RemoveListener(ID string) error {
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM listeners WHERE listener_id = '%s'", ID))
	statement, err := dbConn.Prepare("DELETE FROM listeners WHERE listener_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(ID)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No listener with ID '%s'", ID))
	}

	return nil
}

// Updates listener kill date
func UpdateListenerKillDate(ID string) error {
	var kill_time int64 = time.Now().Unix()
	comunication.Logger.Info(fmt.Sprintf("UPDATE listeners SET kill_date = %d WHERE listener_id = '%s'", kill_time, ID))
	statement, err := dbConn.Prepare("UPDATE listeners SET kill_date = ? WHERE listener_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(kill_time, ID)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No listener with ID '%s'", ID))
	}

	return nil
}

// Add file to the database
func AddFile(info *File_info) error {
	comunication.Logger.Info(fmt.Sprintf("INSERT INTO files (implant_id, file_path, file_name, file_type, file_size, created_at) VALUES ('%s', '%s', '%s', '%s', %d, %d)",
		info.ImplantID,
		info.FilePath,
		info.FileName,
		info.FileType,
		info.FileSize,
		info.CreatedAt))

	statement, err := dbConn.Prepare("INSERT INTO files (implant_id, file_path, file_name, file_type, file_size, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		info.ImplantID,
		info.FilePath,
		info.FileName,
		info.FileType,
		info.FileSize,
		info.CreatedAt)

	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	return nil
}

// Get all files
func GetAllFiles() ([]*File_info, error) {
	comunication.Logger.Info("SELECT * FROM files")
	statement, err := dbConn.Prepare("SELECT * FROM files")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Error("Error querying database: ", err)
		return nil, err
	}

	var files []*File_info
	for rows.Next() {
		Info := new(File_info)
		var id int64 // Temporary variable to hold the id
		err := rows.Scan(
			&id,
			&Info.ImplantID,
			&Info.FilePath,
			&Info.FileName,
			&Info.FileType,
			&Info.FileSize,
			&Info.CreatedAt)
		if err != nil {
			comunication.Logger.Error("Error scanning row: ", err)
			return nil, err
		}
		files = append(files, Info)
	}
	if len(files) == 0 {
		return nil, comunication.ErrNoResults
	}

	return files, nil
}

// Get file by implant ID
func GetFilesByImplantID(ID string) ([]*File_info, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM files WHERE implant_id = '%s'", ID))
	statement, err := dbConn.Prepare("SELECT * FROM files WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(ID)
	if err != nil {
		comunication.Logger.Error("Error querying database: ", err)
		return nil, err
	}

	var files []*File_info
	for rows.Next() {
		Info := new(File_info)
		var id int64 // Temporary variable to hold the id
		err := rows.Scan(
			&id,
			&Info.ImplantID,
			&Info.FilePath,
			&Info.FileName,
			&Info.FileType,
			&Info.FileSize,
			&Info.CreatedAt)
		if err != nil {
			comunication.Logger.Error("Error scanning row: ", err)
			return nil, err
		}
		files = append(files, Info)
	}
	if len(files) == 0 {
		return nil, comunication.ErrNoResults
	}

	return files, nil
}

// Get file by implant ID and file name
func GetFileByImplantIDAndName(ID string, name string) (*File_info, error) {
	comunication.Logger.Info(fmt.Sprintf("SELECT * FROM files WHERE implant_id = '%s' AND file_name = '%s'", ID, name))
	statement, err := dbConn.Prepare("SELECT * FROM files WHERE implant_id = ? AND file_name = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(ID, name)
	var id int64 // Temporary variable to hold the id
	Info := new(File_info)
	err = row.Scan(
		&id,
		&Info.ImplantID,
		&Info.FilePath,
		&Info.FileName,
		&Info.FileType,
		&Info.FileSize,
		&Info.CreatedAt)
	if err != nil {
		comunication.Logger.Error("Error scanning row: ", err)
		return nil, err
	}

	return Info, nil
}

// Remove file by implant ID and file
func RemoveFile(ID string, name string) error {
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM files WHERE implant_id = '%s' AND file_name = '%s'", ID, name))
	statement, err := dbConn.Prepare("DELETE FROM files WHERE implant_id = ? AND file_name = ?")
	if err != nil {
		comunication.Logger.Error("Error preparing statement: ", err)
		return err
	}
	defer statement.Close()

	resp, err := statement.Exec(ID, name)
	if err != nil {
		comunication.Logger.Error("Error executing query: ", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Error("Error getting rows affected: ", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No file with name '%s' for implant '%s'", name, ID))
	}

	return nil
}
