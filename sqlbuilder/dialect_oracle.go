package sqlbuilder

import (
	"fmt"
	"strings"
)

type oracleDialect struct {
	sqlPaginate string
}

func (d *oracleDialect) SetPagination() {
	d.sqlPaginate = " WHERE outer.rn >= ? AND outer.rn <= ?"
}

func (d *oracleDialect) ToSQL(sql string) string {
	if strings.TrimSpace(d.sqlPaginate) != "" {
		query := strings.Builder{}
		query.WriteString(fmt.Sprintf("SELECT outer.* FROM (SELECT ROWNUM rn, inner.* FROM (%s) inner) outer", sql))
		query.WriteString(d.sqlPaginate)
		return query.String()
	}

	return sql
}
