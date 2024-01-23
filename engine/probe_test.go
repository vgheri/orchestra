package engine

import "testing"

func TestGetRequestSucceeded(t *testing.T) {
	if !getRequestSucceeded(200) {
		t.Fatal("Expected HTTP GET probe request to succeed with HTTP status code 200")
	}
	if !getRequestSucceeded(301) {
		t.Fatal("Expected HTTP GET probe request to succeed with HTTP status code 301")
	}
	if getRequestSucceeded(400) {
		t.Fatal("Expected HTTP GET probe request to fail with HTTP status code 400")
	}
	if getRequestSucceeded(500) {
		t.Fatal("Expected HTTP GET probe request to fail with HTTP status code 500")
	}
}
