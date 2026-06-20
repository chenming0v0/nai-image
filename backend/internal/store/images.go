package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Image 图片记录元信息（不含二进制 data）。
type Image struct {
	ID        string `json:"id"`
	Mime      string `json:"mime"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Source    string `json:"source"`
	CreatedAt int64  `json:"created_at"`
}

// SaveImage 存储图片二进制，返回元信息。
func SaveImage(db *sql.DB, data []byte, mime string, width, height int, source string) (*Image, error) {
	img := &Image{
		ID:        uuid.NewString(),
		Mime:      mime,
		Width:     width,
		Height:    height,
		Source:    source,
		CreatedAt: time.Now().Unix(),
	}
	_, err := db.Exec(
		"INSERT INTO images(id, data, mime, width, height, source, created_at) VALUES(?, ?, ?, ?, ?, ?, ?)",
		img.ID, data, mime, width, height, source, img.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// GetImageData 读取图片二进制和 mime。
func GetImageData(db *sql.DB, id string) (data []byte, mime string, err error) {
	err = db.QueryRow("SELECT data, mime FROM images WHERE id = ?", id).Scan(&data, &mime)
	return
}

// GetImageMeta 读取图片元信息。
func GetImageMeta(db *sql.DB, id string) (*Image, error) {
	img := &Image{}
	err := db.QueryRow(
		"SELECT id, mime, width, height, source, created_at FROM images WHERE id = ?", id,
	).Scan(&img.ID, &img.Mime, &img.Width, &img.Height, &img.Source, &img.CreatedAt)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// DeleteImage 删除图片。
func DeleteImage(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM images WHERE id = ?", id)
	return err
}
