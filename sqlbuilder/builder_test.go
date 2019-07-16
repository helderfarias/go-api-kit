package sqlbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCreateQueryWithoutWhere(t *testing.T) {
	sql, args := Select("SELECT u.*, c.* FROM users").Build()

	assert.Equal(t, "SELECT u.*, c.* FROM users", sql)
	assert.Nil(t, args)
}

func TestShouldCreateQueryWithWhereEmpty(t *testing.T) {
	sql, args := Select("SELECT u.*, c.* FROM users").Where(func(args Value) {}).Build()

	assert.Equal(t, "SELECT u.*, c.* FROM users WHERE 1=1", sql)
	assert.NotNil(t, args)
	assert.Equal(t, 0, len(args))
}

func TestShouldCreateQueryWithManyArgsPostgres(t *testing.T) {
	sql, args :=
		Select("SELECT u.*, c.* FROM users u INNER JOIN address a on u.id = a.user_id").
			Where(func(args Value) {
				args.Add("AND a.id = ?", 10)
				args.Add("  AND a.id <> 0  ")
			}).
			OrderBy("ORDER BY id DESC").
			SetPaginate(10, 10).
			Build()

	assert.Equal(t, "SELECT u.*, c.* FROM users u INNER JOIN address a on u.id = a.user_id WHERE 1=1 AND a.id = ? AND a.id <> 0 ORDER BY id DESC LIMIT ? OFFSET ?", sql)
	assert.NotNil(t, args)
	assert.Equal(t, 3, len(args))
}

func TestShouldCreateQueryWithManyArgsOracle(t *testing.T) {
	sql, args :=
		NewSqlBuilder(Oracle()).
			Select("SELECT e.* FROM employee e").
			Where(func(args Value) {
				args.Add("AND e.id = ?", 10)
			}).
			OrderBy("ORDER BY hiredate").
			SetPaginate(10, 10).
			Build()

	assert.Equal(t, "SELECT outer.* FROM (SELECT ROWNUM rn, inner.* FROM (SELECT e.* FROM employee e WHERE 1=1 AND e.id = ? ORDER BY hiredate) inner) outer WHERE outer.rn >= ? AND outer.rn <= ?", sql)
	assert.NotNil(t, args)
	assert.Equal(t, 3, len(args))
}
