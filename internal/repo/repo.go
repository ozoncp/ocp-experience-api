package repo

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	sql "github.com/jmoiron/sqlx"

	"github.com/ozoncp/ocp-experience-api/internal/models"
)

var NotFound = errors.New("experience does not exist")

// Repo is an experience storage interface
type Repo interface {
	Add(ctx context.Context, request models.Experience) (uint64, error)
	AddExperiences(ctx context.Context, request []models.Experience) error
	List(ctx context.Context, limit, offset uint64) ([]models.Experience, error)
	Describe(ctx context.Context, id uint64) (models.Experience, error)
	Remove(ctx context.Context, id uint64) (bool, error)
}

// NewRepo creates a new Repo
func NewRepo(db *sql.DB) Repo {
	cache := sq.NewStmtCache(db)

	return &repo{
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(cache),
	}
}

// Repo impl
type repo struct {
	builder sq.StatementBuilderType
}

// Add adds to db experience and returns its id
func (r *repo) Add(ctx context.Context, experience models.Experience) (uint64, error) {
	query := r.builder.Insert("experiences").
		Columns("user_id", "type", "from", "to", "level").
		Suffix("RETURNING id").
		Values(experience.UserId, experience.Type, experience.From, experience.To, experience.Level)

	var id uint64 = 0
	rows, err := query.QueryContext(ctx)

	if err != nil {
		return id, err
	}

	rows.Next()
	scanErr := rows.Scan(&id)

	if scanErr != nil {
		return id, scanErr
	}

	return id, nil
}

// AddExperiences adds to db experience slice
func (r *repo) AddExperiences(ctx context.Context, experiences []models.Experience) error {
	query := r.builder.Insert("experiences").Columns("user_id", "type", "from", "to", "level")

	for _, experience := range experiences {
		query = query.Values(experience.UserId, experience.Type, experience.From, experience.To, experience.Level)
	}

	_, err := query.ExecContext(ctx)

	if err != nil {
		return err
	}

	return nil
}

// List returns an experience list
func (r *repo) List(ctx context.Context, limit, offset uint64) ([]models.Experience, error) {
	query := r.builder.Select("id, user_id, type, from, to, level").
		From("experiences").
		Offset(offset).
		Limit(limit)

	rows, err := query.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	experiences := make([]models.Experience, 0, limit)

	for rows.Next() {
		var experience models.Experience
		scanErr := rows.Scan(&experience.Id, &experience.UserId, &experience.Type, &experience.From, &experience.To, &experience.Level)

		if scanErr != nil {
			return nil, scanErr
		}

		experiences = append(experiences, experience)
	}

	return experiences, nil
}

// Describe returns experience by id
func (r *repo) Describe(ctx context.Context, id uint64) (models.Experience, error) {
	query := r.builder.Select("id, user_id, type, from, to, level").
		From("experiences").
		Where("id = ?", id)

	row, err := query.QueryContext(ctx)

	if err != nil {
		return models.Experience{}, err
	}

	var experience models.Experience

	if !row.Next() {
		return models.Experience{}, NotFound
	}

	scanErr := row.Scan(&experience.Id, &experience.UserId, &experience.Type, &experience.From, &experience.To, &experience.Level)

	if scanErr != nil {
		return models.Experience{}, scanErr
	}

	return experience, nil
}

// Remove deletes experience by id
func (r *repo) Remove(ctx context.Context, id uint64) (bool, error) {
	query := r.builder.Delete("experiences").Where("id = ?", id)
	ret, err := query.ExecContext(ctx)

	if err != nil {
		return false, err
	}

	rowsDeleted, err := ret.RowsAffected()

	if err != nil {
		return false, err
	}

	return rowsDeleted > 0, err
}
