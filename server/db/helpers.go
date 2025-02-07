package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dovelus/corb4n-c2/server/comunication"
)

// GetAllImplants Get all active implants
func GetAllImplants() ([]*ImplantInfo, error) {
	comunication.Logger.Info("Getting all implants: SELECT * FROM implants")
	statement, err := dbConn.Prepare("SELECT * FROM implants")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Errorf("Error querying database: %v", err)
		return nil, err
	}

	var implants []*ImplantInfo
	for rows.Next() {
		Info := new(ImplantInfo)
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
			comunication.Logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		implants = append(implants, Info)
	}
	if len(implants) == 0 {
		return nil, comunication.ErrNoResults
	}
	return implants, nil
}

// GetImplantByID Get specific implant information for a given ID
func GetImplantByID(ID string) (*ImplantInfo, error) {
	comunication.Logger.Infof("SELECT * FROM implants WHERE implant_id = '%s'", ID)
	statement, err := dbConn.Prepare("SELECT * FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	row := statement.QueryRow(ID)
	Info := new(ImplantInfo)
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
		comunication.Logger.Errorf("Error scanning row: %v", err)
		return nil, err
	}

	return Info, nil
}

func AddImplant(info *ImplantInfo) error {
	// Check if the implant already exists
	comunication.Logger.Infof("Checking if implant with ID '%s' already exists", info.ID)
	statement, err := dbConn.Prepare("SELECT COUNT(*) FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	var count int
	err = statement.QueryRow(info.ID).Scan(&count)
	if err != nil {
		comunication.Logger.Errorf("Error querying database: %v", err)
		return err
	}

	if count > 0 {
		comunication.Logger.Warnf("Implant with ID '%s' already exists", info.ID)
		return comunication.ErrImplantExists
	}

	// Add the implant to the database
	comunication.Logger.Infof("INSERT INTO implants VALUES ('%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%d', '%t', '%d')",
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

	statement, err = dbConn.Prepare("INSERT INTO implants VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

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
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	return nil
}

// RemoveImplant Given an implant ID, removes all information from the database (implant, task)
func RemoveImplant(ID string) error {
	comunication.Logger.Infof("DELETE FROM implants WHERE implant_id = '%s'", ID)
	comunication.Logger.Infof("DELETE FROM tasks WHERE implant_id = '%s'", ID)

	statementImplants, err := dbConn.Prepare("DELETE FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for implants: %v", err)
		return err
	}
	defer func(statementImplants *sql.Stmt) {
		err := statementImplants.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for implants: %v", err)
		}
	}(statementImplants)

	statementTasks, err := dbConn.Prepare("DELETE FROM tasks WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return err
	}
	defer func(statementTasks *sql.Stmt) {
		err := statementTasks.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statementTasks)

	resImplants, err := statementImplants.Exec(ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query for implants: %v", err)
		return err
	}

	resTasks, err := statementTasks.Exec(ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query for tasks: %v", err)
		return err
	}

	implantsAffected, err := resImplants.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
		return err
	}
	tasksAffected, err := resTasks.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
		return err
	}

	if implantsAffected == 0 && tasksAffected == 0 {
		comunication.Logger.Warn("No rows affected")
		return comunication.ErrNoResults
	}

	return nil
}

// SetImplantStatus Given an implant ID, sets the active field to false
func SetImplantStatus(ID string, status bool) error {
	comunication.Logger.Infof("UPDATE implants SET active = '%t' WHERE implant_id = '%s'", status, ID)
	statement, err := dbConn.Prepare("UPDATE implants SET active = ? WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for implants: %v", err)
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for implants: %v", err)
		}
	}(statement)

	res, err := statement.Exec(status, ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query for implants: %v", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
		return err
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return comunication.ErrNoResults
	}

	return nil
}

