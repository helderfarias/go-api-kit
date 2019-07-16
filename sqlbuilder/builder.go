package sqlbuilder

type SqlBuilderFactory interface {
	Select(selectClausule string) SqlBuilder
}

type SqlBuilder interface {
	Where(args WhereArgs) WhereBuilder

	Build() (string, []interface{})
}

type Value interface {
	Add(sql string, args ...interface{})
}

type WhereBuilder interface {
	SetPaginate(offset, limit int64) WhereBuilder

	GroupBy(sql string) WhereBuilder

	OrderBy(sql string) WhereBuilder

	Build() (string, []interface{})
}

type WhereArgs func(args Value)

type Dialect interface {
	SetPagination()

	ToSQL(sql string) string
}
