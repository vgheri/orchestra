package engine

import (
	"time"
)

type Liveness struct {
	Kind         string        `yaml:"kind"`
	Spec         LivenessProbe `yaml:"-"`
	InitialDelay time.Duration `yaml:"initialDelay"`
	Retries      int           `yaml:"retries"`
	RetryDelay   time.Duration `yaml:"retryDelay"`
}

type Process struct {
	Name       string        `yaml:"name"`
	Command    string        `yaml:"command"`
	Liveness   *Liveness     `yaml:"liveness"`
	StartDelay time.Duration `yaml:"startDelay"`
}

func (p *Process) HasStartDelay() bool {
	return p.StartDelay.Nanoseconds() > 0
}

func (p *Process) HasInitialDelay() bool {
	return p.Liveness.InitialDelay.Nanoseconds() > 0
}

func (p *Process) HasRetryDelay() bool {
	return p.Liveness.RetryDelay.Nanoseconds() > 0
}

func (p *Process) HasLiveness() bool {
	return p.Liveness != nil
}
