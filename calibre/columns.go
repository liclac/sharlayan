package calibre

import (
	"database/sql"
	"strconv"
	"strings"
)

var _ sql.Scanner = &IDs{}

type IDs []int

func (v *IDs) Scan(src interface{}) error {
	draft := IDs{}
	if src == nil {
		*v = draft
		return nil
	}
	for _, s := range strings.Split(src.(string), ",") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		draft = append(draft, i)
	}
	*v = draft
	return nil
}
