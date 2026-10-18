package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"api/models"
)

var ErrJournalNotFound = errors.New("journal not found")

type JournalRepository struct {
	db *pgxpool.Pool
}

func NewJournalRepository(db *pgxpool.Pool) *JournalRepository {
	return &JournalRepository{
		db: db,
	}
}

func (r *JournalRepository) CreateJournal(ctx context.Context, journal models.Journal) (models.Journal, error) {
	if journal.Uuid == uuid.Nil() {
		journal.Uuid = uuid.New()
	}

	images, err := json.Marshal(journal.Images)
	if err != nil {
		return models.Journal{}, fmt.Errorf("marshal images: %w", err)
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO journals (
			uuid,
			title,
			category,
			organization,
			date,
			venue,
			role,
			points,
			article,
			images
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		journal.Uuid,
		journal.Title,
		journal.Category,
		journal.Organization,
		journal.Date,
		journal.Venue,
		journal.Role,
		journal.Points,
		journal.Article,
		images,
	)

	if err != nil {
		return models.Journal{}, fmt.Errorf("create journal: %w", err)
	}
	return journal, nil
}

func (r *JournalRepository) FindJournalByID(ctx context.Context, category models.Category, Uuid uuid.UUID) (models.Journal, error) {
	var journal models.Journal
	var images []byte

	err := r.db.QueryRow(ctx,
		`SELECT
			uuid,
			title,
			category,
			organization,
			date,
			venue,
			role,
			points,
			article,
			images
		FROM journals
		WHERE uuid = $1 AND category = $2`,
		Uuid,
		category,
	).Scan(
		&journal.Uuid,
		&journal.Title,
		&journal.Category,
		&journal.Organization,
		&journal.Date,
		&journal.Venue,
		&journal.Role,
		&journal.Points,
		&journal.Article,
		&images,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Journal{}, ErrJournalNotFound
	}
	if err != nil {
		return models.Journal{}, fmt.Errorf("find journal: %w", err)
	}
	if err := json.Unmarshal(images, &journal.Images); err != nil {
		return models.Journal{}, fmt.Errorf("unmarshal images: %w", err)
	}
	return journal, nil
}

func (r *JournalRepository) FindJournalsByCategory(ctx context.Context, category models.Category) ([]models.JournalPreview, error) {
	rows, err := r.db.Query(ctx,
		`SELECT
			uuid,
			title,
			organization
		FROM journals
		WHERE category = $1`,
		category,
	)

	if err != nil {
		return nil, fmt.Errorf("find journals by category: %w", err)
	}

	defer rows.Close()

	var journals []models.JournalPreview

	for rows.Next() {
		var journal models.JournalPreview

		if err := rows.Scan(
			&journal.Uuid,
			&journal.Title,
			&journal.Organization,
		); err != nil {
			return nil, fmt.Errorf("scan journal: %w", err)
		}
		journals = append(journals, journal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate journals: %w", err)
	}
	return journals, nil
}

func (r *JournalRepository) UpdateJournal(ctx context.Context, Uuid uuid.UUID, update models.JournalUpdate) (models.Journal, error) {
	var imagesData []byte
	var journal models.Journal

	var images any
	if update.Images != nil {
		data, err := json.Marshal(*update.Images)
		if err != nil {
			return models.Journal{}, fmt.Errorf("marshal images: %w", err)
		}
		images = data
	}

	row := r.db.QueryRow(ctx,
		`UPDATE journals
		SET
			title = COALESCE($1, title),
			category = COALESCE($2, category),
			organization = COALESCE($3, organization),
			date = COALESCE($4, date),
			venue = COALESCE($5, venue),
			role = COALESCE($6, role),
			points = COALESCE($7, points),
			article = COALESCE($8, article),
			images = COALESCE($9, images)
		WHERE uuid = $10
		RETURNING
			uuid,
			title,
			category,
			organization,
			date,
			venue,
			role,
			points,
			article,
			images
		`,
		update.Title,
		update.Category,
		update.Organization,
		update.Date,
		update.Venue,
		update.Role,
		update.Points,
		update.Article,
		images,
		Uuid,
	)

	err := row.Scan(
		&journal.Uuid,
		&journal.Title,
		&journal.Category,
		&journal.Organization,
		&journal.Date,
		&journal.Venue,
		&journal.Role,
		&journal.Points,
		&journal.Article,
		&imagesData,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Journal{}, ErrJournalNotFound
	}
	if err != nil {
		return models.Journal{}, fmt.Errorf("update journal: %w", err)
	}
	if err := json.Unmarshal(imagesData, &journal.Images); err != nil {
		return models.Journal{}, fmt.Errorf("unmarshal images: %w", err)
	}
	return journal, nil
}

func (r *JournalRepository) DeleteJournal(ctx context.Context, Uuid uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM journals
		WHERE uuid = $1`,
		Uuid,
	)
	if err != nil {
		return fmt.Errorf("delete journal: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrJournalNotFound
	}
	return nil
}
