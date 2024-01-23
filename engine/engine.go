package engine

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vgheri/orchestra/color"
)

type Options struct {
	Verbose bool
}

type Engine struct {
	commands []*exec.Cmd
	Cfg      *Configuration
	Opts     Options
	Done     chan ProcessResult
}

type ProcessResult struct {
	ProcessName string
	Err         error
}

func New(cfg *Configuration, opts Options, done chan ProcessResult) *Engine {
	return &Engine{
		commands: []*exec.Cmd{},
		Cfg:      cfg,
		Opts:     opts,
		Done:     done,
	}
}

func (e *Engine) verbose() bool {
	return e.Opts.Verbose
}

func (e *Engine) Start() error {
	for _, proc := range e.Cfg.Processes {
		// command setup
		c := color.RandomColor()
		cmd := parseCommand(proc.Command)
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			if e.verbose() {
				fmt.Printf("[Orchestra] Start: cannot create StdoutPipe for command %s: %s.\n", proc.Command, err)
			}
			return fmt.Errorf("failed starting command %s", proc.Command)
		}
		scanner := bufio.NewScanner(cmdReader)
		go func(procName string, scanner *bufio.Scanner, c color.Color) {
			for scanner.Scan() {
				line := scanner.Text()
				printLineToStdOut(c, procName, line)
			}
			err = cmd.Wait()
			if err != nil && e.verbose() {
				fmt.Printf("[Orchestra] Start: cannot wait command %s: %s.\n", procName, err)
			}
			e.Done <- ProcessResult{ProcessName: procName, Err: err}
		}(proc.Name, scanner, c)

		// process start
		if proc.HasStartDelay() {
			time.Sleep(proc.StartDelay)
		}
		err = cmd.Start()
		if err != nil {
			if e.verbose() {
				fmt.Printf("[Orchestra] cannot start command %s: %s.\n", proc.Name, err)
			}
			return fmt.Errorf("failed starting command %s", proc.Name)
		}
		e.commands = append(e.commands, cmd)

		// liveness check
		if proc.HasLiveness() {
			if proc.HasInitialDelay() {
				if e.verbose() {
					fmt.Printf("[Orchestra] Initial delay for command %s, sleeping for %s\n", proc.Name, proc.Liveness.InitialDelay.String())
				}
				time.Sleep(proc.Liveness.InitialDelay)
			}
			var ready bool
			var err error
			// probe needs to execute at least once
			for i := 0; i < proc.Liveness.Retries+1 && !ready; i++ {
				ready, err = proc.Liveness.Spec.Probe()
				if err != nil {
					if e.verbose() {
						fmt.Printf("[Orchestra] Error checking liveness of process %s: %s. \n", proc.Name, err)
					}
					return fmt.Errorf("failed starting command %s", proc.Name)
				}
				if !ready && proc.HasRetryDelay() {
					time.Sleep(proc.Liveness.RetryDelay)
				}
			}
			if !ready {
				if e.verbose() {
					fmt.Printf("[Orchestra] Could not confirm process %s liveness within %d retries.", proc.Name, proc.Liveness.Retries)
				}
				return fmt.Errorf("failed starting command %s", proc.Name)
			} else {
				if e.verbose() {
					fmt.Printf("[Orchestra] Process %s activated and confirmed alive\n", proc.Name)
				}
			}
		}
	}
	return nil
}

func (e *Engine) Stop(s os.Signal) {
	for _, c := range e.commands {
		err := c.Process.Signal(s)
		if err != nil && err != os.ErrProcessDone {
			if e.verbose() {
				fmt.Printf("[Orchestra] Error signaling interrupt to process %s: %s\n", c.Path, err)
			}
		}
	}
}

func printLineToStdOut(c color.Color, procName, line string) {
	s := fmt.Sprintf("[%s] %s", procName, line)
	// fmt.Printf("[%s] %s\n", color.Colorise(c, procName), color.Colorise(c, line))
	fmt.Println(color.Colorise(c, s))
}

// parseCommand parses an input string by identifying the command directory (if any), the command and its arguments.
// By convention, after the last "/", the first word is the command to execute and the next one are the arguments to pass to it.
// It means cases where there are multiple commands onto the same line are not supported.
func parseCommand(in string) *exec.Cmd {
	var dir string
	var cmdArgs []string
	script := strings.Split(in, " ")
	if filepath.Dir(script[0]) != "." {
		dir = filepath.Dir(script[0])
	}
	path := script[0]
	if len(script) > 1 {
		cmdArgs = script[1:]
	}
	cmd := exec.Command(path, cmdArgs...)
	cmd.Dir = dir
	return cmd
}
