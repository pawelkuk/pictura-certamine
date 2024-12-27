package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samber/lo"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Read(ctx context.Context, e *model.ContestantEntry) error {
	q := `
	SELECT
		e.id,
		e.contestant_id,
		e.status,
		e.updated_at,
		a.id,
		a.created_at,
		a.key,
		c.email,
		c.first_name,
		c.last_name,
		c.phone_number,
		c.consent_conditions,
		c.consent_marketing
	FROM
		entry AS e
	INNER JOIN
		art_piece AS a
	ON
		a.entry_id = e.id
	INNER JOIN
		contestant AS c
	ON
		c.id = e.contestant_id
	WHERE
		e.id = ?
	`
	rows, err := r.DB.QueryContext(ctx, q, e.ID)
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	entries := map[string]model.ContestantEntry{}
	for rows.Next() {
		var eid, econtestantid, estatus, eupdatedat, acreatedat, akey, cemail, cfirstname, clastname, cphone string
		var aid int64
		var cconsentconditions, cconsentmarketing bool
		err := rows.Scan(&eid, &econtestantid, &estatus, &eupdatedat, &aid, &acreatedat, &akey, &cemail, &cfirstname, &clastname, &cphone, &cconsentconditions, &cconsentmarketing)
		if err != nil {
			return fmt.Errorf("could not scan row: %w", err)
		}
		status := estatus
		createdAt, err := time.Parse(time.RFC3339, acreatedat)
		if err != nil {
			return fmt.Errorf("could not parse created at: %w", err)
		}
		p := model.ArtPiece{
			ID:        aid,
			CreatedAt: createdAt,
			Key:       akey,
		}
		if e, ok := entries[eid]; ok {
			e.ArtPieces = append(e.ArtPieces, p)
			entries[eid] = e
		} else {
			e := model.ContestantEntry{
				ID:                eid,
				ContestantID:      econtestantid,
				Status:            status,
				Email:             cemail,
				ArtPieces:         []model.ArtPiece{p},
				FirstName:         cfirstname,
				LastName:          clastname,
				PhoneNumber:       cphone,
				ConsentConditions: cconsentconditions,
				ConsentMarketing:  cconsentmarketing,
				UpdatedAt:         eupdatedat,
			}
			entries[eid] = e
		}
	}
	err = rows.Close()
	if err != nil {
		return fmt.Errorf("could not close rows: %w", err)
	}
	*e = entries[e.ID]
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.ContestantEntryQueryFilter) ([]model.ContestantEntry, error) {
	q := `
	SELECT
		e.id,
		e.contestant_id,
		e.status,
		e.updated_at,
		a.id,
		a.created_at,
		a.key,
		c.email,
		c.first_name,
		c.last_name,
		c.phone_number,
		c.consent_conditions,
		c.consent_marketing
	FROM
		entry AS e
	INNER JOIN
		art_piece AS a
	ON
		a.entry_id = e.id
	INNER JOIN
		contestant AS c
	ON
		c.id = e.contestant_id
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	entries := map[string]model.ContestantEntry{}
	for rows.Next() {
		var eid, econtestantid, estatus, eupdatedat, acreatedat, akey, cemail, cfirstname, clastname, cphone string
		var aid int64
		var cconsentconditions, cconsentmarketing bool
		err := rows.Scan(&eid, &econtestantid, &estatus, &eupdatedat, &aid, &acreatedat, &akey, &cemail, &cfirstname, &clastname, &cphone, &cconsentconditions, &cconsentmarketing)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		status := estatus
		updatedat, err := time.Parse(time.RFC3339, eupdatedat)
		if err != nil {
			return nil, fmt.Errorf("could not parse updated at: %w", err)
		}
		eupdatedat = updatedat.Format(time.DateTime)
		createdAt, err := time.Parse(time.RFC3339, acreatedat)
		if err != nil {
			return nil, fmt.Errorf("could not parse created at: %w", err)
		}
		p := model.ArtPiece{
			ID:        aid,
			CreatedAt: createdAt,
			Key:       akey,
		}
		if e, ok := entries[eid]; ok {
			e.ArtPieces = append(e.ArtPieces, p)
			entries[eid] = e
		} else {
			e := model.ContestantEntry{
				ID:                eid,
				ContestantID:      econtestantid,
				Status:            status,
				Email:             cemail,
				ArtPieces:         []model.ArtPiece{p},
				FirstName:         cfirstname,
				LastName:          clastname,
				PhoneNumber:       cphone,
				ConsentConditions: cconsentconditions,
				ConsentMarketing:  cconsentmarketing,
				UpdatedAt:         eupdatedat,
			}
			entries[eid] = e
		}
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return lo.MapToSlice(entries, func(key string, entry model.ContestantEntry) model.ContestantEntry { return entry }), nil
}
