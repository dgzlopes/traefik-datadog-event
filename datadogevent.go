// Package traefik_datadog_event can trigger Datadog events when some pattern matches
package traefik_datadog_event

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
)

const endpoint = "https://api.datadoghq.com/api/v1/events?api_key="

// Config the plugin configuration.
type Config struct {
	APIKey   string `yaml:"APIKey"`
	Code     int    `yaml:"Code"`
	Title    string `yaml:"Title"`
	Message  string `yaml:"Message"`
	Priority string `yaml:"Priority"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Title:    "Default error",
		Message:  "Default error",
		Priority: "normal",
	}
}

// DatadogEvent plugin
type DatadogEvent struct {
	next     http.Handler
	APIKey   string
	Code     int
	Title    string
	Message  string
	Priority string
	name     string
}

// New created a new plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &DatadogEvent{
		APIKey:   config.APIKey,
		Code:     config.Code,
		Title:    config.Title,
		Message:  config.Message,
		Priority: config.Priority,
		next:     next,
		name:     name,
	}, nil
}

func (a *DatadogEvent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	recorder := httptest.NewRecorder()
	a.next.ServeHTTP(recorder, req)
	_, _ = rw.Write(recorder.Body.Bytes())

	if recorder.Code == a.Code {
		req, err := http.NewRequest("POST", endpoint+a.APIKey, generateEventPayload(a))
		if err != nil {
			log.Printf("traefik-datadog-event: failed to create NewRequest %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		httpClient := &http.Client{}
		_, ok := httpClient.Do(req)
		if ok != nil {
			log.Printf("traefik-datadog-event: event request denied %w", err)
		}
	}
}

func generateEventPayload(a *DatadogEvent) *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"title":"` + a.Title + `","text":"` + a.Message + `","priority":"` + a.Priority + `"}`))
}
