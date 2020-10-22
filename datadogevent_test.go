package traefik_datadog_event

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeHTTPp(t *testing.T) {
	tests := []struct {
		desc          string
		codePattern   int
		bodyPattern   string
		title         string
		message       string
		priority      string
		expEvent      string
		expStatusCode int
	}{
		{
			desc:          "generate event payload",
			codePattern:   400,
			bodyPattern:   "[a-z]+",
			title:         "test",
			message:       "test",
			priority:      "normal",
			expEvent:      "{\"title\":\"test\",\"text\":\"test\",\"priority\":\"normal\"}",
			expStatusCode: http.StatusOK,
		}, {
			desc:          "generate event payload",
			codePattern:   100,
			bodyPattern:   "[a-z]+",
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
				CodePattern: test.codePattern,
				BodyPattern: test.bodyPattern,
				Title:       test.title,
				Message:     test.message,
				Priority:    test.priority,
				next:        next,
				name:        "event",
			}

			exampleRequestBody := "This is my request body"
			exampleIncorrectRequestBody := "12354 514 251241"

			assert.Equal(t, GenerateEventPayload(cfg).String(), test.expEvent, "Both events should be equal.")
			assert.True(t, CheckPattern(cfg.BodyPattern, exampleRequestBody), "Pattern should match.")
			assert.False(t, CheckPattern(cfg.BodyPattern, exampleIncorrectRequestBody), "Pattern shouldn't match.")
		})
	}
}
