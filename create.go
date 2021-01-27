package squbix

import (
	"errors"
	"fmt"
	"strings"
)

type createQueryBuilder struct {
	cteFragments            []string
	intoFragment            string
	fieldFragments          []string
	valueFragments          []string
	valueWithSelectFragment string
	onConflictFragment      string
}

// NewCreateQuery creates new sql builder instance for insert operation.
func NewCreateQuery(table string) *createQueryBuilder {
	return &createQueryBuilder{
		intoFragment: table,
	}
}

// AddCTE adds common table expression to include in generated query.
func (ths *createQueryBuilder) AddCTE(CTEs ...string) *createQueryBuilder {
	ths.cteFragments = append(ths.cteFragments, CTEs...)

	return ths
}

// AddField adds field to insert in generated query.
func (ths *createQueryBuilder) AddField(fields ...string) *createQueryBuilder {
	ths.fieldFragments = append(ths.fieldFragments, fields...)

	return ths
}

// AddValue adds value to insert in generated query.
func (ths *createQueryBuilder) AddValue(values ...string) *createQueryBuilder {
	ths.valueFragments = append(ths.valueFragments, values...)

	return ths
}

// AddValueWithSelect adds value with select in generated query.
func (ths *createQueryBuilder) AddValueWithSelect(valueWithSelect string) *createQueryBuilder {
	ths.valueWithSelectFragment = valueWithSelect

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
	if len(ths.valueFragments) == 0 && len(ths.valueWithSelectFragment) == 0 {
		return "", errors.New("no value(s) to be inserted, add it using AddValue or AddValueWithSelect method")
	}
	if len(ths.valueFragments) > 0 && len(ths.valueWithSelectFragment) > 0 {
		return "", errors.New("use only AddValue or AddValueWithSelect to add value(s)")
	}

	queryFragments := []string{}

	if len(ths.cteFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"WITH %s",
			strings.Join(ths.cteFragments, ", "),
		))
	}

	queryFragments = append(queryFragments, fmt.Sprintf(
		"INSERT INTO %s (%s)",
		ths.intoFragment,
		strings.Join(ths.fieldFragments, ", "),
	))

	if len(ths.valueFragments) > 0 {
		queryFragments = append(queryFragments, fmt.Sprintf(
			"VALUES %s",
			strings.Join(ths.valueFragments, ", "),
		))
	}

	if len(ths.valueWithSelectFragment) > 0 {
		queryFragments = append(queryFragments, ths.valueWithSelectFragment)
	}

	if len(ths.onConflictFragment) > 0 {
		queryFragments = append(queryFragments, ths.onConflictFragment)
	}

	query := strings.Join(queryFragments, " ")
	query = whitespaceNormalizer.ReplaceAllString(query, " ")

	return query, nil
}