// GetImplantStatus Given an implant ID, returns the status of the implant
func GetImplantStatus(ID string) (bool, error) {
	comunication.Logger.Info("SELECT active FROM implants WHERE implant_id= '%s'", ID)
	statement, err := dbConn.Prepare("SELECT active FROM implants WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for implants: %v", err)
		return false, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for implants: %v", err)
		}
	}(statement)

	row := statement.QueryRow(ID)
	var status bool
	err = row.Scan(&status)
	if err != nil {
		comunication.Logger.Errorf("Error scanning row: %v", err)
		return false, err
	}

	return status, nil
}

// UpdateImplantKillDate Given an implant ID, sets the killDate field to the current time
func UpdateImplantKillDate(ID string) error {
	var killTime = time.Now().Unix()
	comunication.Logger.Infof("UPDATE implants SET kill_date = '%d' WHERE implant_id = '%s'", killTime, ID)
	statement, err := dbConn.Prepare("UPDATE implants SET kill_date = ? WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for implants: %v", err)
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for implants: %v", err)
		}
	}(statement)

	res, err := statement.Exec(killTime, ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query for implants: %v", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
		return err
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return comunication.ErrNoResults
	}

	return nil
}

// GetImplantTasks Given an implant ID, returns all its tasks
func GetImplantTasks(ID string, completed bool) ([]*ImplantTask, error) {
	comunication.Logger.Infof("SELECT * FROM tasks WHERE implant_id='%s' AND completed='%t'", ID, completed)

	// Check the total number of tasks in the database
	var totalTasks int
	err := dbConn.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&totalTasks)
	if err != nil {
		comunication.Logger.Errorf("Error counting tasks: %v", err)
		return nil, err
	}

	if totalTasks > MAX_N_TASKS {
		return nil, fmt.Errorf("total number of tasks exceeds the maximum limit of %d", MAX_N_TASKS)
	}

	var statement *sql.Stmt
	if completed {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE implant_id = ? AND completed = TRUE")
	} else {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE implant_id = ? AND completed = FALSE")
	}
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statement)

	rows, err := statement.Query(ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query for tasks: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing rows: %v", err)
		}
	}(rows)

	var tasks []*ImplantTask
	for rows.Next() {
		task := new(ImplantTask)
		err := rows.Scan(
			&task.TaskID,
			&task.ImplantID,
			&task.FileID,
			&task.TaskType,
			&task.TaskData,
			&task.CreatedAt,
			&task.Completed,
			&task.CompletedAt,
			&task.TaskResult)
		if err != nil {
			comunication.Logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, comunication.ErrNoResults
	}

	return tasks, nil
}

// GetAllTasks Returns all tasks for all implants
func GetAllTasks(completed bool) ([]*ImplantTask, error) {
	comunication.Logger.Infof("SELECT * FROM tasks WHERE completed='%t'", completed)
	var statement *sql.Stmt
	var err error
	if completed {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE completed = TRUE")
	} else {
		statement, err = dbConn.Prepare("SELECT * FROM tasks WHERE completed = FALSE")
	}
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statement)

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Errorf("Error executing query for tasks: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing rows: %v", err)
		}
	}(rows)

	var tasks = make([]*ImplantTask, MAX_N_TASKS)
	for rows.Next() {
		task := new(ImplantTask)
		err := rows.Scan(
			&task.TaskID,
			&task.ImplantID,
			&task.FileID,
			&task.TaskType,
			&task.TaskData,
			&task.CreatedAt,
			&task.Completed,
			&task.CompletedAt,
			&task.TaskResult)
		if err != nil {
			comunication.Logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil, comunication.ErrNoResults
	}

	return tasks, nil
}

// AddTask Adds a task to the database
func AddTask(task *ImplantTask) error {
	comunication.Logger.Infof("INSERT INTO tasks VALUES ('%s', '%s', '%s' ,'%d', '%s', '%d', '%t', '%d', '%s')",
		task.TaskID,
		task.ImplantID,
		task.FileID,
		task.TaskType,
		task.TaskData,
		task.CreatedAt,
		task.Completed,
		task.CompletedAt,
		task.TaskResult)

	statement, err := dbConn.Prepare("INSERT INTO tasks VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	_, err = statement.Exec(
		task.TaskID,
		task.ImplantID,
		task.FileID,
		task.TaskType,
		task.TaskData,
		task.CreatedAt,
		task.Completed,
		task.CompletedAt,
		task.TaskResult)

	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	return nil
}

