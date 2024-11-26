package filters

import (
	"fmt"
	"strings"
)

const (
	OperatorEq  = "eq"
	OperatorNeq = "neq"
	OperatorLt  = "lt"
	OperatorGt  = "gt"
)

var OperatorChars = map[string]string{
	OperatorEq:  "=",
	OperatorNeq: "!=",
	OperatorLt:  "<",
	OperatorGt:  ">",
}

// filter for DB (WHERE ...)
type Filter struct {
	Operator string // operators eq/neq/lt/gt
	Name     string // name of the field in database
	Value    any    // value (string or int)
}

type FilterSettings struct {
	PageSize   int
	Pagination int
	F          []Filter
	FieldOrder string
}

// get filtered query fo db
func (s *FilterSettings) GetFilterWithPagination() string {
	res := ""
	for _, f := range s.F {
		char, ok := OperatorChars[f.Operator]
		if ok {
			if res != "" {
				res += " AND "
			} else {
				res += "WHERE "
			}
			if v, ok := f.Value.(string); ok {
				res += fmt.Sprintf(`%s %s '%s'`, f.Name, char, v)
			} else if v, ok := f.Value.(int); ok {
				res += fmt.Sprintf(`%s %s %d`, f.Name, char, v)
			}
		} else {
			continue
		}
	}
	if s.FieldOrder != "" {
		res += fmt.Sprintf(" ORDER BY %s ", s.FieldOrder)
	}
	if s.PageSize > 0 && s.Pagination > 0 {
		itemsPerPage := s.PageSize
		offset := (s.Pagination - 1) * itemsPerPage

		res += fmt.Sprintf(" LIMIT %d OFFSET %d", itemsPerPage, offset)
	}
	return res
}

func GetOperatorAndValue(s string) (string, string) {
	if strings.Contains(s, ":") {
		spl := strings.Split(s, ":")
		return spl[0], spl[1]
	} else {
		return OperatorEq, s
	}
}
