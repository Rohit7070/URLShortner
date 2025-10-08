package service

import (
	"URLShortnerAssignment/internal/interfaces"
	"URLShortnerAssignment/internal/models"
	"URLShortnerAssignment/internal/repository"
	"crypto/rand"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"strings"
	"time"
)

var (
	ErrAlreadyExists = errors.New("short code already exists")
	ErrNotFound      = errors.New("not found")
)

// Service implements URL shortening business logic.
type Service struct {
	repo       interfaces.URLRepository
	baseDomain string
	codeLen    int
}

// NewURLService initializes a new Service instance.
func NewURLService(repo interfaces.URLRepository, baseDomain string) *Service {
	return &Service{repo: repo, baseDomain: baseDomain, codeLen: 6}
}

// Shorten creates a new short URL, or returns an existing one if already present.
func (s *Service) Shorten(longURL string, custom string) (*models.URL, error) {
	if longURL == "" {
		return nil, errors.New("long URL required")
	}

	// --- Handle custom alias case ---
	if custom != "" {
		// Check if alias already exists
		if _, err := s.repo.GetByShortCode(custom); err == nil {
			return nil, ErrAlreadyExists
		}

		u := &models.URL{
			ShortCode: custom,
			LongURL:   longURL,
			CreatedAt: time.Now(),
		}
		err := s.repo.Create(u)
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrAlreadyExists
		}
		return u, err
	}

	// --- Return existing if same longURL exists ---
	if existing, err := s.repo.GetByLongURL(longURL); err == nil {
		return existing, nil
	}

	// --- Auto-generate new code ---
	for i := 0; i < 5; i++ {
		code, _ := generateShortCode(s.codeLen)
		u := &models.URL{
			ShortCode: code,
			LongURL:   longURL,
			CreatedAt: time.Now(),
		}

		err := s.repo.Create(u)
		if err == nil {
			return u, nil
		}
		// If collision (duplicate code), retry
		if errors.Is(err, repository.ErrDuplicate) {
			continue
		}
		// Other DB error
		return nil, err
	}
	return nil, fmt.Errorf("failed to create short code after retries")
}

// Resolve finds the original long URL by its short code.
func (s *Service) Resolve(code string) (*models.URL, error) {
	u, err := s.repo.GetByShortCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, repository.ErrDuplicate) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// Increment hit count (best-effort)
	_ = s.repo.IncrementHits(code)
	return u, nil
}

// Stats returns the statistics for a given short code.
func (s *Service) Stats(code string) (*models.URL, error) {
	u, err := s.repo.GetByShortCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

// BaseDomain exposes the configured base domain.
func (s *Service) BaseDomain() string {
	return s.baseDomain
}

// --- Short code generator ---

const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generateShortCode creates a random alphanumeric string of given length.
func generateShortCode(n int) (string, error) {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(base62[num.Int64()])
	}
	return sb.String(), nil
}
