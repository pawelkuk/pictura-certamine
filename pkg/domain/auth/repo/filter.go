package repo

import (
	"bytes"
	"strings"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
)

func applyFilter(filter model.QueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.UserID != nil {
		wc = append(wc, "user_id = ?")
		*args = append(*args, filter.UserID)
	}

	if filter.Token != nil {
		wc = append(wc, "token = ?")
		*args = append(*args, filter.Token.Value)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
