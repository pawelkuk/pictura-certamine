package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, c *model.Contest) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO contest(
			id,
			name,
			slug,
			start_time,
			end_time,
			is_active
		) VALUES(?, ?, ?, ?, ?, ?)
		RETURNING id`,
		c.ID,
		c.Name,
		c.Slug.Value,
		c.Start.Format(time.RFC3339),
		c.End.Format(time.RFC3339),
		c.IsActive,
	)
	if err != nil {
		return fmt.Errorf("could not create contest: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, c *model.Contest) error {
	row := r.DB.QueryRowContext(ctx,
		`SELECT 
			name,
			slug,
			start_time,
			end_time,
			is_active
		FROM 
			contest
		WHERE id = ?`,
		c.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query row with id=%s: %w", c.ID, row.Err())
	}
	var name, slugRaw, startStr, endStr string
	var isActive bool
	err := row.Scan(&name, &slugRaw, &startStr, &endStr, &isActive)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}
	slug := model.ParseSlug(slugRaw)
	c.Name = name
	c.Slug = slug
	c.Start = start
	c.End = end
	c.IsActive = isActive
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, c *model.Contest) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE 
			contest 
		SET 
			name = ?,
			slug = ?,
			start_time = ?,
			end_time = ?,
			is_active = ?
		WHERE
			id = ?`,
		c.Name,
		c.Slug.Value,
		c.Start.Format(time.RFC3339),
		c.End.Format(time.RFC3339),
		c.IsActive,
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update contest: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, c *model.Contest) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM contest WHERE id = ?`, c.ID)
	if err != nil {
		return fmt.Errorf("could not delete contest: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.ContestQueryFilter) ([]model.Contest, error) {
	q := `
	select
		id,
		name,
		slug,
		start_time,
		end_time,
		is_active
	from
		contest
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	contests := []model.Contest{}
	for rows.Next() {
		c := &model.Contest{}
		var id, name, slugRaw, startStr, endStr string
		var isActive bool
		err := rows.Scan(&id, &name, &slugRaw, &startStr, &endStr, &isActive)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}
		slug := model.ParseSlug(slugRaw)
		c.ID = id
		c.Name = name
		c.Slug = slug
		c.Start = start
		c.End = end
		c.IsActive = isActive
		contests = append(contests, *c)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return contests, nil
}
