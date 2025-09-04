# TaskReload - CRUD API для управления задачами

## 1. Введение

### Описание задания
Разработан полнофункциональный CRUD API для управления задачами с веб-интерфейсом. API позволяет создавать, читать, обновлять и удалять задачи, а также фильтровать их по статусу и приоритету.

### Выбранная СУБД
**SQLite** - легковесная встраиваемая СУБД, которая идеально подходит для:
- Простых приложений и прототипов
- Отсутствия необходимости в установке отдельного сервера БД
- Хранения данных в одном файле
- Быстрого развертывания

## 2. Структура проекта

```
TaskReload/
├── main.go                    # Главный файл приложения
├── go.mod                     # Зависимости Go модуля
├── go.sum                     # Хеши зависимостей
├── tasks.db                   # Файл базы данных SQLite
├── internal/
│   ├── models/
│   │   └── task.go           # Структуры данных и модели
│   ├── database/
│   │   └── database.go       # Функции для работы с БД
│   ├── handlers/
│   │   ├── task_handlers.go  # HTTP обработчики API
│   │   └── home_handler.go   # Обработчик главной страницы
│   └── middleware/
│       └── middleware.go     # CORS и логирование
├── templates/
│   └── index.html            # Веб-интерфейс
├── docs/                     # Swagger документация
└── README.md                 # Документация проекта
```

## 3. Реализация

### База данных

#### Подключение к БД
```go
// internal/database/database.go
func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite", "./tasks.db")
    if err != nil {
        return nil, err
    }
    
    if err = db.Ping(); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

#### Создание таблицы
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'pending',
    priority TEXT DEFAULT 'medium',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Go код для создания таблицы
```go
func CreateTables(db *sql.DB) error {
    query := `
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            status TEXT DEFAULT 'pending',
            priority TEXT DEFAULT 'medium',
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `
    
    _, err := db.Exec(query)
    return err
}
```

### API эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/tasks` | Получить список всех задач |
| GET | `/api/tasks/{id}` | Получить задачу по ID |
| POST | `/api/tasks` | Создать новую задачу |
| PUT | `/api/tasks/{id}` | Обновить существующую задачу |
| DELETE | `/api/tasks/{id}` | Удалить задачу |

#### Пример обработчика создания задачи
```go
// internal/handlers/task_handlers.go
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
```

#### Структуры данных
```go
// internal/models/task.go
type Task struct {
    ID          int       `json:"id" db:"id"`
    Title       string    `json:"title" db:"title" binding:"required"`
    Description string    `json:"description" db:"description"`
    Status      string    `json:"status" db:"status"`
    Priority    string    `json:"priority" db:"priority"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTaskRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
    Status      string `json:"status"`
    Priority    string `json:"priority"`
}
```

### Дополнительные возможности

#### Фильтрация задач
```go
// GET /api/tasks?status=completed&priority=high
type TaskFilter struct {
    Status   string `form:"status"`
    Priority string `form:"priority"`
}

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
    // ... обработка результата
}
```

#### Валидация данных
```go
type Task struct {
    Title string `json:"title" binding:"required"` // Обязательное поле
    // ... другие поля
}
```

#### Логирование и CORS
```go
// internal/middleware/middleware.go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

## 4. Инструкция по запуску

### Предварительные требования
- Go 1.19 или выше
- Доступ к интернету для загрузки зависимостей

### Установка и запуск

#### Вариант 1: Прямой запуск

1. **Клонирование проекта**
```bash
git clone <repository-url>
cd TaskReload
```

2. **Установка зависимостей**
```bash
go mod download
```

3. **Запуск приложения**
```bash
go run main.go
```

4. **Доступ к приложению**
- Веб-интерфейс: http://localhost:8080
- Swagger документация: http://localhost:8080/api/swagger/index.html

#### Вариант 2: Запуск через Docker

1. **Сборка и запуск с Docker Compose**
```bash
# Сборка и запуск
docker-compose up --build

# Запуск в фоновом режиме
docker-compose up -d --build

# Остановка
docker-compose down
```

2. **Запуск только с Docker**
```bash
# Сборка образа
docker build -t taskreload .

# Запуск контейнера
docker run -d -p 8080:8080 -v $(pwd)/data:/root/data --name taskreload-app taskreload

# Остановка контейнера
docker stop taskreload-app
docker rm taskreload-app
```

3. **Доступ к приложению**
- Веб-интерфейс: http://localhost:8080
- Swagger документация: http://localhost:8080/api/swagger/index.html

### Тестирование API

#### Создание задачи
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Изучить Go",
    "description": "Изучить основы языка Go",
    "status": "pending",
    "priority": "high"
  }'
```

#### Получение списка задач
```bash
curl http://localhost:8080/api/tasks
```

#### Фильтрация по статусу
```bash
curl "http://localhost:8080/api/tasks?status=completed"
```

## 5. Заключение

### Сложности и их решения

1. **Проблема с SQLite драйвером**
   - **Сложность**: Изначально использовался `github.com/mattn/go-sqlite3`, который требует CGO и компилятор C
   - **Решение**: Переход на `modernc.org/sqlite` - чистый Go драйвер без зависимостей от C

2. **Ошибка фронтенда с undefined данными**
   - **Сложность**: Поле `Data` в JSON ответе могло отсутствовать из-за тега `omitempty`
   - **Решение**: Убрал тег `omitempty` и добавил проверки на фронтенде

3. **Структура проекта**
   - **Сложность**: Организация кода в соответствии с Go best practices
   - **Решение**: Разделение на логические модули: models, database, handlers, middleware

### Возможные улучшения

1. **Функциональность**
   - Добавление аутентификации и авторизации
   - Поддержка вложенных задач (subtasks)
   - Система тегов для задач
   - Экспорт задач в различные форматы (CSV, PDF)

2. **Технические улучшения**
   - Добавление unit и integration тестов
   - Конфигурация через environment variables
   - Логирование в файл
   - Метрики и мониторинг

3. **Производительность**
   - Кэширование с Redis
   - Пагинация для больших списков задач
   - Оптимизация SQL запросов
   - Индексы в базе данных

4. **Фронтенд**
   - Переход на современный фреймворк (React, Vue.js)
   - PWA функциональность
   - Офлайн режим
   - Drag & drop для изменения статуса задач

### Итоги
Проект успешно демонстрирует:
- Понимание REST API принципов
- Работу с реляционными базами данных
- Структурирование Go приложений
- Создание простого, но функционального веб-интерфейса
- Обработку ошибок и валидацию данных

Код написан в соответствии с Go идиомами, имеет четкую структуру и готов для дальнейшего развития.

