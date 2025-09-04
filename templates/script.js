// Глобальные переменные для управления задачами
let tasks = [];
let editingTaskId = null;

// Загрузка задач при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    loadTasks();
});

// Загрузка задач с сервера
async function loadTasks() {
    try {
        const statusFilter = document.getElementById('statusFilter').value;
        const priorityFilter = document.getElementById('priorityFilter').value;
        
        let url = '/api/tasks';
        const params = new URLSearchParams();
        if (statusFilter) params.append('status', statusFilter);
        if (priorityFilter) params.append('priority', priorityFilter);
        
        if (params.toString()) {
            url += '?' + params.toString();
        }

        const response = await fetch(url);
        const data = await response.json();
        
        if (data.success && data.data) {
            tasks = Array.isArray(data.data) ? data.data : [];
            renderTasks();
        } else {
            const errorMsg = data.error || 'Неизвестная ошибка';
            showNotification('Ошибка загрузки задач: ' + errorMsg, 'error');
            tasks = [];
            renderTasks();
        }
    } catch (error) {
        showNotification('Ошибка загрузки задач: ' + error.message, 'error');
    }
}

// Фильтрация задач
function filterTasks() {
    loadTasks();
}

// Обновление списка задач
function refreshTasks() {
    loadTasks();
}

// Отображение задач
function renderTasks() {
    const container = document.getElementById('tasksContainer');
    
    // Убеждаемся, что tasks - это массив
    if (!Array.isArray(tasks)) {
        tasks = [];
    }
    
    if (tasks.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-clipboard-list"></i>
                <h3>Нет задач</h3>
                <p>Создайте первую задачу, чтобы начать работу</p>
            </div>
        `;
        return;
    }

    container.innerHTML = tasks.map(task => `
        <div class="task-card">
            <div class="task-header">
                <div>
                    <div class="task-title">${escapeHtml(task.title)}</div>
                    <div class="task-description">${escapeHtml(task.description || '')}</div>
                </div>
            </div>
            <div class="task-meta">
                <span class="badge badge-status">${getStatusText(task.status)}</span>
                <span class="badge badge-priority ${task.priority}">${getPriorityText(task.priority)}</span>
                <small style="color: #7f8c8d; margin-left: auto;">
                    ${new Date(task.created_at).toLocaleDateString('ru-RU')}
                </small>
            </div>
            <div class="task-actions">
                <button class="btn btn-primary btn-sm" onclick="editTask(${task.id})">
                    <i class="fas fa-edit"></i> Изменить
                </button>
                <button class="btn btn-danger btn-sm" onclick="deleteTask(${task.id})">
                    <i class="fas fa-trash"></i> Удалить
                </button>
            </div>
        </div>
    `).join('');
}

// Открытие модального окна для создания задачи
function openCreateModal() {
    editingTaskId = null;
    document.getElementById('modalTitle').textContent = 'Новая задача';
    document.getElementById('taskForm').reset();
    document.getElementById('taskStatus').value = 'pending';
    document.getElementById('taskPriority').value = 'medium';
    document.getElementById('taskModal').style.display = 'block';
}

// Открытие модального окна для редактирования задачи
function editTask(taskId) {
    if (!Array.isArray(tasks)) {
        showNotification('Ошибка: список задач недоступен', 'error');
        return;
    }
    
    const task = tasks.find(t => t.id === taskId);
    if (!task) return;

    editingTaskId = taskId;
    document.getElementById('modalTitle').textContent = 'Редактировать задачу';
    document.getElementById('taskTitle').value = task.title;
    document.getElementById('taskDescription').value = task.description || '';
    document.getElementById('taskStatus').value = task.status;
    document.getElementById('taskPriority').value = task.priority;
    document.getElementById('taskModal').style.display = 'block';
}

// Закрытие модального окна
function closeModal() {
    document.getElementById('taskModal').style.display = 'none';
    editingTaskId = null;
}

// Обработка отправки формы
document.getElementById('taskForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const taskData = {
        title: document.getElementById('taskTitle').value.trim(),
        description: document.getElementById('taskDescription').value.trim(),
        status: document.getElementById('taskStatus').value,
        priority: document.getElementById('taskPriority').value
    };

    if (!taskData.title) {
        showNotification('Название задачи обязательно', 'error');
        return;
    }

    try {
        const url = editingTaskId ? `/api/tasks/${editingTaskId}` : '/api/tasks';
        const method = editingTaskId ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(taskData)
        });

        const data = await response.json();
        
        if (data.success) {
            showNotification(
                editingTaskId ? 'Задача обновлена успешно' : 'Задача создана успешно', 
                'success'
            );
            closeModal();
            loadTasks();
        } else {
            showNotification('Ошибка: ' + data.error, 'error');
        }
    } catch (error) {
        showNotification('Ошибка: ' + error.message, 'error');
    }
});

// Удаление задачи
async function deleteTask(taskId) {
    if (!Array.isArray(tasks)) {
        showNotification('Ошибка: список задач недоступен', 'error');
        return;
    }
    
    if (!confirm('Вы уверены, что хотите удалить эту задачу?')) {
        return;
    }

    try {
        const response = await fetch(`/api/tasks/${taskId}`, {
            method: 'DELETE'
        });

        const data = await response.json();
        
        if (data.success) {
            showNotification('Задача удалена успешно', 'success');
            loadTasks();
        } else {
            showNotification('Ошибка удаления: ' + data.error, 'error');
        }
    } catch (error) {
        showNotification('Ошибка удаления: ' + error.message, 'error');
    }
}

// Вспомогательные функции
function getStatusText(status) {
    const statusMap = {
        'pending': 'Ожидает',
        'in_progress': 'В работе',
        'completed': 'Завершено',
        'cancelled': 'Отменено'
    };
    return statusMap[status] || status;
}

function getPriorityText(priority) {
    const priorityMap = {
        'low': 'Низкий',
        'medium': 'Средний',
        'high': 'Высокий'
    };
    return priorityMap[priority] || priority;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showNotification(message, type) {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.classList.add('show');
    }, 100);
    
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 3000);
}

// Закрытие модального окна при клике вне его
window.onclick = function(event) {
    const modal = document.getElementById('taskModal');
    if (event.target === modal) {
        closeModal();
    }
}
