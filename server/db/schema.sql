CREATE TABLE IF NOT EXISTS implants (
    implant_id TEXT PRIMARY KEY,
    hostname TEXT,
    private_ip TEXT,
    public_ip TEXT,
    os TEXT,
    pid INTEGER,
    user TEXT,
    protection_name TEXT,
    last_check_in INTEGER,
    active BOOLEAN,
    kill_date INTEGER /*BANG!*/
);

CREATE TABLE IF NOT EXISTS listeners (
    listener_id TEXT PRIMARY KEY,
    config BLOB,
    host TEXT,
    port INTEGER,
    created_at INTEGER,
    kill_date INTEGER
);

/* Holds Implant Tasks information */
CREATE TABLE IF NOT EXISTS tasks (
    task_id TEXT PRIMARY KEY,
    implant_id TEXT,
    file_id INTEGER,
    task_type INTEGER,
    task_data BLOB,
    created_at INTEGER,
    completed BOOLEAN,
    completed_at INTEGER,
    task_result BLOB,
    FOREIGN KEY (implant_id) REFERENCES implants (implant_id),
    FOREIGN KEY (file_id) REFERENCES files (id)
);

/* Holds file information */
CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    implant_id TEXT,
    file_path TEXT,
    file_name TEXT,
    file_type TEXT,
    file_size INTEGER,
    created_at INTEGER,
    FOREIGN KEY (implant_id) REFERENCES implants (implant_id)
);
