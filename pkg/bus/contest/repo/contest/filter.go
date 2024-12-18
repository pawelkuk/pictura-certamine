package repo

import (
	"bytes"

	"github.com/pawelkuk/pictura-certamine/pkg/bus/contest/model"
)

func applyFilter(filter model.QueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, filter.ID)
	}

}
