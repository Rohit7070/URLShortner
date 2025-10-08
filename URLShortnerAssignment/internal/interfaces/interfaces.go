package interfaces

import "URLShortnerAssignment/internal/models"

type URLRepository interface {
	Create(url *models.URL) error
	GetByShortCode(code string) (*models.URL, error)
	GetByLongURL(long string) (*models.URL, error)
	IncrementHits(code string) error
}

type URLService interface {
	Shorten(longURL string, custom string) (*models.URL, error)
	Resolve(code string) (*models.URL, error)
	Stats(code string) (*models.URL, error)
}
