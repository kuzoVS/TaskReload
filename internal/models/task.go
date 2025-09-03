package models

import (
	"time"
)

// Task представляет задачу
type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	Priority    string    `json:"priority" db:"priority"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTaskRequest представляет запрос на создание задачи
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

// UpdateTaskRequest представляет запрос на обновление задачи
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

// TaskFilter представляет фильтр для задач
type TaskFilter struct {
	Status   string `form:"status"`
	Priority string `form:"priority"`
}

// TaskResponse представляет ответ с задачей
type TaskResponse struct {
	Success bool        `json:"success"`
	Data    *Task      `json:"data"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TasksResponse представляет ответ со списком задач
type TasksResponse struct {
	Success bool        `json:"success"`
	Data    []Task     `json:"data"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Total   int         `json:"total"`
}

// Константы для статусов и приоритетов
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusCancelled  = "cancelled"

	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)