// RemovePendingTasksImplant Removes non completed tasks for a given implant-id and taskid
func RemovePendingTasksImplant(ID string) error {
	comunication.Logger.Info(fmt.Sprintf("DELETE FROM tasks WHERE implant_id = '%s' AND completed = FALSE", ID))
	statement, err := dbConn.Prepare("DELETE FROM tasks WHERE implant_id = ? AND completed = FALSE")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected for tasks: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No task with implant-ID: '%s'", ID))
	}

	return nil
}

// RemoveAllTasksImplant Removes all tasks for a given implant-id
func RemoveAllTasksImplant(ID string) error {
	comunication.Logger.Infof("DELETE FROM tasks WHERE implant_id = '%s'", ID)
	statement, err := dbConn.Prepare("DELETE FROM tasks WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected for tasks: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No task with implant-ID: '%s'", ID))
	}

	return nil
}

func GetTask(taskID string) (*ImplantTask, error) {
	comunication.Logger.Info("SELECT * FROM tasks WHERE task_id = '%s'", taskID)
	statement, err := dbConn.Prepare("SELECT * FROM tasks WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statement)

	row := statement.QueryRow(taskID)
	task := new(ImplantTask)
	err = row.Scan(
		&task.TaskID,
		&task.ImplantID,
		&task.FileID,
		&task.TaskType,
		&task.TaskData,
		&task.CreatedAt,
		&task.Completed,
		&task.CompletedAt,
		&task.TaskResult)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comunication.Logger.Errorf("No task found with the given taskID: %v", taskID)
			return nil, comunication.ErrNoResults
		}
		comunication.Logger.Errorf("Error scanning row: %v", err)
		return nil, err
	}

	return task, nil
}

// RemoveTask Remove task based on taskID
func RemoveTask(taskID string) error {
	comunication.Logger.Infof("DELETE FROM tasks WHERE task_id = '%s'", taskID)
	statement, err := dbConn.Prepare("DELETE FROM tasks WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(taskID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected for tasks: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No database entry with ID '%s'", taskID))
	}

	return nil
}

func CompleteTask(taskID string, taskResult []byte) error {
	var completeTime = time.Now().Unix()
	comunication.Logger.Infof("UPDATE tasks SET completed = TRUE, completed_at = '%d', task_result = ? WHERE task_id = '%s'", completeTime, taskID)
	statement, err := dbConn.Prepare("UPDATE tasks SET completed = TRUE, completed_at = ?, task_result = ? WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(completeTime, taskResult, taskID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected for tasks: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No task with ID '%s'", taskID))
	}

	return nil
}

// CompleteTaskWithFile Completes a task with a file
func CompleteTaskWithFile(taskID string, fileID int64) error {
	var completeTime = time.Now().Unix()
	comunication.Logger.Infof("UPDATE tasks SET completed = TRUE, completed_at = '%d', file_id = '%d' WHERE task_id = '%s'", completeTime, fileID, taskID)
	statement, err := dbConn.Prepare("UPDATE tasks SET completed = TRUE, completed_at = ?, file_id = ? WHERE task_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for tasks: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for tasks: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(completeTime, fileID, taskID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected for tasks: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No task with ID '%s'", taskID))
	}

	return nil
}

// UpdateImplantCheckin Updates last checkin time for an implant
func UpdateImplantCheckin(ID string) error {
	var checkinTime = time.Now().Unix()
	comunication.Logger.Infof("UPDATE implants SET last_check_in = '%d' WHERE implant_id = '%s'", checkinTime, ID)
	statement, err := dbConn.Prepare("UPDATE implants SET last_check_in = ? WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement for implants: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement for implants: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(checkinTime, ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected for implants: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No implant with ID '%s'", ID))
	}

	return nil
}

