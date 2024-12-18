package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, e *model.Entry) error {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO entry(
			id,
			contestant_id,
			session_id,
			status,
		) VALUES(?, ?, ?, ?,)
		RETURNING id`,
		e.ID,
		e.ContestantID,
		e.SessionID,
		string(e.Status),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not create entry: %w", err)
	}
	for _, p := range e.ArtPieces {
		p.CreatedAt = time.Now()
		result, err := tx.ExecContext(ctx,
			`INSERT INTO
				art_piece(
					entry_id,
					created_at,
					key
				)
			VALUES(?, ?, ?, ?,)
			RETURNING
				id`,
			e.ID,
			p.CreatedAt.Format(time.RFC3339),
			p.Key,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not create art piece: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not get last inserted id: %w", err)
		}
		p.ID = id
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, e *model.Entry) error {
	row := r.DB.QueryRowContext(ctx,
		`SELECT
			contestant_id,
			session_id,
			status,
		FROM 
			entry
		WHERE id = ?`,
		e.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query row with id=%s: %w", e.ID, row.Err())
	}
	var contestantid, sessionid, statusStr string
	err := row.Scan(&contestantid, &sessionid, &statusStr)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	status, err := model.ParseStatus(statusStr)
	if err != nil {
		return fmt.Errorf("could not parse status: %w", err)
	}
	e.ContestantID = contestantid
	e.SessionID = sessionid
	e.Status = status
	rows, err := r.DB.QueryContext(ctx,
		`SELECT
			id,
			created_at,
			key
		FROM
			entry
		WHERE entry_id = ?`,
		e.ID)

	if err != nil {
		return fmt.Errorf("could not query rows with entry_id=%s: %w", e.ID, row.Err())
	}
	for rows.Next() {
		p := model.ArtPiece{}
		var createdAtStr, key string
		var id int64
		err := rows.Scan(&id, &createdAtStr, &key)
		if err != nil {
			return fmt.Errorf("could not scan row: %w", err)
		}
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return fmt.Errorf("could not parse created at: %w", err)
		}
		p.ID = id
		p.Key = key
		p.CreatedAt = createdAt
		e.ArtPieces = append(e.ArtPieces, p)
	}
	err = rows.Close()
	if err != nil {
		return fmt.Errorf("could not close rows: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, e *model.Entry) error {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE 
			entry 
		SET 
			contestant_id = ?,
			session_id = ?,
			status = ?
		WHERE
			id = ?`,
		e.ContestantID,
		e.SessionID,
		string(e.Status),
		e.ID,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not update entry: %w", err)
	}
	for _, p := range e.ArtPieces {
		p.CreatedAt = time.Now()
		_, err := r.DB.ExecContext(ctx,
			`UPDATE
				art_piece
			SET
				entry_id = ?,
				created_at = ?,
				key = ?
			WHERE
				id = ?`,
			e.ID,
			p.CreatedAt.Format(time.RFC3339),
			p.Key,
			p.ID,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("could not update art piece: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, c *model.Entry) error {
	// Art pieces are cascade deleted
	_, err := r.DB.ExecContext(ctx, `DELETE FROM entry WHERE id = ?`, c.ID)
	if err != nil {
		return fmt.Errorf("could not delete entry: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.EntryQueryFilter) ([]model.Entry, error) {
	q := `
	SELECT
		entry.id,
		entry.contestant_id,
		entry.session_id,
		entry.status,
		art_pice.id,
		art_pice.created_at,
		art_pice.key
	FROM
		entry
	INNER JOIN
		art_piece
	ON
		art_pice.entry_id = entry.id
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	entries := map[string]model.Entry{}
	for rows.Next() {
		var eid, econtestantid, esessionid, estatus, aid, acreatedat, akey string
		err := rows.Scan(&eid, &econtestantid, &esessionid, &estatus, &aid, &acreatedat, &akey)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}

		entries = append(entries, *e)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return entries, nil
}
