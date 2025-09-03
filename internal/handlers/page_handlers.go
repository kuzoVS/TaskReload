package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HomePage отображает главную страницу приложения
func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "TaskReload - Управление задачами",
	})
}