// GetAllListeners Get all listeners
func GetAllListeners() ([]*ListenerInfo, error) {
	comunication.Logger.Info("SELECT * FROM listeners")
	statement, err := dbConn.Prepare("SELECT * FROM listeners")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Errorf("Error querying database: %v", err)
		return nil, err
	}

	var listeners []*ListenerInfo
	for rows.Next() {
		Info := new(ListenerInfo)
		err := rows.Scan(
			&Info.ListenerID,
			&Info.Config,
			&Info.Host,
			&Info.Port,
			&Info.CreatedAt,
			&Info.KillDate)
		if err != nil {
			comunication.Logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		listeners = append(listeners, Info)
	}
	if len(listeners) == 0 {
		return nil, comunication.ErrNoResults
	}

	return listeners, nil
}

// GetListenerByID Get listener by ID
func GetListenerByID(ID string) (*ListenerInfo, error) {
	comunication.Logger.Infof("SELECT * FROM listeners WHERE listener_id = '%s'", ID)
	statement, err := dbConn.Prepare("SELECT * FROM listeners WHERE listener_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	row := statement.QueryRow(ID)
	Info := new(ListenerInfo)
	err = row.Scan(
		&Info.ListenerID,
		&Info.Config,
		&Info.Host,
		&Info.Port,
		&Info.CreatedAt,
		&Info.KillDate)
	if err != nil {
		comunication.Logger.Errorf("Error scanning row: %v", err)
		return nil, err
	}

	return Info, nil
}

// AddListener Add a listener to the database
func AddListener(info *ListenerInfo) error {
	comunication.Logger.Infof("INSERT INTO listeners VALUES ('%s', '%s', '%s', '%d', '%d', '%d')",
		info.ListenerID,
		info.Config,
		info.Host,
		info.Port,
		info.CreatedAt,
		info.KillDate)

	statement, err := dbConn.Prepare("INSERT INTO listeners VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	_, err = statement.Exec(
		info.ListenerID,
		info.Config,
		info.Host,
		info.Port,
		info.CreatedAt,
		info.KillDate)

	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	return nil
}

// RemoveListener Remove listener by ID
func RemoveListener(ID string) error {
	comunication.Logger.Infof("DELETE FROM listeners WHERE listener_id = '%s'", ID)
	statement, err := dbConn.Prepare("DELETE FROM listeners WHERE listener_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No listener with ID '%s'", ID))
	}

	return nil
}

// UpdateListenerKillDate Updates listener kill date
func UpdateListenerKillDate(ID string) error {
	var killTime = time.Now().Unix()
	comunication.Logger.Infof("UPDATE listeners SET kill_date = '%d' WHERE listener_id = '%s'", killTime, ID)
	statement, err := dbConn.Prepare("UPDATE listeners SET kill_date = ? WHERE listener_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(killTime, ID)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No listener with ID '%s'", ID))
	}

	return nil
}

// AddFile Add file to the database
func AddFile(info *FileInfo) error {
	comunication.Logger.Infof("INSERT INTO files (implant_id, file_path, file_name, file_type, file_size, created_at) VALUES ('%s', '%s', '%s', '%s', '%d', '%d')",
		info.ImplantID,
		info.FilePath,
		info.FileName,
		info.FileType,
		info.FileSize,
		info.CreatedAt)

	statement, err := dbConn.Prepare("INSERT INTO files (implant_id, file_path, file_name, file_type, file_size, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	_, err = statement.Exec(
		info.ImplantID,
		info.FilePath,
		info.FileName,
		info.FileType,
		info.FileSize,
		info.CreatedAt)

	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	return nil
}

// GetAllFiles Get all files
func GetAllFiles() ([]*FileInfo, error) {
	comunication.Logger.Info("SELECT * FROM files")
	statement, err := dbConn.Prepare("SELECT * FROM files")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	rows, err := statement.Query()
	if err != nil {
		comunication.Logger.Errorf("Error querying database: %v", err)
		return nil, err
	}

	var files []*FileInfo
	for rows.Next() {
		Info := new(FileInfo)
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
			comunication.Logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		files = append(files, Info)
	}
	if len(files) == 0 {
		return nil, comunication.ErrNoResults
	}

	return files, nil
}

