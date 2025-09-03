package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"taskreload/internal/database"
	"taskreload/internal/models"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

// SetDB устанавливает подключение к базе данных для обработчиков
func SetDB(database *sql.DB) {
	db = database
}

// GetTasks получает список всех задач с возможностью фильтрации
// @Summary Получить список задач
// @Description Получить список всех задач с возможностью фильтрации по статусу и приоритету
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу"
// @Param priority query string false "Фильтр по приоритету"
// @Success 200 {object} models.TasksResponse
// @Router /api/tasks [get]
func GetTasks(c *gin.Context) {
	var filter models.TaskFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.TasksResponse{
			Success: false,
			Error:   "Неверные параметры фильтрации",
		})
		return
	}

	tasks, err := database.GetAllTasks(db, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TasksResponse{
			Success: false,
			Error:   "Ошибка получения задач из базы данных",
		})
		return
	}

	// Убеждаемся, что tasks не nil
	if tasks == nil {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, models.TasksResponse{
		Success: true,
		Data:    tasks,
		Total:   len(tasks),
		Message: "Задачи получены успешно",
	})
}

// GetTaskByID получает задачу по ID
// @Summary Получить задачу по ID
// @Description Получить задачу по указанному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} models.TaskResponse
// @Failure 404 {object} models.TaskResponse
// @Router /api/tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TaskResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		})
		return
	}

	task, err := database.GetTaskByID(db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.TaskResponse{
				Success: false,
				Error:   "Задача не найдена",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка получения задачи из базы данных",
		})
		return
	}

	c.JSON(http.StatusOK, models.TaskResponse{
		Success: true,
		Data:    task,
		Message: "Задача получена успешно",
	})
}

// CreateTask создает новую задачу
// @Summary Создать задачу
// @Description Создать новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Данные задачи"
// @Success 201 {object} models.TaskResponse
// @Failure 400 {object} models.TaskResponse
// @Router /api/tasks [post]
func CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.TaskResponse{
			Success: false,
			Error:   "Неверные данные задачи",
		})
		return
	}

	// Установка значений по умолчанию
	if req.Status == "" {
		req.Status = models.StatusPending
	}
	if req.Priority == "" {
		req.Priority = models.PriorityMedium
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
	}

	if err := database.InsertTask(db, task); err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка создания задачи в базе данных",
		})
		return
	}

	c.JSON(http.StatusCreated, models.TaskResponse{
		Success: true,
		Data:    task,
		Message: "Задача создана успешно",
	})
}

// UpdateTask обновляет существующую задачу
// @Summary Обновить задачу
// @Description Обновить существующую задачу по ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body models.UpdateTaskRequest true "Данные для обновления"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} models.TaskResponse
// @Failure 404 {object} models.TaskResponse
// @Router /api/tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TaskResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		})
		return
	}

	// Проверка существования задачи
	exists, err := database.TaskExists(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка проверки существования задачи",
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, models.TaskResponse{
			Success: false,
			Error:   "Задача не найдена",
		})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.TaskResponse{
			Success: false,
			Error:   "Неверные данные для обновления",
		})
		return
	}

	// Получение текущей задачи
	currentTask, err := database.GetTaskByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка получения текущей задачи",
		})
		return
	}

	// Обновление полей
	if req.Title != "" {
		currentTask.Title = req.Title
	}
	if req.Description != "" {
		currentTask.Description = req.Description
	}
	if req.Status != "" {
		currentTask.Status = req.Status
	}
	if req.Priority != "" {
		currentTask.Priority = req.Priority
	}

	if err := database.UpdateTask(db, id, currentTask); err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка обновления задачи в базе данных",
		})
		return
	}

	// Получение обновленной задачи
	updatedTask, err := database.GetTaskByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка получения обновленной задачи",
		})
		return
	}

	c.JSON(http.StatusOK, models.TaskResponse{
		Success: true,
		Data:    updatedTask,
		Message: "Задача обновлена успешно",
	})
}

// DeleteTask удаляет задачу по ID
// @Summary Удалить задачу
// @Description Удалить задачу по указанному идентификатору
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} models.TaskResponse
// @Failure 404 {object} models.TaskResponse
// @Router /api/tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.TaskResponse{
			Success: false,
			Error:   "Неверный ID задачи",
		})
		return
	}

	// Проверка существования задачи
	exists, err := database.TaskExists(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка проверки существования задачи",
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, models.TaskResponse{
			Success: false,
			Error:   "Задача не найдена",
		})
		return
	}

	if err := database.DeleteTask(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, models.TaskResponse{
			Success: false,
			Error:   "Ошибка удаления задачи из базы данных",
		})
		return
	}

	c.JSON(http.StatusOK, models.TaskResponse{
		Success: true,
		Message: "Задача удалена успешно",
	})
}
