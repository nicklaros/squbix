package squbix

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewReadQuery(t *testing.T) {
	Convey("If we dont specify table to read", t, func() {
		builder := &queryBuilder{}
		query, err := builder.BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no table specified for query")
			So(query, ShouldEqual, "")
		})
	})

	Convey("If we dont specify field to select", t, func() {
		query, err := NewReadQuery("table_a").
			BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no field selected, add it using AddSelect method")
			So(query, ShouldEqual, "")
		})
	})

	Convey("Given many tables to select", t, func() {
		query, err := NewReadQuery("table_a").
			AddFrom("table_b").
			AddSelect(
				"field_a",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a, table_b")
		})
	})

	Convey("Given field to select", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a")
		})
	})

	Convey("Given many fields to select", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
				"field_b",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a, field_b FROM table_a")
		})
	})

	Convey("Given many fields to select while using backtick for multiline string", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
				`JSON_BUILD_OBJECT (
					'price', 1000
				) AS field_b`,
			).
			BuildQuery()

		Convey("It should normalize any whitespace character and returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a, JSON_BUILD_OBJECT ( 'price', 1000 ) AS field_b FROM table_a")
		})
	})

	Convey("Given field to select and join expression", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id")
		})
	})

	Convey("Given field to select and many join expressions", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
				"LEFT JOIN table_c ON table_c.row_id = table_b.row_id",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id LEFT JOIN table_c ON table_c.row_id = table_b.row_id")
		})
	})

	Convey("Given field to select, join expression and where clause", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id")
		})
	})

	Convey("Given field to select, join expression and many where clauses", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
				"table_a.name = :name",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id AND table_a.name = :name")
		})
	})

	Convey("Given field to select, join expression, where clause, and group by field", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id")
		})
	})

	Convey("Given field to select, join expression, where clause, and many group by fields", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
				"name",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id, name")
		})
	})

	Convey("Given field to select, join expression, where clause, group by field and CTE", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			AddCTE(
				"table_expression_a AS (SELECT 0)",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "WITH table_expression_a AS (SELECT 0) SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id")
		})
	})

	Convey("Given field to select, join expression, where clause, group by field and many CTEs", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			AddCTE(
				"table_expression_a AS (SELECT 0)",
				"table_expression_b AS (SELECT 1)",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "WITH table_expression_a AS (SELECT 0), table_expression_b AS (SELECT 1) SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id")
		})
	})

	Convey("Given field to select, join expression, where clause, group by field, CTE and an order by field", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			AddCTE(
				"table_expression_a AS (SELECT 0)",
			).
			AddOrderBy(
				"field_a ASC",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "WITH table_expression_a AS (SELECT 0) SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id ORDER BY field_a ASC")
		})
	})

	Convey("Given field to select, join expression, where clause, group by field, CTE and many order by fields", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			AddCTE(
				"table_expression_a AS (SELECT 0)",
			).
			AddOrderBy(
				"field_a ASC",
				"field_b DESC",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "WITH table_expression_a AS (SELECT 0) SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id ORDER BY field_a ASC, field_b DESC")
		})
	})

	Convey("Given field to select, join expression, where clause, group by field, CTE, order by field and limit", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			AddCTE(
				"table_expression_a AS (SELECT 0)",
			).
			AddOrderBy(
				"field_a ASC",
			).
			AddLimit(10).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "WITH table_expression_a AS (SELECT 0) SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id ORDER BY field_a ASC LIMIT 10")
		})
	})

	Convey("Given field to select, join expression, where clause, group by field, CTE, order by field, limit and offset", t, func() {
		query, err := NewReadQuery("table_a").
			AddSelect(
				"field_a",
			).
			AddJoin(
				"LEFT JOIN table_b ON table_b.id = table_a.id",
			).
			AddWhere(
				"table_a.id = :id",
			).
			AddGroupBy(
				"id",
			).
			AddCTE(
				"table_expression_a AS (SELECT 0)",
			).
			AddOrderBy(
				"field_a ASC",
			).
			AddLimit(10).
			AddOffset(5).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "WITH table_expression_a AS (SELECT 0) SELECT field_a FROM table_a LEFT JOIN table_b ON table_b.id = table_a.id WHERE table_a.id = :id GROUP BY id ORDER BY field_a ASC LIMIT 10 OFFSET 5")
		})
	})
}
