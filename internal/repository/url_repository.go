package repository

import (
	"URLShortener/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

// URLRepository handles database operations for URLs
type URLRepository struct {
	db *sql.DB
}

// NewURLRepository creates a new URL repository
func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

// GetByShortCode retrieves a URL by its short code
func (r *URLRepository) GetByShortCode(shortCode string) (*models.URL, error) {
	query := `
		SELECT id, short_code, long_url, created_at, clicks 
		FROM urls 
		WHERE short_code = $1
	`

	url := &models.URL{}
	err := r.db.QueryRow(query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.LongURL,
		&url.CreatedAt,
		&url.Clicks,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("URL not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	return url, nil
}

// GetByLongURL checks if a long URL already exists
func (r *URLRepository) GetByLongURL(longURL string) (*models.URL, error) {
	query := `
		SELECT id, short_code, long_url, created_at, clicks 
		FROM urls 
		WHERE long_url = $1
	`

	url := &models.URL{}
	err := r.db.QueryRow(query, longURL).Scan(
		&url.ID,
		&url.ShortCode,
		&url.LongURL,
		&url.CreatedAt,
		&url.Clicks,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not an error, just doesn't exist
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check long URL: %w", err)
	}

	return url, nil
}

// Create inserts a new shortened URL
func (r *URLRepository) Create(shortCode string, longURL string) (*models.URL, error) {
	query := `
		INSERT INTO urls (short_code, long_url) 
		VALUES ($1, $2) 
		RETURNING id, short_code, long_url, created_at, clicks
	`

	url := &models.URL{}
	err := r.db.QueryRow(query, shortCode, longURL).Scan(
		&url.ID,
		&url.ShortCode,
		&url.LongURL,
		&url.CreatedAt,
		&url.Clicks,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	return url, nil
}

// IncrementClicks increments the click count for a URL
func (r *URLRepository) IncrementClicks(shortCode string) error {
	query := `
		UPDATE urls 
		SET clicks = clicks + 1 
		WHERE short_code = $1
	`

	result, err := r.db.Exec(query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("URL not found")
	}

	return nil
}

// Delete removes a URL by short code (useful for testing)
func (r *URLRepository) Delete(shortCode string) error {
	query := `DELETE FROM urls WHERE short_code = $1`

	result, err := r.db.Exec(query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("URL not found")
	}

	return nil
}

// GetAll retrieves all URLs (for testing/admin)
func (r *URLRepository) GetAll() ([]*models.URL, error) {
	query := `
		SELECT id, short_code, long_url, created_at, clicks 
		FROM urls 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all URLs: %w", err)
	}
	defer rows.Close()

	var urls []*models.URL
	for rows.Next() {
		url := &models.URL{}
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.LongURL,
			&url.CreatedAt,
			&url.Clicks,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return urls, nil
}
