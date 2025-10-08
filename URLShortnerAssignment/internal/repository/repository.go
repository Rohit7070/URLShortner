package repository

import (
	"URLShortnerAssignment/internal/models"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrDuplicate = errors.New("duplicate short code")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(url *models.URL) error {
	err := r.db.Create(url).Error
	if err != nil {
		// Detect unique constraint violation (works for SQLite, Postgres, MySQL)
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "unique") || strings.Contains(lower, "duplicate") {
			return ErrDuplicate
		}
	}
	return err
}

func (r *Repository) GetByShortCode(code string) (*models.URL, error) {
	var u models.URL
	res := r.db.Where("short_code = ?", code).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &u, res.Error
}

func (r *Repository) GetByLongURL(long string) (*models.URL, error) {
	var u models.URL
	res := r.db.Where("long_url = ?", long).First(&u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, sql.ErrNoRows
	}
	return &u, res.Error
}

func (r *Repository) IncrementHits(code string) error {
	res := r.db.Model(&models.URL{}).
		Where("short_code = ?", code).
		UpdateColumn("hits", gorm.Expr("hits + ?", 1))
	if res.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return res.Error
}
