package repo

import (
	"bytes"
	"strings"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

func applyFilter(filter model.EntryQueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ContestantID != nil {
		wc = append(wc, "contestant_id = ?")
		*args = append(*args, *filter.ContestantID)
	}
	if filter.SessionID != nil {
		wc = append(wc, "session_id = ?")
		*args = append(*args, *filter.SessionID)
	}
	if filter.Status != nil {
		wc = append(wc, "status = ?")
		*args = append(*args, string(*filter.Status))
	}
	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, *filter.ID)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
