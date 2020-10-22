package traefik_datadog_event

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTPp(t *testing.T) {
	tests := []struct {
		desc          string
		code          int
		title         string
		message       string
		priority      string
		expEvent      string
		expStatusCode int
	}{
		{
			desc:          "generate event payload",
			code:          400,
			title:         "test",
			message:       "test",
			priority:      "normal",
			expEvent:      "{\"title\":\"test\",\"text\":\"test\",\"priority\":\"normal\"}",
			expStatusCode: http.StatusOK,
		}, {
			desc:          "generate event payload",
			code:          100,
			title:         "Something bad happened!",
			message:       "We need help **here**",
			priority:      "low",
			expEvent:      "{\"title\":\"Something bad happened!\",\"text\":\"We need help **here**\",\"priority\":\"low\"}",
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			cfg := &DatadogEvent{
				Code:     test.code,
				Title:    test.title,
				Message:  test.message,
				Priority: test.priority,
				next:     next,
				name:     "event",
			}
			assert.Equal(t, generateEventPayload(cfg).String(), test.expEvent, "Both events should be equal.")
		})
	}
}
