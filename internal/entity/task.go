package entity

import "time"

// タスクステータスの定数定義
const (
	StatusTodo  = "todo"
	StatusDoing = "doing"
	StatusDone  = "done"
)

// Task タスクエンティティ
type Task struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}
