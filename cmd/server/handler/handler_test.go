package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"

	"testing"
)

func TestMetric_HandleCaunter(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		url         string
		requestType string
		metricValue string
		metricType  string
		metricName  string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
				url:         "/update/",
				requestType: "POST",
				metricType:  "gauge/",
				metricName:  "GOOS/",
				metricValue: "123",
			},
		},
		{
			name: "positive test #2",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
				url:         "/update/",
				requestType: "POST",
				metricType:  "counter/",
				metricName:  "GOOS/",
				metricValue: "123",
			},
		},
		{
			name: "negative test #1",
			want: want{
				code: 400,

				contentType: "application/json",
				url:         "/update/",
				requestType: "POST",
				metricType:  "uncown/",
				metricName:  "GOOS/",
				metricValue: "123",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 404,

				contentType: "application/json",
				url:         "/update/",
				requestType: "POST",
				metricType:  "counter/",
				metricName:  "/",
				metricValue: "123",
			},
		},
		{
			name: "negative test #2",
			want: want{
				code: 400,

				contentType: "application/json",
				url:         "/update/",
				requestType: "POST",
				metricType:  "counter/",
				metricName:  "GOOS/",
				metricValue: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.want.url = fmt.Sprintf("%s%s%s%s", test.want.url, test.want.metricType, test.want.metricName, test.want.metricValue)
			request := httptest.NewRequest(test.want.requestType, test.want.url, nil)

			w := httptest.NewRecorder()
			RequestHandle(w, request)

			res := w.Result()
			res.Body.Close()

			assert.Equal(t, res.StatusCode, test.want.code)

		})
	}

}
