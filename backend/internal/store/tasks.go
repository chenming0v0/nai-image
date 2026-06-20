package store

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Task 历史记录。
type Task struct {
	ID            string          `json:"id"`
	Model         string          `json:"model"`
	Params        json.RawMessage `json:"params"`
	Status        string          `json:"status"`
	Error         string          `json:"error,omitempty"`
	OutputImageID string          `json:"output_image_id,omitempty"`
	Seeds         json.RawMessage `json:"seeds,omitempty"`
	VibeCacheIDs  json.RawMessage `json:"vibe_cache_ids,omitempty"`
	Usage         json.RawMessage `json:"usage,omitempty"`
	UpstreamID    string          `json:"upstream_id,omitempty"`
	CreatedAt     int64           `json:"created_at"`
	FinishedAt    int64           `json:"finished_at,omitempty"`
}

// CreateTask 创建一条 running 状态的历史记录。
func CreateTask(db *sql.DB, model string, params json.RawMessage) (*Task, error) {
	t := &Task{
		ID:        uuid.NewString(),
		Model:     model,
		Params:    params,
		Status:    "running",
		CreatedAt: time.Now().Unix(),
	}
	_, err := db.Exec(
		"INSERT INTO tasks(id, model, params, status, created_at) VALUES(?, ?, ?, ?, ?)",
		t.ID, t.Model, string(t.Params), t.Status, t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// TaskResult 用于更新任务最终结果。
type TaskResult struct {
	Status        string // "done" / "error"
	Error         string
	OutputImageID string
	Seeds         json.RawMessage
	VibeCacheIDs  json.RawMessage
	Usage         json.RawMessage
	UpstreamID    string
}

// FinishTask 更新任务最终状态。
func FinishTask(db *sql.DB, id string, r *TaskResult) error {
	finishedAt := time.Now().Unix()
	seeds := ""
	if len(r.Seeds) > 0 {
		seeds = string(r.Seeds)
	}
	vibe := ""
	if len(r.VibeCacheIDs) > 0 {
		vibe = string(r.VibeCacheIDs)
	}
	usage := ""
	if len(r.Usage) > 0 {
		usage = string(r.Usage)
	}
	_, err := db.Exec(
		`UPDATE tasks SET status = ?, error = ?, output_image_id = ?, seeds = ?, vibe_cache_ids = ?, usage = ?, upstream_id = ?, finished_at = ? WHERE id = ?`,
		r.Status, r.Error, r.OutputImageID, seeds, vibe, usage, r.UpstreamID, finishedAt, id,
	)
	return err
}

// GetTask 读取单条任务。
func GetTask(db *sql.DB, id string) (*Task, error) {
	t := &Task{}
	var outputImageID, seeds, vibe, usage, upstreamID sql.NullString
	var errMsg sql.NullString
	var finishedAt sql.NullInt64
	err := db.QueryRow(
		`SELECT id, model, params, status, error, output_image_id, seeds, vibe_cache_ids, usage, upstream_id, created_at, finished_at FROM tasks WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Model, &t.Params, &t.Status, &errMsg, &outputImageID, &seeds, &vibe, &usage, &upstreamID, &t.CreatedAt, &finishedAt)
	if err != nil {
		return nil, err
	}
	t.Error = errMsg.String
	t.OutputImageID = outputImageID.String
	if seeds.Valid {
		t.Seeds = json.RawMessage(seeds.String)
	}
	if vibe.Valid {
		t.VibeCacheIDs = json.RawMessage(vibe.String)
	}
	if usage.Valid {
		t.Usage = json.RawMessage(usage.String)
	}
	t.UpstreamID = upstreamID.String
	if finishedAt.Valid {
		t.FinishedAt = finishedAt.Int64
	}
	return t, nil
}

// ListTasks 按创建时间倒序分页查询任务。
func ListTasks(db *sql.DB, limit, offset int) ([]*Task, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := db.Query(
		`SELECT id, model, params, status, error, output_image_id, seeds, vibe_cache_ids, usage, upstream_id, created_at, finished_at FROM tasks ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		var outputImageID, seeds, vibe, usage, upstreamID sql.NullString
		var errMsg sql.NullString
		var finishedAt sql.NullInt64
		if err := rows.Scan(&t.ID, &t.Model, &t.Params, &t.Status, &errMsg, &outputImageID, &seeds, &vibe, &usage, &upstreamID, &t.CreatedAt, &finishedAt); err != nil {
			return nil, 0, err
		}
		t.Error = errMsg.String
		t.OutputImageID = outputImageID.String
		if seeds.Valid {
			t.Seeds = json.RawMessage(seeds.String)
		}
		if vibe.Valid {
			t.VibeCacheIDs = json.RawMessage(vibe.String)
		}
		if usage.Valid {
			t.Usage = json.RawMessage(usage.String)
		}
		t.UpstreamID = upstreamID.String
		if finishedAt.Valid {
			t.FinishedAt = finishedAt.Int64
		}
		tasks = append(tasks, t)
	}
	return tasks, total, nil
}

// DeleteTask 删除单条任务。
func DeleteTask(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// DeleteAllTasks 清空所有任务。
func DeleteAllTasks(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM tasks")
	return err
}
