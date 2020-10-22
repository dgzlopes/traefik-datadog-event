// Package traefik_datadog_event can trigger Datadog events when some pattern matches
package traefik_datadog_event

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
)

const endpoint = "https://api.datadoghq.com/api/v1/events?api_key="

// Config the plugin configuration.
type Config struct {
	APIKey      string `yaml:"APIKey"`
	Code        int    `yaml:"Code"`
	BodyPattern string `yaml:"BodyPattern"`
	Title       string `yaml:"Title"`
	Message     string `yaml:"Message"`
	Priority    string `yaml:"Priority"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Code:     -1,
		Title:    "Default error",
		Message:  "Default error",
		Priority: "normal",
	}
}

// DatadogEvent plugin
type DatadogEvent struct {
	next        http.Handler
	APIKey      string
	Code        int
	BodyPattern string
	Title       string
	Message     string
	Priority    string
	name        string
}

// New created a new plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("traefik-datadog-event: you need to specity your Datadog APIKey")
	}
	return &DatadogEvent{
		APIKey:      config.APIKey,
		Code:        config.Code,
		BodyPattern: config.BodyPattern,
		Title:       config.Title,
		Message:     config.Message,
		Priority:    config.Priority,
		next:        next,
		name:        name,
	}, nil
}

func (a *DatadogEvent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	recorder := httptest.NewRecorder()
	a.next.ServeHTTP(recorder, req)
	_, _ = rw.Write(recorder.Body.Bytes())

	if recorder.Code == a.Code {
		if a.BodyPattern != "" {
			if CheckPattern(a.BodyPattern, recorder.Body.String()) {
				SendEvent(a)
			}
		} else {
			SendEvent(a)
		}

	} else {
		if a.BodyPattern != "" && a.Code == -1 {
			fmt.Println("hwllo")
			if CheckPattern(a.BodyPattern, recorder.Body.String()) {
				SendEvent(a)
			}
		}
	}
}

// GenerateEventPayload generates the JSON payload required by the Datadog API
func GenerateEventPayload(a *DatadogEvent) *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"title":"` + a.Title + `","text":"` + a.Message + `","priority":"` + a.Priority + `"}`))
}

// SendEvent sends the HTTP event to Datadog
func SendEvent(a *DatadogEvent) {
	req, err := http.NewRequest("POST", endpoint+a.APIKey, GenerateEventPayload(a))
	if err != nil {
		log.Fatal("traefik-datadog-event: failed to create NewRequest %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	_, ok := httpClient.Do(req)
	if ok != nil {
		log.Fatal("traefik-datadog-event: event request denied %w", err)
	}
}

// CheckPattern checks if some pattern matches using regexp
func CheckPattern(pattern, value string) bool {
	res, _ := regexp.MatchString(pattern, value)
	return res
}
