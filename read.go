package squbix

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var whitespaceNormalizer = regexp.MustCompile(`\s+`)

type queryBuilder struct {
	cteFragments     []string
	selectFragments  []string
	fromFragments    []string
	joinFragments    []string
	whereFragments   []string
	groupByFragments []string
	orderByFragments []string
	limit            *int32
	offset           *int32
}

// NewReadQuery creates new sql builder instance for select operation.
func NewReadQuery(table string) *queryBuilder {
	return &queryBuilder{
		fromFragments: []string{table},
	}
}

// AddCTE adds common table expression to include in generated query.
func (ths *queryBuilder) AddCTE(CTEs ...string) *queryBuilder {
	ths.cteFragments = append(ths.cteFragments, CTEs...)

	return ths
}

// AddSelect adds field to select in generated query.
func (ths *queryBuilder) AddSelect(fields ...string) *queryBuilder {
	ths.selectFragments = append(ths.selectFragments, fields...)

	return ths
}

// AddFrom adds table to select from in generated query.
func (ths *queryBuilder) AddFrom(tables ...string) *queryBuilder {
	ths.fromFragments = append(ths.fromFragments, tables...)

	return ths
}

// AddJoin adds table to join.
func (ths *queryBuilder) AddJoin(tables ...string) *queryBuilder {
	ths.joinFragments = append(ths.joinFragments, tables...)

	return ths
}

// AddWhere adds where clause in generated query.
func (ths *queryBuilder) AddWhere(where ...string) *queryBuilder {
	ths.whereFragments = append(ths.whereFragments, where...)

	return ths
}

// AddGroupBy adds field to group by in generated query.
func (ths *queryBuilder) AddGroupBy(groupBy ...string) *queryBuilder {
	ths.groupByFragments = append(ths.groupByFragments, groupBy...)

	return ths
}

// AddOrderBy adds field to order by in generated query.
func (ths *queryBuilder) AddOrderBy(orderBy ...string) *queryBuilder {
	ths.orderByFragments = append(ths.orderByFragments, orderBy...)

	return ths
}

// AddLimit adds limit in generated query.
func (ths *queryBuilder) AddLimit(limit int32) *queryBuilder {
	ths.limit = &limit

	return ths
}

// AddOffset adds offset in generated query.
func (ths *queryBuilder) AddOffset(offset int32) *queryBuilder {
	ths.offset = &offset

	return ths
}

// BuildQuery generates final query string.
func (ths *queryBuilder) BuildQuery() (string, error) {
	if len(ths.fromFragments) == 0 {
		return "", errors.New("no table specified for query")
	}
	if len(ths.selectFragments) == 0 {
		return "", errors.New("no field selected, add it using AddSelect method")
	}

	queryFragments := []string{}

	if len(ths.cteFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"WITH %s",
			strings.Join(ths.cteFragments, ", "),
		))
	}

	queryFragments = append(queryFragments, fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(ths.selectFragments, ", "),
		strings.Join(ths.fromFragments, ", "),
	))

	if len(ths.joinFragments) > 0 {
		queryFragments = append(queryFragments, strings.Join(ths.joinFragments, " "))
	}

	if len(ths.whereFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"WHERE %s",
			strings.Join(ths.whereFragments, " AND "),
		))
	}

	if len(ths.groupByFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"GROUP BY %s",
			strings.Join(ths.groupByFragments, ", "),
		))
	}

	if len(ths.orderByFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"ORDER BY %s",
			strings.Join(ths.orderByFragments, ", "),
		))
	}

	if ths.limit != nil {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"LIMIT %s",
			fmt.Sprint(*ths.limit),
		))
	}

	if ths.offset != nil {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"OFFSET %s",
			fmt.Sprint(*ths.offset),
		))
	}

	query := strings.Join(queryFragments, " ")
	query = whitespaceNormalizer.ReplaceAllString(query, " ")

	return query, nil
}
