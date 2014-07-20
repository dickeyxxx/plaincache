package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	Convey("With a cache", t, func() {
		c := NewCache()

		Convey("Cache miss", func() {
			_, ok := c.Read("foo")
			So(ok, ShouldEqual, false)
		})

		Convey("With cache in \"foo\"", func() {
			c.Write("foo", "bar")

			Convey("Cache read", func() {
				result, _ := c.Read("foo")
				So(result, ShouldEqual, "bar")
			})

			Convey("Cache delete", func() {
				c.Delete("foo")
				_, ok := c.Read("foo")
				So(ok, ShouldEqual, false)
			})
		})

		Convey("With many caches", func() {
			c.Write("foo", "bar1")
			c.Write("foo", "bar2")
			c.Write("bar", "baz1")
			c.Write("bar", "baz2")
			c.Write("foo", "bar3")

			Convey("Cache read", func() {
				result, _ := c.Read("foo")
				So(result, ShouldEqual, "bar3")
			})
		})
	})
}
