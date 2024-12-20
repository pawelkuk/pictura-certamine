package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, c *model.Contestant) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO contestant(
			id,
			email,
			first_name,
			last_name,
			birthdate,
			policy_accepted
		) VALUES(?, ?, ?, ?, ?, ?)
		RETURNING id`,
		c.ID,
		c.Email.Address,
		c.FirstName,
		c.LastName,
		c.Birthdate.Format(time.DateOnly),
		c.PolicyAccepted,
	)
	if err != nil {
		return fmt.Errorf("could not create contestant: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, c *model.Contestant) error {
	row := r.DB.QueryRowContext(ctx,
		`SELECT 
			email,
			first_name,
			last_name,
			birthdate,
			policy_accepted
		FROM 
			contestant
		WHERE id = ?`,
		c.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query row with id=%s: %w", c.ID, row.Err())
	}
	var emailStr, firstname, lastName, endStr string
	var isActive bool
	err := row.Scan(&emailStr, &firstname, &lastName, &endStr, &isActive)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return fmt.Errorf("could not parse timestamp: %w", err)
	}
	email, err := mail.ParseAddress(emailStr)
	if err != nil {
		return fmt.Errorf("could not parse email: %w", err)
	}
	c.Email = *email
	c.FirstName = firstname
	c.LastName = lastName
	c.Birthdate = end
	c.PolicyAccepted = isActive
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, c *model.Contestant) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE 
			contestant 
		SET 
			email = ?,
			first_name = ?,
			last_name = ?,
			birthdate = ?,
			policy_accepted = ?
		WHERE
			id = ?`,
		c.Email.Address,
		c.FirstName,
		c.LastName,
		c.Birthdate.Format(time.DateOnly),
		c.PolicyAccepted,
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update contestant: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, c *model.Contestant) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM contestant WHERE id = ?`, c.ID)
	if err != nil {
		return fmt.Errorf("could not delete contestant: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.ContestantQueryFilter) ([]model.Contestant, error) {
	q := `
	select
		id,
		email,
		first_name,
		last_name,
		birthdate,
		policy_accepted
	from
		contestant
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	contestants := []model.Contestant{}
	for rows.Next() {
		c := &model.Contestant{}
		var id, emailStr, firstname, lastname, endStr string
		var isActive bool
		err := rows.Scan(&id, &emailStr, &firstname, &lastname, &endStr, &isActive)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse timestamp: %w", err)
		}
		email, err := mail.ParseAddress(emailStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse email: %w", err)
		}
		c.ID = id
		c.Email = *email
		c.FirstName = firstname
		c.LastName = lastname
		c.Birthdate = end
		c.PolicyAccepted = isActive
		contestants = append(contestants, *c)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return contestants, nil
}
