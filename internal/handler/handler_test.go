package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Feinot/metric-and-allert/internal/storage"
	"github.com/stretchr/testify/assert"
)

type routeTest struct {
	title string // title of the test
	//route           *Route            // the route being tested
	types           string            // a request to test the route
	vars            map[string]string // the expected vars of the match
	scheme          string            // the expected scheme of the built URL
	host            string            // the expected host of the built URL
	path            string            // the expected path of the built URL
	query           string            // the expected query string of the built URL
	pathTemplate    string            // the expected path template of the route
	hostTemplate    string            // the expected host template of the route
	queriesTemplate string            // the expected query template of the route
	methods         []string          // the expected route methods
	pathRegexp      string            // the expected path regexp
	queriesRegexp   string            // the expected query regexp
	shouldMatch     bool              // whether the request is expected to match the route at all
	shouldRedirect  bool              // whether the request should result in a redirect
	statusCode      int
	name            string
	mType           string
	valueCounter    int64
	valueGauge      float64
}

func TestMetric_UpdateHandler(t *testing.T) {
	tests := []routeTest{
		{
			title: "Positive test#1 guage ",
			types: "POST",

			host:       "http://localhost:8080/update/gauge/asd/123",
			statusCode: 200,
		},
		{
			title: "Positive test#2 ",
			types: "POST",

			host:       "http://localhost:8080/update/gauge/asd/123",
			statusCode: 200,
		},
		{
			title: "Nigative test#1 ",
			types: "POST",

			host:       "http://localhost:8080/update/gauge/asd/",
			statusCode: 400,
		},
		{
			title: "Nigative test#2 ",
			types: "POST",

			host:       "http://localhost:8080/update/gauge/asd/",
			statusCode: 400,
		},
		{
			title: "Nigative test#3 ",
			types: "POST",

			host:       "http://localhost:8080/update/gauge//132",
			statusCode: 404,
		},
		{
			title: "Nigative test#4 ",
			types: "POST",

			host:       "http://localhost:8080/update//asd/132",
			statusCode: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {

			r, err := http.NewRequest(test.types, test.host, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			w := httptest.NewRecorder()

			//Hack to try to fake gorilla/mux vars

			// CHANGE THIS LINE!!!

			RequestUpdateHandle(w, r)

			assert.Equal(t, test.statusCode, w.Code)
			//assert.Equal(t, []byte("abcd"), w.Body.Bytes())
		})
	}
}
func TestMetric_ValueHandler(t *testing.T) {
	tests := []routeTest{
		{
			title:      "Positive test#1 gauge ",
			mType:      "gauge",
			types:      "GET",
			name:       "asd",
			valueGauge: 123.22,
			host:       "http://localhost:8080/value/gauge/asd",
			statusCode: 200,
		},
		{
			title:        "Positive test#2 ",
			mType:        "counter",
			types:        "GET",
			name:         "asd",
			valueCounter: 123,
			host:         "http://localhost:8080/value/counter/asd",
			statusCode:   200,
		},
		{
			title: "Nigative test#1 ",
			types: "GET",
			mType: "gauge",

			name:       "asd",
			valueGauge: 0,

			host:       "http://localhost:8080/value/gauge/",
			statusCode: 404,
		},
		{
			title: "Nigative test#2 ",
			types: "GET",
			mType: "gauge",

			name:       "asd",
			valueGauge: 0,
			host:       "http://localhost:8080/value/gauge/qwe",
			statusCode: 404,
		},
		{
			title: "Nigative test#3 ",
			types: "GET",

			host:       "http://localhost:8080/value/none",
			statusCode: 404,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.mType == "gauge" {

				storage.Gauge[test.name] = test.valueGauge
			}
			if test.mType == "counter" {

				storage.Counter[test.name] = test.valueCounter
			}

			r, err := http.NewRequest(test.types, test.host, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			w := httptest.NewRecorder()

			//Hack to try to fake gorilla/mux vars

			// CHANGE THIS LINE!!!

			RequestValueHandle(w, r)

			assert.Equal(t, test.statusCode, w.Code)
			//assert.Equal(t, []byte("abcd"), w.Body.Bytes())
		})
	}
}
