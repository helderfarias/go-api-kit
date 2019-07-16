package sqlbuilder

import (
	"strings"
)

type postgresDialect struct {
	sqlPaginate string
}

func (d *postgresDialect) SetPagination() {
	d.sqlPaginate = " LIMIT ? OFFSET ?"
}

func (d *postgresDialect) ToSQL(sql string) string {
	if strings.TrimSpace(d.sqlPaginate) != "" {
		query := strings.Builder{}
		query.WriteString(sql)
		query.WriteString(d.sqlPaginate)
		return query.String()
	}

	return sql
}
