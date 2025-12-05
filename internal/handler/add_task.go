package handler

import (
	"net/http"
	"strings"
	"time"

	"go_web_sample/internal/entity"
	"go_web_sample/internal/store"

	"github.com/labstack/echo/v4"
)

// AddTaskRequest タスク追加リクエスト
type AddTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

// AddTaskResponse タスク追加レスポンス
type AddTaskResponse struct {
	ID string `json:"id"`
}

// AddTask タスクを追加するハンドラー
func AddTask(c echo.Context) error {
	// 1. JSONリクエストボディをデコード
	var req AddTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// 2. バリデーション
	if strings.TrimSpace(req.Title) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
	}

	// Statusのバリデーション
	if req.Status != "" {
		validStatuses := []string{entity.StatusTodo, entity.StatusDoing, entity.StatusDone}
		isValid := false
		for _, status := range validStatuses {
			if req.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "status must be one of: todo, doing, done",
			})
		}
	} else {
		// Statusが指定されていない場合はデフォルトで "todo" を設定
		req.Status = entity.StatusTodo
	}

	// 3. タスクエンティティを生成
	taskID, err := store.GenerateTaskID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate task ID",
		})
	}

	task := &entity.Task{
		ID:      taskID,
		Title:   strings.TrimSpace(req.Title),
		Status:  req.Status,
		Created: time.Now(),
	}

	// 4. Redisに保存
	if err := store.SaveTask(task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to save task",
		})
	}

	// 5. 成功レスポンス
	return c.JSON(http.StatusOK, AddTaskResponse{
		ID: taskID,
	})
}
