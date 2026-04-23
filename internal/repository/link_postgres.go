package repository

import (
	"database/sql"
	"link-storage-service/internal/model"
	"time"
)

type PostgresLinkRepository struct {
	db *sql.DB
}

func NewPostgresLinkRepository(db *sql.DB) *PostgresLinkRepository {
	return &PostgresLinkRepository{db: db}
}

func (r *PostgresLinkRepository) Create(link *model.Link) error {
	query := `
    INSERT INTO links (short_code, original_url, created_at, visits)
    VALUES ($1, $2, $3, $4)
    RETURNING id`

	return r.db.QueryRow(
		query,
		link.ShortCode,
		link.OriginalURL,
		time.Now(),
		0,
	).Scan(&link.ID)
}

func (r *PostgresLinkRepository) GetByShortCode(code string) (*model.Link, error) {
	query := `
    SELECT id, short_code, original_url, created_at, visits
    FROM links
    WHERE short_code = $1`

	row := r.db.QueryRow(query, code)

	var link model.Link
	err := row.Scan(
		&link.ID,
		&link.ShortCode,
		&link.OriginalURL,
		&link.CreatedAt,
		&link.Visits,
	)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *PostgresLinkRepository) GetAll(limit, offset int) ([]model.Link, error) {
	query := `
    SELECT id, short_code, original_url, created_at, visits
    FROM links
    ORDER BY created_at DESC
    LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []model.Link

	for rows.Next() {
		var link model.Link
		err := rows.Scan(
			&link.ID,
			&link.ShortCode,
			&link.OriginalURL,
			&link.CreatedAt,
			&link.Visits,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (r *PostgresLinkRepository) Delete(code string) error {
	query := `DELETE FROM links WHERE short_code = $1`
	_, err := r.db.Exec(query, code)
	return err
}

func (r *PostgresLinkRepository) IncrementVisits(code string) (int64, error) {
	query := `
    UPDATE links
    SET visits = visits + 1
    WHERE short_code = $1
    RETURNING visits`

	var visits int64
	err := r.db.QueryRow(query, code).Scan(&visits)
	if err != nil {
		return 0, err
	}
	return visits, nil
}

func (r *PostgresLinkRepository) UpdateShortCode(id int64, code string) error {
	query := `
    UPDATE links
    SET short_code = $1
    WHERE id = $2`

	_, err := r.db.Exec(query, code, id)
	return err
}
