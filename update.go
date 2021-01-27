package squbix

import (
	"errors"
	"fmt"
	"strings"
)

type updateQueryBuilder struct {
	intoFragment   string
	setFragments   []string
	whereFragments []string
}

// NewUpdateQuery creates new sql builder instance for update operation.
func NewUpdateQuery(table string) *updateQueryBuilder {
	return &updateQueryBuilder{
		intoFragment: table,
	}
}

// AddSetField adds field to update in generated query.
func (ths *updateQueryBuilder) AddSetField(fields ...string) *updateQueryBuilder {
	ths.setFragments = append(ths.setFragments, fields...)

	return ths
}

// AddWhere adds where clause in generated query.
func (ths *updateQueryBuilder) AddWhere(where ...string) *updateQueryBuilder {
	ths.whereFragments = append(ths.whereFragments, where...)

	return ths
}

// BuildQuery generates final query string.
func (ths *updateQueryBuilder) BuildQuery() (string, error) {
	if len(ths.intoFragment) == 0 {
		return "", errors.New("no table specified for query")
	}
	if len(ths.setFragments) == 0 {
		return "", errors.New("no field specified, add it using AddSetField method")
	}
	if len(ths.whereFragments) == 0 {
		return "", errors.New("no update condition specified, this is DANGEROUS, add it using AddWhere method")
	}

	queryFragments := []string{}

	queryFragments = append(queryFragments, fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		ths.intoFragment,
		strings.Join(ths.setFragments, ", "),
		strings.Join(ths.whereFragments, " AND "),
	))

	query := strings.Join(queryFragments, " ")
	query = whitespaceNormalizer.ReplaceAllString(query, " ")

	return query, nil
}
