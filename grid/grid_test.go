package grid

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMove(t *testing.T) {
	Convey("Init data", t, func() {
		Convey("Init grid data1", func() {
			g := &Grid{Data: [4][4]int{
				{2, 2, 0, 4},
				{2, 2, 4, 0},
				{2, 2, 2, 2},
				{2, 2, 4, 4},
			}}
			Convey("Move Up", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(UP)
				So(result, ShouldEqual, true)
				So(gu.Data, ShouldEqual, [4][4]int{
					{4, 4, 4, 4},
					{4, 4, 2, 2},
					{0, 0, 4, 4},
					{0, 0, 0, 0},
				})
			})
			Convey("Move Down", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(DOWN)
				So(result, ShouldEqual, true)
				So(gu.Data, ShouldEqual, [4][4]int{
					{0, 0, 0, 0},
					{0, 0, 4, 4},
					{4, 4, 2, 2},
					{4, 4, 4, 4},
				})
			})
			Convey("Move Left", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(LEFT)
				So(result, ShouldEqual, true)
				So(gu.Data, ShouldEqual, [4][4]int{
					{4, 4, 0, 0},
					{4, 4, 0, 0},
					{4, 4, 0, 0},
					{4, 8, 0, 0},
				})
			})
			Convey("Move Right", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(RIGHT)
				So(result, ShouldEqual, true)
				So(gu.Data, ShouldEqual, [4][4]int{
					{0, 0, 4, 4},
					{0, 0, 4, 4},
					{0, 0, 4, 4},
					{0, 0, 4, 8},
				})
			})
		})
		Convey("Init grid data2", func() {
			g := &Grid{Data: [4][4]int{
				{2, 4, 2, 4},
				{4, 8, 4, 8},
				{2, 4, 2, 4},
				{4, 8, 4, 8},
			}}
			Convey("Move Up", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(UP)
				So(result, ShouldEqual, false)
			})
			Convey("Move Down", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(DOWN)
				So(result, ShouldEqual, false)
			})
			Convey("Move Left", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(LEFT)
				So(result, ShouldEqual, false)
			})
			Convey("Move Right", func() {
				gu := &Grid{}
				*gu = *g
				result := gu.Move(RIGHT)
				So(result, ShouldEqual, false)
			})
		})
	})
}
