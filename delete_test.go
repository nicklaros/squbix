package squbix

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDeleteQuery(t *testing.T) {
	Convey("If we dont specify table to read", t, func() {
		builder := &deleteQueryBuilder{}
		query, err := builder.BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no table specified for query")
			So(query, ShouldEqual, "")
		})
	})

	Convey("Given where clauses", t, func() {
		query, err := NewDeleteQuery("table_a").
			AddWhere(
				"table_a.id = :id",
				"table_a.name = :name",
			).
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "DELETE FROM table_a WHERE table_a.id = :id AND table_a.name = :name")
		})
	})
}