// GetFilesByImplantID Get file by implant ID
func GetFilesByImplantID(ID string) ([]*FileInfo, error) {
	comunication.Logger.Infof("SELECT * FROM files WHERE implant_id = '%s'", ID)
	statement, err := dbConn.Prepare("SELECT * FROM files WHERE implant_id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	rows, err := statement.Query(ID)
	if err != nil {
		comunication.Logger.Errorf("Error querying database: %v", err)
		return nil, err
	}

	var files []*FileInfo
	for rows.Next() {
		Info := new(FileInfo)
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
			comunication.Logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}
		files = append(files, Info)
	}
	if len(files) == 0 {
		return nil, comunication.ErrNoResults
	}

	return files, nil
}

// GetFileByFID Get file by ID
func GetFileByFID(ID uint64) (*FileInfo, error) {
	comunication.Logger.Infof("SELECT * FROM files WHERE id = '%d'", ID)
	statement, err := dbConn.Prepare("SELECT * FROM files WHERE id = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	row := statement.QueryRow(ID)
	Info := new(FileInfo)
	err = row.Scan(
		&Info.ID,
		&Info.ImplantID,
		&Info.FilePath,
		&Info.FileName,
		&Info.FileType,
		&Info.FileSize,
		&Info.CreatedAt)
	if err != nil {
		comunication.Logger.Errorf("Error scanning row: %v", err)
		return nil, err
	}

	return Info, nil
}

// GetFileID Get file ID by implant ID and file name
func GetFileID(fileInfo *FileInfo) (int64, error) {
	var fileID int64
	statement, err := dbConn.Prepare("SELECT id FROM files WHERE implant_id = ? AND file_name = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return 0, err
	}
	defer statement.Close()

	err = statement.QueryRow(fileInfo.ImplantID, fileInfo.FileName).Scan(&fileID)
	if err != nil {
		comunication.Logger.Errorf("Error querying file ID: %v", err)
		return 0, err
	}

	return fileID, nil
}

// GetFileByImplantIDAndName Get file by implant ID and file name
func GetFileByImplantIDAndName(ID string, name string) (*FileInfo, error) {
	comunication.Logger.Infof("SELECT * FROM files WHERE implant_id = '%s' AND file_name = '%s'", ID, name)
	statement, err := dbConn.Prepare("SELECT * FROM files WHERE implant_id = ? AND file_name = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return nil, err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	row := statement.QueryRow(ID, name)
	var id int64 // Temporary variable to hold the id
	Info := new(FileInfo)
	err = row.Scan(
		&id,
		&Info.ImplantID,
		&Info.FilePath,
		&Info.FileName,
		&Info.FileType,
		&Info.FileSize,
		&Info.CreatedAt)
	if err != nil {
		comunication.Logger.Errorf("Error scanning row: %v", err)
		return nil, err
	}

	return Info, nil
}

// RemoveFile Remove file by implant ID and file
func RemoveFile(ID string, name string) error {
	comunication.Logger.Infof("DELETE FROM files WHERE implant_id = '%s' AND file_name = '%s'", ID, name)
	statement, err := dbConn.Prepare("DELETE FROM files WHERE implant_id = ? AND file_name = ?")
	if err != nil {
		comunication.Logger.Errorf("Error preparing statement: %v", err)
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			comunication.Logger.Errorf("Error closing statement: %v", err)
		}
	}(statement)

	resp, err := statement.Exec(ID, name)
	if err != nil {
		comunication.Logger.Errorf("Error executing query: %v", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		comunication.Logger.Errorf("Error getting rows affected: %v", err)
	}

	if rows == 0 {
		comunication.Logger.Warn("No rows affected")
		return errors.New(fmt.Sprintf("No file with name '%s' for implant '%s'", name, ID))
	}

	return nil
}
