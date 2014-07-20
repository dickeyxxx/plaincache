package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type fakeCache struct{}

func (c *fakeCache) Delete(key string)            {}
func (c *fakeCache) Write(key string, val string) {}
func (c *fakeCache) Read(key string) (string, bool) {
	if key == "/foo" {
		return "cache result", true
	} else {
		return "", false
	}
}

func TestServer(t *testing.T) {
	Convey("With a server", t, func() {
		s := NewServer()
		s.cache = &fakeCache{}
		handler := s.handler()

		Convey("The handler", func() {
			Convey("handles GET requests", func() {
				Convey("With existing key", func() {
					recorder := httptest.NewRecorder()
					req, err := http.NewRequest("GET", "/foo", nil)
					So(err, ShouldBeNil)
					handler.ServeHTTP(recorder, req)
					So(recorder.Code, ShouldEqual, 200)
				})

				Convey("Without existing key", func() {
					recorder := httptest.NewRecorder()
					req, err := http.NewRequest("GET", "/bar", nil)
					So(err, ShouldBeNil)
					handler.ServeHTTP(recorder, req)
					So(recorder.Code, ShouldEqual, 404)
				})
			})

			Convey("handles POST requests", func() {
				recorder := httptest.NewRecorder()
				req, err := http.NewRequest("POST", "/foo", strings.NewReader("foobar"))
				So(err, ShouldBeNil)
				handler.ServeHTTP(recorder, req)
				So(recorder.Code, ShouldEqual, 200)
			})

			Convey("handles DELETE requests", func() {
				recorder := httptest.NewRecorder()
				req, err := http.NewRequest("DELETE", "/foo", nil)
				So(err, ShouldBeNil)
				handler.ServeHTTP(recorder, req)
				So(recorder.Code, ShouldEqual, 200)
			})

			Convey("handles FOOBAR requests", func() {
				recorder := httptest.NewRecorder()
				req, err := http.NewRequest("FOOBAR", "/foo", nil)
				So(err, ShouldBeNil)
				handler.ServeHTTP(recorder, req)
				So(recorder.Code, ShouldEqual, 405)
			})
		})
	})
}
