package database

import (
	"database/sql"
	"log"
	"taskreload/internal/models"

	_ "modernc.org/sqlite"
)

// InitDB инициализирует подключение к базе данных
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Подключение к базе данных установлено")
	return db, nil
}

// CreateTables создает необходимые таблицы
func CreateTables(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT DEFAULT 'pending',
		priority TEXT DEFAULT 'medium',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Таблицы созданы успешно")
	return nil
}

// InsertTask добавляет новую задачу в базу данных
func InsertTask(db *sql.DB, task *models.Task) error {
	query := `
	INSERT INTO tasks (title, description, status, priority, created_at, updated_at)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := db.Exec(query, task.Title, task.Description, task.Status, task.Priority)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = int(id)
	return nil
}

// GetTaskByID получает задачу по ID
func GetTaskByID(db *sql.DB, id int) (*models.Task, error) {
	query := `SELECT id, title, description, status, priority, created_at, updated_at FROM tasks WHERE id = ?`
	
	var task models.Task
	err := db.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.CreatedAt, &task.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &task, nil
}

// GetAllTasks получает все задачи с возможностью фильтрации
func GetAllTasks(db *sql.DB, filter *models.TaskFilter) ([]models.Task, error) {
	query := `SELECT id, title, description, status, priority, created_at, updated_at FROM tasks WHERE 1=1`
	var args []interface{}

	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}

	if filter.Priority != "" {
		query += ` AND priority = ?`
		args = append(args, filter.Priority)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Всегда возвращаем слайс, даже если он пустой
	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(db *sql.DB, id int, task *models.Task) error {
	query := `
	UPDATE tasks 
	SET title = ?, description = ?, status = ?, priority = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	_, err := db.Exec(query, task.Title, task.Description, task.Status, task.Priority, id)
	return err
}

// DeleteTask удаляет задачу по ID
func DeleteTask(db *sql.DB, id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// TaskExists проверяет существование задачи
func TaskExists(db *sql.DB, id int) (bool, error) {
	var exists int
	query := `SELECT 1 FROM tasks WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
