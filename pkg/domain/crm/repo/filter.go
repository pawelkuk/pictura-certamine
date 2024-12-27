package repo

import (
	"bytes"
	"strings"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/crm/model"
)

func applyFilter(filter model.ContestantEntryQueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	// Contestant
	if filter.ID != nil {
		wc = append(wc, "c.id = ?")
		*args = append(*args, *filter.ID)
	}
	if filter.Email != nil {
		wc = append(wc, "c.email = ?")
		*args = append(*args, filter.Email)
	}
	if filter.FirstName != nil {
		wc = append(wc, "c.first_name = ?")
		*args = append(*args, *filter.FirstName)
	}
	if filter.LastName != nil {
		wc = append(wc, "c.last_name = ?")
		*args = append(*args, *filter.LastName)
	}
	if filter.ConsentConditions != nil {
		wc = append(wc, "c.consent_conditions = ?")
		*args = append(*args, *filter.ConsentConditions)
	}
	if filter.ConsentConditions != nil {
		wc = append(wc, "c.consent_marketing = ?")
		*args = append(*args, *filter.ConsentMarketing)
	}

	// Entry
	if filter.ContestantID != nil {
		wc = append(wc, "e.contestant_id = ?")
		*args = append(*args, *filter.ContestantID)
	}
	if filter.Status != nil {
		wc = append(wc, "e.status = ?")
		*args = append(*args, *filter.Status)
	}
	if filter.ID != nil {
		wc = append(wc, "e.id = ?")
		*args = append(*args, *filter.ID)
	}
	if filter.Token != nil {
		wc = append(wc, "e.token = ?")
		*args = append(*args, *filter.Token)
	}

	if filter.TokenExpiry != nil {
		wc = append(wc, "e.token_expiry = ?")
		*args = append(*args, filter.TokenExpiry.Format(time.RFC3339))
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
