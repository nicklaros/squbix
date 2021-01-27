package squbix

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewUpdateQuery(t *testing.T) {
	Convey("If we dont specify table to read", t, func() {
		builder := &updateQueryBuilder{}
		query, err := builder.BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no table specified for query")
			So(query, ShouldEqual, "")
		})
	})

	Convey("If we dont specify field to select", t, func() {
		query, err := NewUpdateQuery("table_a").
			BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no field specified, add it using AddSetField method")
			So(query, ShouldEqual, "")
		})
	})

	Convey("Given field to set but no update condition", t, func() {
		query, err := NewUpdateQuery("table_a").
			AddSetField("field_a = 'A'").
			BuildQuery()

		Convey("It should returns error", func() {
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "no update condition specified, this is DANGEROUS, add it using AddWhere method")
			So(query, ShouldEqual, "")
		})
	})

	Convey("Given field to set and update condition", t, func() {
		query, err := NewUpdateQuery("table_a").
			AddSetField("field_a = 'A'").
			AddWhere("field_a IS NULL").
			BuildQuery()

		Convey("It should returns generated query", func() {
			So(err, ShouldBeNil)
			So(query, ShouldEqual, "UPDATE table_a SET field_a = 'A' WHERE field_a IS NULL")
		})
	})
}
