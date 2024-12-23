package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/mail"

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
			consent_conditions,
			consent_marketing
		) VALUES(?, ?, ?, ?, ?, ?)
		RETURNING id`,
		c.ID,
		c.Email.Address,
		c.FirstName,
		c.LastName,
		c.ConsentConditions,
		c.ConsentMarketing,
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
			consent_conditions,
			consent_marketing
		FROM 
			contestant
		WHERE id = ?`,
		c.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query row with id=%s: %w", c.ID, row.Err())
	}
	var emailStr, firstname, lastName, endStr string
	var cc, cm bool
	err := row.Scan(&emailStr, &firstname, &lastName, &endStr, &cc, &cm)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	email, err := mail.ParseAddress(emailStr)
	if err != nil {
		return fmt.Errorf("could not parse email: %w", err)
	}
	c.Email = *email
	c.FirstName = firstname
	c.LastName = lastName
	c.ConsentConditions = cc
	c.ConsentMarketing = cm
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
			consent_conditions = ?,
			consent_marketing = ?
		WHERE
			id = ?`,
		c.Email.Address,
		c.FirstName,
		c.LastName,
		c.ConsentConditions,
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
		consent_conditions,
		consent_marketing
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
		var cc, cm bool
		err := rows.Scan(&id, &emailStr, &firstname, &lastname, &endStr, &cc, &cm)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		email, err := mail.ParseAddress(emailStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse email: %w", err)
		}
		c.ID = id
		c.Email = *email
		c.FirstName = firstname
		c.LastName = lastname
		c.ConsentConditions = cc
		c.ConsentMarketing = cm
		contestants = append(contestants, *c)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return contestants, nil
}
