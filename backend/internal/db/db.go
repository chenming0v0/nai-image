package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB 是全局数据库连接，由 Init 初始化。
var DB *sql.DB

// Init 打开 SQLite 数据库并执行迁移。
func Init(path string) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create db dir: %w", err)
		}
	}
	// modernc.org/sqlite 使用 ?_pragma=xxx 进行 pragma 配置
	dsn := path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	d, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	// SQLite 单写并发限制
	d.SetMaxOpenConns(1)
	if err := d.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	if err := migrate(d); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	DB = d
	return nil
}

func migrate(d *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS settings (
	key   TEXT PRIMARY KEY,
	value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
	id         TEXT PRIMARY KEY,
	data       BLOB    NOT NULL,
	mime       TEXT    NOT NULL,
	width      INTEGER,
	height     INTEGER,
	source     TEXT    NOT NULL DEFAULT 'generated',
	created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at DESC);

CREATE TABLE IF NOT EXISTS tasks (
	id              TEXT PRIMARY KEY,
	model           TEXT    NOT NULL,
	params          TEXT    NOT NULL,
	status          TEXT    NOT NULL,
	error           TEXT,
	output_image_id TEXT,
	seeds           TEXT,
	vibe_cache_ids  TEXT,
	usage           TEXT,
	upstream_id     TEXT,
	created_at      INTEGER NOT NULL,
	finished_at     INTEGER
);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
`
	_, err := d.Exec(schema)
	return err
}

// Close 关闭数据库连接。
func Close() error {
	if DB == nil {
		return nil
	}
	return DB.Close()
}
