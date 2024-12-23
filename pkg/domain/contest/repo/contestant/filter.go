package repo

import (
	"bytes"
	"strings"

	"github.com/pawelkuk/pictura-certamine/pkg/domain/contest/model"
)

func applyFilter(filter model.ContestantQueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, *filter.ID)
	}
	if filter.Email != nil {
		wc = append(wc, "email = ?")
		*args = append(*args, filter.Email.Address)
	}
	if filter.FirstName != nil {
		wc = append(wc, "first_name = ?")
		*args = append(*args, *filter.FirstName)
	}
	if filter.LastName != nil {
		wc = append(wc, "last_name = ?")
		*args = append(*args, *filter.LastName)
	}
	if filter.ConsentConditions != nil {
		wc = append(wc, "consent_conditions = ?")
		*args = append(*args, *filter.ConsentConditions)
	}
	if filter.ConsentConditions != nil {
		wc = append(wc, "consent_marketing = ?")
		*args = append(*args, *filter.ConsentMarketing)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
