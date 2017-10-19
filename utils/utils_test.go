package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDiff(t *testing.T) {
	Convey("Init data", t, func() {
		d1 := [4][4]int{
			{2, 2, 0, 4},
			{2, 2, 4, 0},
			{2, 2, 2, 2},
			{2, 2, 4, 4},
		}
		d2 := [4][4]int{
			{2, 2, 0, 4},
			{2, 2, 4, 0},
			{2, 2, 2, 2},
			{2, 2, 4, 4},
		}
		d3 := [4][4]int{
			{2, 2, 0, 4},
			{2, 2, 2, 0},
			{2, 2, 2, 2},
			{2, 2, 4, 4},
		}
		So(Diff(d1, d2), ShouldEqual, false)
		So(Diff(d2, d3), ShouldEqual, true)
	})
}
