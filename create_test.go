package squbix

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCreateQuery(t *testing.T) {
	Convey("If table not specified", t, func() {
		builder := &createQueryBuilder{}
		query, err := builder.BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no table specified for query")
			So(query, ShouldEqual, "")
		})
	})

	Convey("If field names not specified", t, func() {
		query, err := NewCreateQuery("table_a").BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no field specified, add it using AddField method")
			So(query, ShouldEqual, "")
		})
	})

	Convey("If values to be inserted not specified", t, func() {
		query, err := NewCreateQuery("table_a").
			AddField("field_1", "field_2").
			BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no value(s) to be inserted, add it using AddValue or AddValueWithSelect method")
			So(query, ShouldEqual, "")
		})
	})

	Convey("If 1 value inserted", t, func() {
		query, err := NewCreateQuery("table_a").
			AddField("field_1", "field_2").
			AddValue(
				"(value_1, value_2)",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, `INSERT INTO table_a (field_1, field_2) VALUES (value_1, value_2)`)
		})
	})

	Convey("If multiple values inserted", t, func() {
		query, err := NewCreateQuery("table_a").
			AddField("field_1", "field_2").
			AddValue(
				"(value_1, value_2)",
				"(value_3, value_4)",
				"(value_5, value_6)",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, `INSERT INTO table_a (field_1, field_2) VALUES (value_1, value_2), (value_3, value_4), (value_5, value_6)`)
		})
	})

	Convey("If on conflict injected to the query", t, func() {
		query, err := NewCreateQuery("table_a").
			AddField("field_1", "field_2").
			AddValue(
				"(value_1, value_2)",
				"(value_3, value_4)",
				"(value_5, value_6)",
			).
			AddOnConflict(`
				ON CONFLICT ON CONSTRAINT table_a_pkey
				DO UPDATE SET field_2 = EXCLUDED.field_2`,
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, `INSERT INTO table_a (field_1, field_2) VALUES (value_1, value_2), (value_3, value_4), (value_5, value_6) ON CONFLICT ON CONSTRAINT table_a_pkey DO UPDATE SET field_2 = EXCLUDED.field_2`)
		})
	})

	Convey("If AddValue and AddValueWithSelect used in the same time", t, func() {
		query, err := NewCreateQuery("table_a").
			AddField("field_1", "field_2").
			AddValue(
				"(value_1, value_2)",
				"(value_3, value_4)",
				"(value_5, value_6)",
			).
			AddValueWithSelect(`
				SELECT
					a,
					b`,
			).
			BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "use only AddValue or AddValueWithSelect to add value(s)")
			So(query, ShouldEqual, "")
		})
	})

	Convey("CTEs added", t, func() {
		query, err := NewCreateQuery("table_a").
			AddCTE(
				`cte_1 AS (SELECT 1)`,
				`cte_2 AS (SELECT 2)`,
			).
			AddField("field_1", "field_2").
			AddValue(
				"(value_1, value_2)",
				"(value_3, value_4)",
				"(value_5, value_6)",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, `WITH cte_1 AS (SELECT 1), cte_2 AS (SELECT 2) INSERT INTO table_a (field_1, field_2) VALUES (value_1, value_2), (value_3, value_4), (value_5, value_6)`)
		})
	})

	Convey("Input value with AddValueWithSelect", t, func() {
		query, err := NewCreateQuery("table_a").
			AddField("field_1", "field_2").
			AddValueWithSelect(`
				SELECT
					a,
					b`,
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, `INSERT INTO table_a (field_1, field_2) SELECT a, b`)
		})
	})
}
