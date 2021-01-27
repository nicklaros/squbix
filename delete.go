package squbix

import (
	"errors"
	"fmt"
	"strings"
)

type deleteQueryBuilder struct {
	fromFragment   string
	whereFragments []string
}

// NewDeleteQuery creates new sql builder instance for delete operation.
func NewDeleteQuery(table string) *deleteQueryBuilder {
	return &deleteQueryBuilder{
		fromFragment: table,
	}
}

// AddWhere adds where clause in generated query.
func (ths *deleteQueryBuilder) AddWhere(where ...string) *deleteQueryBuilder {
	ths.whereFragments = append(ths.whereFragments, where...)

	return ths
}

// BuildQuery generates final query string.
func (ths *deleteQueryBuilder) BuildQuery() (string, error) {
	if len(ths.fromFragment) == 0 {
		return "", errors.New("no table specified for query")
	}

	queryFragments := []string{}

	queryFragments = append(queryFragments, fmt.Sprintf(
		"DELETE FROM %s",
		ths.fromFragment,
	))

	if len(ths.whereFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"WHERE %s",
			strings.Join(ths.whereFragments, " AND "),
		))
	}

	query := strings.Join(queryFragments, " ")
	query = whitespaceNormalizer.ReplaceAllString(query, " ")

	return query, nil
}
