package repo

import (
	"bytes"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

func applyFilter(filter model.EntryQueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, filter.ID)
	}

}
