package engine

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		Command string
		Dir     string
		Name    string
		Args    []string
	}{
		{Command: "ls -l", Dir: "", Name: "ls", Args: []string{"-l"}},
		{Command: "ping", Dir: "", Name: "ping", Args: nil},
		{Command: "ping -t 10 www.google.com", Dir: "", Name: "ping", Args: []string{"-t", "10", "www.google.com"}},
		{Command: "/test/dir/testbin --verbose", Dir: "/test/dir", Name: "/test/dir/testbin", Args: []string{"--verbose"}},
		{Command: "/Users/vgheri/go/src/github.com/vgheri/testapi", Dir: "/Users/vgheri/go/src/github.com/vgheri", Name: "/Users/vgheri/go/src/github.com/vgheri/testapi", Args: nil},
	}

	for _, tc := range tests {
		cmd := parseCommand(tc.Command)
		var args []string
		if len(cmd.Args) > 1 {
			args = cmd.Args[1:]
		}
		if cmd.Dir != tc.Dir || cmd.Args[0] != tc.Name || !reflect.DeepEqual(args, tc.Args) {
			t.Fatalf("Expected directory to be %s, got %s, application to be %s, got %s, expected arguments %s, got %s\n",
				tc.Dir, cmd.Dir, tc.Name, cmd.Args[0], tc.Args, args)
		}
	}
}

func TestStart(t *testing.T) {
	//	// I want to test
	//	// - Given a configuration where Process A -> Process B, Process B only starts after Process A
	//	// - Given a configuration where Process A has a liveness test -> Process B, Process B only starts after Process A's liveness probe completes
	//	var simpleCfg = `
	//
	// processes:
	//   - name: touch A
	//     command: touch a.txt
	//   - name: touch B
	//     command: touch b.txt
	//
	// `
	//
	//	simpleConfig, err := readConf([]byte(simpleCfg))
	//	if err != nil {
	//		t.Fatalf("Expected error to be nil, got %s instead", err.Error())
	//	}
	//	var cfgWithLiveness = `
	//
	// processes:
	//   - name: touch A
	//     command: touch a.txt
	//     liveness:
	//     kind: exec
	//     spec:
	//     command: ls -l
	//     initialDelay: 3s
	//   - name: touch B
	//     command: touch b.txt
	//
	// `
	//
	//	configWithLiveness, err := readConf([]byte(cfgWithLiveness))
	//	if err != nil {
	//		t.Fatalf("Expected error to be nil, got %s instead", err.Error())
	//	}
	//	tests := []struct {
	//		cfg *Configuration
	//	}{
	//		{cfg: simpleConfig},
	//		{cfg: configWithLiveness},
	//	}
	//	for _, tc := range tests {
	//		err := Start(tc.cfg)
	//		if err != nil {
	//			t.Fatalf("Test failed: %s", err)
	//		}
	//		// check a.txt creation date < b.txt creation date
	//	}
}
