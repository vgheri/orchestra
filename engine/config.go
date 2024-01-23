package engine

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Processes []*Process `yaml:"processes"`
}

func (l *Liveness) UnmarshalYAML(n *yaml.Node) error {
	type L Liveness
	type T struct {
		*L   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}
	obj := &T{L: (*L)(l)}
	if err := n.Decode(obj); err != nil {
		return err
	}
	switch l.Kind {
	case "httpGet":
		l.Spec = new(HttpGetProbe)
	case "exec":
		l.Spec = new(ExecProbe)
	default:
		panic("kind unknown")
	}
	return obj.Spec.Decode(l.Spec)
}

func ReadConfiguration(path string) (*Configuration, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read configuration, cannot read file: %w", err)
	}
	t, err := readConf(buf)
	return t, err
}

func readConf(content []byte) (*Configuration, error) {
	var conf Configuration
	err := yaml.Unmarshal(content, &conf)
	if err != nil {
		return nil, fmt.Errorf("read configuration, cannot unmarshal content: %w", err)
	}
	return &conf, nil
}

func (cfg *Configuration) String() string {
	var print string
	print += "processes:\n"
	for _, p := range cfg.Processes {
		print += fmt.Sprintf("- name: %s\n", p.Name)
		print += fmt.Sprintf("  command: %s\n", p.Command)
		if p.Liveness != nil {
			print += "  liveness:\n"
			print += fmt.Sprintf("    kind: %s\n", p.Liveness.Kind)
			print += "    spec:\n"
			if p.Liveness.Kind == "exec" {
				p, ok := p.Liveness.Spec.(*ExecProbe)
				if !ok {
					panic("Unexpected error type casting from LivenessProbe to ExecProbe")
				}
				print += fmt.Sprintf("      command: %s\n", p.Command)
			} else {
				p, ok := p.Liveness.Spec.(*HttpGetProbe)
				if !ok {
					panic("Unexpected error type casting from LivenessProbe to HttpGetProbe")
				}
				print += fmt.Sprintf("      path: %s\n", p.Path)
				print += fmt.Sprintf("      port: %s\n", p.Port)
				print += "    httpHeaders:\n"
				for _, h := range p.HttpHeaders {
					print += fmt.Sprintf("      - name: %s\n", h.Name)
					print += fmt.Sprintf("        value: %s\n", h.Value)
				}
			}
			print += fmt.Sprintf("    initialDelay: %s\n", p.Liveness.InitialDelay)
			print += fmt.Sprintf("    retries: %d\n", p.Liveness.Retries)
			print += fmt.Sprintf("    retryDelay: %s\n", p.Liveness.RetryDelay)
		}
	}
	return print
}
