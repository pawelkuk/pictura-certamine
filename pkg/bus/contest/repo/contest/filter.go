package repo

import (
	"bytes"
	"strings"
	"time"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

func applyFilter(filter model.ContestQueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, *filter.ID)
	}
	if filter.Name != nil {
		wc = append(wc, "name = ?")
		*args = append(*args, *filter.Name)
	}
	if filter.Slug != nil {
		wc = append(wc, "slug = ?")
		*args = append(*args, filter.Slug.Value)
	}
	if filter.Start != nil {
		wc = append(wc, "start_time = ?")
		*args = append(*args, filter.Start.Format(time.RFC3339))
	}
	if filter.End != nil {
		wc = append(wc, "end_time = ?")
		*args = append(*args, filter.End.Format(time.RFC3339))
	}
	if filter.IsActive != nil {
		wc = append(wc, "is_active = ?")
		*args = append(*args, *filter.IsActive)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
