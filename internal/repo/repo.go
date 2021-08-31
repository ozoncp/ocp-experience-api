package repo

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	sql "github.com/jmoiron/sqlx"

	"github.com/ozoncp/ocp-experience-api/internal/models"
)

var NotFound = errors.New("experience does not exist")

// IRepo is an experience storage interface
type IRepo interface {
	Add(ctx context.Context, request models.Experience) (uint64, error)
	AddExperiences(ctx context.Context, request []models.Experience) ([]uint64, error)
	List(ctx context.Context, limit, offset uint64) ([]models.Experience, error)
	Describe(ctx context.Context, id uint64) (models.Experience, error)
	Remove(ctx context.Context, id uint64) (bool, error)
	Update(ctx context.Context, experience models.Experience) error
}

// NewRepo creates a new Repo
func NewRepo(db *sql.DB) *Repo {
	cache := sq.NewStmtCache(db)

	return &Repo{
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(cache),
	}
}

// Repo is IRepo impl
type Repo struct {
	builder sq.StatementBuilderType
}

// Add adds to db experience and returns its id
func (r *Repo) Add(ctx context.Context, experience models.Experience) (uint64, error) {
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
func (r *Repo) AddExperiences(ctx context.Context, experiences []models.Experience) ([]uint64, error) {
	query := r.builder.Insert("experiences").Columns("user_id", "type", "from", "to", "level").Suffix("RETURNING id")

	for _, experience := range experiences {
		query = query.Values(experience.UserId, experience.Type, experience.From, experience.To, experience.Level)
	}

	rows, err := query.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	newIds := make([]uint64, 0, len(experiences))

	for rows.Next() {
		var id uint64 = 0
		rows.Scan(&id)

		newIds = append(newIds, id)
	}

	return newIds, nil
}

// List returns an experience list
func (r *Repo) List(ctx context.Context, limit, offset uint64) ([]models.Experience, error) {
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
func (r *Repo) Describe(ctx context.Context, id uint64) (models.Experience, error) {
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
func (r *Repo) Remove(ctx context.Context, id uint64) (bool, error) {
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

// Update updates existing experience, returns NotFound error if request does not exist
func (r *Repo) Update(ctx context.Context, experience models.Experience) error {
	query := r.builder.Update("experiences")

	if experience.UserId != 0 {
		query = query.Set("user_id", experience.UserId)
	}

	if experience.Type != 0 {
		query = query.Set("type", experience.Type)
	}

	if !experience.From.IsZero() {
		query = query.Set("from", experience.From)
	}

	if !experience.To.IsZero() {
		query = query.Set("to", experience.To)
	}

	if experience.Level != 0 {
		query = query.Set("level", experience.Level)
	}

	query = query.Where("id = ?", experience.Id)
	ret, err := query.ExecContext(ctx)

	if err != nil {
		return err
	}

	rowsUpdated, err := ret.RowsAffected()

	if err != nil {
		return err
	}

	if rowsUpdated == 0 {
		return NotFound
	}

	return nil
}
