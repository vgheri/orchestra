package engine

import (
	"reflect"
	"testing"
	"time"
)

func TestReadConf(t *testing.T) {
	var data = `
processes:
- name: DB
  command: ./bin/db
  liveness:
    kind: exec
    spec:
      command: pg_isready
    initialDelay: 3s
    retries: 3
    retryDelay: 10s
- name: API
  command: ./bin/api
  liveness:
    kind: httpGet 
    spec: 
      path: /health
      port: 8080
      httpHeaders:
      - name: Custom-Header
        value: Awesome
    initialDelay: 3s
    retries: 3
    retryDelay: 10s
- name: Website
  command: npm run start
  startDelay: 1m
`

	expectedProcesses := []*Process{
		{
			Name:     "DB",
			Command:  "./bin/db",
			Liveness: &Liveness{Kind: "exec", Spec: &ExecProbe{Command: "pg_isready"}, Retries: 3, InitialDelay: time.Duration(3) * time.Second, RetryDelay: time.Duration(10) * time.Second},
		},
		{
			Name:     "API",
			Command:  "./bin/api",
			Liveness: &Liveness{Kind: "httpGet", Spec: &HttpGetProbe{Path: "/health", Port: "8080", HttpHeaders: []*HttpHeader{{Name: "Custom-Header", Value: "Awesome"}}}, Retries: 3, InitialDelay: time.Duration(3) * time.Second, RetryDelay: time.Duration(10) * time.Second},
		},
		{
			Name:       "Website",
			Command:    "npm run start",
			StartDelay: time.Duration(1) * time.Minute,
		},
	}
	conf, err := readConf([]byte(data))
	if err != nil {
		t.Fatalf("Expected error to be nil, got %s instead", err.Error())
	}

	if len(conf.Processes) != len(expectedProcesses) {
		t.Fatalf("Expected to have %d processes, got %d instead", len(expectedProcesses), len(conf.Processes))
	}
	for i, p := range conf.Processes {
		if !reflect.DeepEqual(p, expectedProcesses[i]) {
			t.Fatalf("Expected processes %+v, got %+v instead", p, expectedProcesses[i])
		}
	}
}
