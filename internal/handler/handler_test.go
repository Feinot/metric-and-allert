package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Feinot/metric-and-allert/internal/forms"
	"github.com/Feinot/metric-and-allert/internal/logger"
	"github.com/Feinot/metric-and-allert/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

type routeTest struct {
	title                string // title of the test
	val                  string
	types                string // a request to test the route
	host                 string // the expected host of the built URL
	statusCode           int
	name                 string
	mType                string
	valueCounter         int64
	valueGauge           float64
	expectedName         string
	expectedType         string
	expectedValueCounter int64
	expectedvalueGauge   float64
}

func TestMetric_UpdateHandler(t *testing.T) {
	logger.Init()
	tests := []routeTest{
		{
			//pkease talk a little bit about test.title
			title: "Positive test#1 guage ",
			types: "POST",
			mType: "gauge",
			name:  "test",
			val:   "12",

			host:       "/update/gauge/asd/123",
			statusCode: 200,
		},
		{
			title: "Positive test#2 ",

			types: "POST",
			mType: "counter",
			name:  "test",
			val:   "12",

			host:       "/update/counter/asd/123",
			statusCode: 200,
		},
		{
			title: "Nigative test#1 ",
			types: "POST",
			mType: "gauge",
			name:  "test",
			val:   "",

			host:       "http://localhost:8080/update/gauge/asd",
			statusCode: 400,
		},
		{
			title: "Nigative test#2 ",
			types: "POST",

			mType: "gauge",
			name:  "test",
			val:   "",

			host:       "http://localhost:8080/update/counter/asd/",
			statusCode: 400,
		},
		{
			title: "Nigative test#3 ",
			types: "POST",
			mType: "gauge",
			val:   "132",
			name:  "",

			host:       "http://localhost:8080/update/gauge//132",
			statusCode: 404,
		},
		{
			title: "Nigative test#4 ",
			types: "POST",
			name:  "test",
			mType: "",
			val:   "132",

			host:       "http://localhost:8080/update//asd/132",
			statusCode: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.types, test.host, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("Mtype", test.mType)
			rctx.URLParams.Add("Mname", test.name)
			rctx.URLParams.Add("Mvalue", test.val)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			RequestUpdateHandle(w, r)

			assert.Equal(t, test.statusCode, w.Code)

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

				storage.ServerGauge[test.name] = test.valueGauge
			}
			if test.mType == "counter" {

				storage.ServerCounter[test.name] = test.valueCounter
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.types, test.host, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("Mtype", test.mType)
			rctx.URLParams.Add("Mname", test.name)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			RequestValueHandle(w, r)

			assert.Equal(t, test.statusCode, w.Code)

		})
	}
}
func TestMetric_HandleValue(t *testing.T) {
	var metrics forms.Metrics

	tests := []routeTest{
		{
			title:              "Positive test#1 gauge ",
			mType:              "gauge",
			types:              "POST",
			name:               "test1",
			valueGauge:         123.22,
			expectedName:       "test1",
			expectedType:       "gauge",
			expectedvalueGauge: 123.22,
			host:               "http://localhost:8080/value/",
			statusCode:         200,
		},
		{
			title:                "Positive test#2 ",
			mType:                "counter",
			types:                "POST",
			name:                 "test2",
			valueCounter:         123,
			expectedName:         "test2",
			expectedType:         "counter",
			expectedValueCounter: 123,
			host:                 "http://localhost:8080/value/",
			statusCode:           200,
		},
		{
			title:                "Nigative test#1 counter ",
			mType:                "",
			types:                "POST",
			name:                 "test3",
			expectedName:         "test3",
			expectedType:         "counter",
			expectedValueCounter: 123,
			valueCounter:         123,

			host:       "http://localhost:8080/value/",
			statusCode: 400,
		},
		{
			title:              "Nigative test#1_gauge ",
			mType:              "",
			types:              "POST",
			name:               "test4",
			expectedName:       "test4",
			expectedType:       "gauge",
			expectedvalueGauge: 123.123,
			valueGauge:         123.123,

			host:       "http://localhost:8080/value/",
			statusCode: 400,
		},
		{
			title:                "Nigative test#2_counter ",
			mType:                "counter",
			types:                "POST",
			name:                 "",
			expectedName:         "test5",
			expectedType:         "counter",
			expectedValueCounter: 123,
			valueCounter:         123,
			valueGauge:           0,
			host:                 "http://localhost:8080/value/",
			statusCode:           404,
		},
		{
			title:              "Nigative test#2_gauge ",
			mType:              "gauge",
			types:              "POST",
			name:               "",
			expectedName:       "test6",
			expectedType:       "gauge",
			expectedvalueGauge: 123,

			valueGauge: 123,
			host:       "http://localhost:8080/value/",
			statusCode: 404,
		},
		{
			title:                "Nigative test#3 ",
			mType:                "counter",
			types:                "POST",
			name:                 "test7",
			expectedName:         "test8",
			expectedType:         "counter",
			expectedValueCounter: 123,

			host:       "http://localhost:8080/value/",
			statusCode: 404,
		},
		{
			title:              "Nigative test#3 ",
			mType:              "gauge",
			types:              "POST",
			name:               "test9",
			expectedName:       "test10",
			expectedType:       "gauge",
			expectedvalueGauge: 123,

			host:       "http://localhost:8080/value/",
			statusCode: 404,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.mType == "gauge" {

				storage.ServerGauge[test.expectedName] = test.expectedvalueGauge
				metrics.Value = &test.valueGauge
			}
			if test.mType == "counter" {

				storage.ServerCounter[test.expectedName] = test.expectedValueCounter
				metrics.Delta = &test.valueCounter
			}
			metrics.ID = test.name
			metrics.MType = test.mType
			resp, err := json.Marshal(metrics)
			if err != nil {
				fmt.Println("Cannot Marshal: ", err)
			}

			w := httptest.NewRecorder()

			r := httptest.NewRequest(test.types, test.host, bytes.NewReader(resp))

			HandleValue(w, r)

			assert.Equal(t, test.statusCode, w.Code)

		})
	}

}
func TestMetric_RequestValueHandle(t *testing.T) {

	tests := []routeTest{
		{
			//pkease talk a little bit about test.title
			title:              "Positive test#1 guage ",
			types:              "GET",
			mType:              "gauge",
			name:               "test",
			val:                "12",
			expectedName:       "test",
			expectedvalueGauge: 12,
			host:               "http://localhost:8080/value/gauge/test",
			statusCode:         200,
		},
		{
			title: "Positive test#2 ",

			types: "GET",
			mType: "counter",
			name:  "test1",

			expectedName:         "test1",
			expectedValueCounter: 12,
			host:                 "http://localhost:8080/value/counter/test1",
			statusCode:           200,
		},
		{
			title: "Nigative test#1 ",
			types: "GET",
			mType: "gauge",
			name:  "test2",

			expectedName:       "assdasd",
			expectedvalueGauge: 12,

			host:       "http://localhost:8080/value/gauge/asd",
			statusCode: 404,
		},
		{
			title: "Nigative test#2 ",
			types: "GET",

			mType:              "gauge",
			name:               "test3",
			val:                "",
			expectedName:       "fdgdfg",
			expectedvalueGauge: 12,
			host:               "http://localhost:8080/value/counter/asd/",
			statusCode:         404,
		},
		{
			title:              "Nigative test#3 ",
			types:              "GET",
			mType:              "gauge",
			val:                "132",
			name:               "",
			expectedName:       "test4",
			expectedvalueGauge: 12,

			host:       "http://localhost:8080/value/gauge/",
			statusCode: 404,
		},
		{
			title:              "Nigative test#4 ",
			types:              "GET",
			name:               "test5",
			mType:              "",
			val:                "132",
			expectedName:       "test5",
			expectedvalueGauge: 12,
			host:               "http://localhost:8080/value//asd",
			statusCode:         404,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.mType == "gauge" {

				storage.ServerGauge[test.expectedName] = test.expectedvalueGauge

			}
			if test.mType == "counter" {

				storage.ServerCounter[test.expectedName] = test.expectedValueCounter

			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(test.types, test.host, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("Mtype", test.mType)
			rctx.URLParams.Add("Mname", test.name)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			RequestValueHandle(w, r)

			assert.Equal(t, test.statusCode, w.Code)

		})
	}
}
