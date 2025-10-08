package models

import "time"

type URL struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	ShortCode string    `gorm:"uniqueIndex;size:64" json:"short_code"`
	LongURL   string    `gorm:"not null;size:2000" json:"long_url"`
	Hits      uint64    `gorm:"default:0" json:"hits"`
	CreatedAt time.Time `json:"created_at"`
}
