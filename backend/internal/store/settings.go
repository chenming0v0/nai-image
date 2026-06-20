package store

import (
	"database/sql"
	"strconv"
)

// Settings 后端配置（与 models.Settings 一致，但放 store 包避免循环依赖）。
type Settings struct {
	UpstreamBaseURL string `json:"upstream_base_url"`
	UpstreamAPIKey  string `json:"upstream_api_key"`
	DefaultModel    string `json:"default_model"`
	RequestTimeout  int    `json:"request_timeout_seconds"`
}

func settingGet(db *sql.DB, key string) (string, bool) {
	var v string
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&v)
	if err != nil {
		return "", false
	}
	return v, true
}

func settingSet(db *sql.DB, key, value string) error {
	_, err := db.Exec("INSERT INTO settings(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value", key, value)
	return err
}

// GetSettings 读取配置，defaults 作为未设置项的回退值。
func GetSettings(db *sql.DB, defaults Settings) (Settings, error) {
	s := defaults
	if v, ok := settingGet(db, "upstream_base_url"); ok {
		s.UpstreamBaseURL = v
	}
	if v, ok := settingGet(db, "upstream_api_key"); ok {
		s.UpstreamAPIKey = v
	}
	if v, ok := settingGet(db, "default_model"); ok {
		s.DefaultModel = v
	}
	if v, ok := settingGet(db, "request_timeout_seconds"); ok {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			s.RequestTimeout = n
		}
	}
	return s, nil
}

// UpdateSettings 全量更新配置（空字符串也会写入）。
func UpdateSettings(db *sql.DB, s Settings) error {
	if err := settingSet(db, "upstream_base_url", s.UpstreamBaseURL); err != nil {
		return err
	}
	if err := settingSet(db, "upstream_api_key", s.UpstreamAPIKey); err != nil {
		return err
	}
	if err := settingSet(db, "default_model", s.DefaultModel); err != nil {
		return err
	}
	if err := settingSet(db, "request_timeout_seconds", strconv.Itoa(s.RequestTimeout)); err != nil {
		return err
	}
	return nil
}
