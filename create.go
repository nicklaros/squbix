package squbix

import (
	"errors"
	"fmt"
	"strings"
)

type createQueryBuilder struct {
	intoFragment       string
	fieldFragments     []string
	valueFragments     []string
	onConflictFragment string
}

// NewCreateQuery creates new sql builder instance for insert operation.
func NewCreateQuery(table string) *createQueryBuilder {
	return &createQueryBuilder{
		intoFragment: table,
	}
}

// AddField adds field to insert in generated query.
func (ths *createQueryBuilder) AddField(fields ...string) *createQueryBuilder {
	ths.fieldFragments = append(ths.fieldFragments, fields...)

	return ths
}

// AddValue adds field to insert in generated query.
func (ths *createQueryBuilder) AddValue(values ...string) *createQueryBuilder {
	ths.valueFragments = append(ths.valueFragments, values...)

	return ths
}

// AddOnConflict adds on conflict in generated query.
func (ths *createQueryBuilder) AddOnConflict(onConflict string) *createQueryBuilder {
	ths.onConflictFragment = onConflict

	return ths
}

// BuildQuery generates final query string.
func (ths *createQueryBuilder) BuildQuery() (string, error) {
	if len(ths.intoFragment) == 0 {
		return "", errors.New("no table specified for query")
	}
	if len(ths.fieldFragments) == 0 {
		return "", errors.New("no field specified, add it using AddField method")
	}
	if len(ths.valueFragments) == 0 {
		return "", errors.New("no value to be inserted, add it using AddValue method")
	}

	queryFragments := []string{}

	queryFragments = append(queryFragments, fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		ths.intoFragment,
		strings.Join(ths.fieldFragments, ", "),
		strings.Join(ths.valueFragments, ", "),
	))

	if len(ths.onConflictFragment) > 0 {
		queryFragments = append(queryFragments, ths.onConflictFragment)
	}

	query := strings.Join(queryFragments, " ")
	query = whitespaceNormalizer.ReplaceAllString(query, " ")

	return query, nil
}
