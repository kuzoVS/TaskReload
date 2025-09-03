package main

import (
	"log"
	"taskreload/internal/database"
	"taskreload/internal/handlers"
	"taskreload/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "taskreload/docs"
)

// @title TaskReload API
// @version 1.0
// @description API для управления списком задач
// @host localhost:8080
// @BasePath /
func main() {
	// Инициализация базы данных
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	defer db.Close()

	// Создание таблиц
	if err := database.CreateTables(db); err != nil {
		log.Fatal("Ошибка создания таблиц:", err)
	}

	// Передача подключения к БД в обработчики
	handlers.SetDB(db)

	// Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Статические файлы для фронтенда
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Главная страница
	r.GET("/", handlers.HomePage)

	// API роуты
	api := r.Group("/api")
	{
		// Swagger документация
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// CRUD операции для задач
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTaskByID)
			tasks.POST("", handlers.CreateTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}
	}

	log.Println("Сервер запущен на http://localhost:8080")
	log.Println("Swagger документация: http://localhost:8080/api/swagger/index.html")
	
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
