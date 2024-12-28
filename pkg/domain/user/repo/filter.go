package repo

import (
	"bytes"
	"strings"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
)

func applyFilter(filter model.QueryFilter, args *[]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		wc = append(wc, "id = ?")
		*args = append(*args, filter.ID)
	}

	if filter.Email != nil {
		wc = append(wc, "email LIKE ?")
		*args = append(*args, filter.Email.Address)
	}

	if filter.AuthorizationToken != nil {
		wc = append(wc, "authorization_token LIKE ?")
		*args = append(*args, filter.AuthorizationToken)
	}

	if filter.ActivationToken != nil {
		wc = append(wc, "activation_token LIKE ?")
		*args = append(*args, filter.ActivationToken)
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
