package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type fakeListener struct{}

func (f *fakeListener) ListenAndServe() error { return nil }

func TestMain(t *testing.T) {
	Convey("With fake listener", t, func() {
		server = &fakeListener{}

		Convey("it works", func() {
			main()
		})
	})
}
