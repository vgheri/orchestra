package engine

import (
	"fmt"
	"net/http"
	"time"
)

type HttpHeader struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type HttpGetProbe struct {
	Path        string        `yaml:"path"`
	Port        string        `yaml:"port"`
	HttpHeaders []*HttpHeader `yaml:"httpHeaders"`
}

type ExecProbe struct {
	Command string `yaml:"command"`
}

type LivenessProbe interface {
	Probe() (bool, error)
}

// Any code greater than or equal to 200 and less than 400 indicates success. Any other code indicates failure.
func getRequestSucceeded(statusCode int) bool {
	return statusCode >= 200 && statusCode < 400
}

func (p HttpGetProbe) Probe() (bool, error) {
	url := fmt.Sprintf("http://localhost:%s/%s", p.Port, p.Path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("HttpGetProbe: cannot create new requst: %w", err)
	}
	for _, h := range p.HttpHeaders {
		req.Header.Add(h.Name, h.Value)
	}
	client := &http.Client{Timeout: time.Duration(3) * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("HttpGetProbe: request failed: %w", err)
	}
	return getRequestSucceeded(res.StatusCode), nil
}

func (p ExecProbe) Probe() (bool, error) {
	cmd := parseCommand(p.Command)
	err := cmd.Start()
	if err != nil {
		return false, fmt.Errorf("Probe, cannot start command %s: %w", p.Command, err)
	}
	err = cmd.Wait()
	if err != nil {
		return false, fmt.Errorf("Probe, command %s failed: %w", p.Command, err)
	}
	return true, nil
}
