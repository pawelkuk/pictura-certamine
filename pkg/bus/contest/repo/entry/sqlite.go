package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
	"github.com/samber/lo"
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
			status
		)
		VALUES(?, ?, ?, ?)
		RETURNING id`,
		e.ID,
		e.ContestantID,
		e.SessionID,
		string(e.Status),
	)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("could not rollback: %w", err)
		}
		return fmt.Errorf("could not create entry: %w", err)
	}
	for idx := range e.ArtPieces {
		e.ArtPieces[idx].CreatedAt = time.Now()
	}
	for idx, p := range e.ArtPieces {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO
				art_piece(
					entry_id,
					created_at,
					key
				)
			VALUES(?, ?, ?)
			RETURNING
				id`,
			e.ID,
			p.CreatedAt.Format(time.RFC3339),
			p.Key,
		)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return fmt.Errorf("could not rollback: %w", err)
			}
			return fmt.Errorf("could not create art piece: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return fmt.Errorf("could not rollback: %w", err)
			}
			return fmt.Errorf("could not get last inserted id: %w", err)
		}
		e.ArtPieces[idx].ID = id
	}
	err = tx.Commit()
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("could not rollback: %w", err)
		}
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, e *model.Entry) error {
	row := r.DB.QueryRowContext(ctx,
		`SELECT
			contestant_id,
			session_id,
			status
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
			art_piece
		WHERE
			entry_id = ?`,
		e.ID)

	if err != nil {
		return fmt.Errorf("could not query rows with entry_id=%s: %w", e.ID, err)
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
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("could not rollback: %w", err)
		}
		return fmt.Errorf("could not update entry: %w", err)
	}
	remainingIds := lo.FilterMap(e.ArtPieces, func(ap model.ArtPiece, idx int) (int64, bool) {
		if ap.ID == 0 {
			return -1, false
		} else {
			return ap.ID, true
		}
	})
	remainingIdStr := lo.Reduce(remainingIds, func(acc string, val int64, idx int) string {
		if idx == 0 {
			return strconv.Itoa(int(val))
		} else {
			return acc + ", " + strconv.Itoa(int(val))
		}
	}, "")
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf(`DELETE FROM
			art_piece
		WHERE
			entry_id = ? AND
			id NOT IN (%s)`, remainingIdStr),
		e.ID,
	)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("could not rollback: %w", err)
		}
		return fmt.Errorf("could not delete art piece: %w", err)
	}
	for _, p := range e.ArtPieces {
		for idx := range e.ArtPieces {
			if !e.ArtPieces[idx].CreatedAt.IsZero() {
				continue
			}
			e.ArtPieces[idx].CreatedAt = time.Now()
		}
		if p.ID == 0 {
			result, err := tx.ExecContext(ctx,
				`INSERT INTO
					art_piece(
						entry_id,
						created_at,
						key
					)
				VALUES( ?, ?, ?)
				RETURNING
					id`,
				e.ID,
				p.CreatedAt.Format(time.RFC3339),
				p.Key,
			)
			if err != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					return fmt.Errorf("could not rollback: %w", err)
				}
				return fmt.Errorf("could not create art piece: %w", err)
			}
			id, err := result.LastInsertId()
			if err != nil {
				errRollback := tx.Rollback()
				if errRollback != nil {
					return fmt.Errorf("could not rollback: %w", err)
				}
				return fmt.Errorf("could not get last inserted id: %w", err)
			}
			p.ID = id
		} else {
			_, err := tx.ExecContext(ctx,
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
				errRollback := tx.Rollback()
				if errRollback != nil {
					return fmt.Errorf("could not rollback: %w", err)
				}
				return fmt.Errorf("could not update art piece: %w", err)
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("could not rollback: %w", err)
		}
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, c *model.Entry) error {
	// necessary for on cascade deletion
	_, err := r.DB.ExecContext(ctx, `PRAGMA foreign_keys = ON;`)
	if err != nil {
		return fmt.Errorf("could not turn on foreign keys: %w", err)
	}
	_, err = r.DB.ExecContext(ctx, `DELETE FROM entry WHERE id = ?`, c.ID)
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
		var eid, econtestantid, esessionid, estatus, acreatedat, akey string
		var aid int64
		err := rows.Scan(&eid, &econtestantid, &esessionid, &estatus, &aid, &acreatedat, &akey)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		status, err := model.ParseStatus(estatus)
		if err != nil {
			return nil, fmt.Errorf("could not parse status: %w", err)
		}
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
		} else {
			e := model.Entry{
				ID:           eid,
				ContestantID: econtestantid,
				SessionID:    esessionid,
				Status:       status,
				ArtPieces:    []model.ArtPiece{p},
			}
			entries[eid] = e
		}
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return lo.MapToSlice(entries, func(key string, entry model.Entry) model.Entry { return entry }), nil
}
