// Copyright (c) 2015, Alcatel-Lucent Inc.
// All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
// * Neither the name of bambou nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bambou

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError_NewError(t *testing.T) {

	Convey("Given I create a new Error", t, func() {
		e := NewError(42, "a big error")

		Convey("Then Code should be 42", func() {
			So(e.Code, ShouldEqual, 42)
		})

		Convey("Then Message should be 'a big error'", func() {
			So(e.Message, ShouldEqual, "a big error")
		})

	})
}

func TestError_Error(t *testing.T) {

	Convey("Given I create a new Error", t, func() {
		e := NewError(42, "a big error")

		Convey("Then Error() output should be should be '<Error: 42, message: a big error>'", func() {
			So(e.Error(), ShouldEqual, "<Error: 42, message: a big error>")
		})
	})
}

func TestError_FromJSON(t *testing.T) {

	Convey("Given I create a new Error", t, func() {
		e := NewError(42, "a big error")

		Convey("When I unmarshal some data", func() {
			d := `{"property": "prop", "type": "iznogood", "descriptions": [{"title": "oula", "description": "pas bon"}]}`
			json.NewDecoder(bytes.NewBuffer([]byte(d))).Decode(e)

			Convey("Then the Message should be 'iznogood'", func() {
				So(e.Message, ShouldEqual, "iznogood")
			})

			Convey("Then the Property should be 'prop'", func() {
				So(e.Property, ShouldEqual, "prop")
			})

			Convey("Then the lenght of Descriptions should be 1", func() {
				So(len(e.Descriptions), ShouldEqual, 1)
			})

			Convey("When I retrieve the content of Descriptions", func() {
				d := e.Descriptions[0]

				Convey("Then the Title should be 'oula'", func() {
					So(d.Title, ShouldEqual, "oula")
				})

				Convey("Then the Description should be ' bon'", func() {
					So(d.Description, ShouldEqual, "pas bon")
				})

			})
		})
	})
}
