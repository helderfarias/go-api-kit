package sqlbuilder

import (
	"bytes"
)

type sqlBuilder struct {
	sql     string
	dialect Dialect
}

type sqlBuilderFactory struct {
	dialect Dialect
}

func NewSqlBuilder(dialect Dialect) SqlBuilderFactory {
	return &sqlBuilderFactory{dialect: dialect}
}

func Oracle() Dialect {
	return &oracleDialect{}
}

func Postgres() Dialect {
	return &postgresDialect{}
}

func Select(selectClausule string) SqlBuilder {
	return &sqlBuilder{sql: selectClausule, dialect: Postgres()}
}

func (b sqlBuilder) Where(whereClausule WhereArgs) WhereBuilder {
	values := &valueBuilder{
		args: []interface{}{},
		sql:  bytes.NewBufferString(" WHERE 1=1"),
	}

	whereClausule(values)

	return &whereBuilder{parent: b, args: values.args, sql: values.sql.String()}
}

func (f *sqlBuilderFactory) Select(selectClausule string) SqlBuilder {
	return &sqlBuilder{sql: selectClausule, dialect: f.dialect}
}
