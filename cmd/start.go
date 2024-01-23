/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vgheri/orchestra/engine"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Runs processes based on the order defined in the orchestra.yml configuration file",
	Long: `Usage: orchestra start
	
	This will start all defined processes in the defined order contained in the orchestra.yml configuration file.
	By default orchestra looks for the file in the current directory ($pwd).
	Use --config flag to specify a different configuration path.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("start called")
		runStart()
	},
}

var verbose bool
var configPath string

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")
	startCmd.PersistentFlags().StringVarP(&configPath, "configPath", "c", "", "use this to specify a custom configuration path")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readConf() (*engine.Configuration, error) {
	var err error
	dir := configPath
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			logFatalMessageAndQuit(fmt.Sprintf("cannot get current directory: %s", err))
		}
	}
	return engine.ReadConfiguration(filepath.Join(dir, "orchestra.yml"))
}

func runStart() {

	// Subscribe early to SIGINT signal
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, os.Interrupt)

	// read configuration
	cfg, err := readConf()
	if err != nil {
		logFatalMessageAndQuit(fmt.Sprintf("cannot read configuration: %s", err))
	}

	done := make(chan engine.ProcessResult, 3)
	defer close(done)

	// start processes
	e := engine.New(cfg, engine.Options{Verbose: verbose}, done)
	err = e.Start()
	if err != nil {
		logFatalMessageAndQuit(fmt.Sprintf("%s", err))
	}

	var resultsRcvd int
	expectedResults := len(cfg.Processes)
	var halt bool
	for {
		select {
		case res := <-done:
			resultsRcvd++
			if res.Err != nil {
				fmt.Printf("[Orchestra] Command %s failed: %s\n", res.ProcessName, res.Err)
			} else {
				fmt.Printf("[Orchestra] Command %s completed\n", res.ProcessName)
			}
		case s := <-quit:
			fmt.Println("[Orchestra] Received interrupt signal, shutting down...")
			e.Stop(s)
			halt = true
		}
		if halt || resultsRcvd == expectedResults {
			break
		}
	}

	fmt.Println("[Orchestra] All done, bye")
	os.Exit(0)
}

func logFatalMessageAndQuit(msg string) {
	fmt.Printf("[Orchestra] Fatal error, %s\n", msg)
	fmt.Println("[Orchestra] Quitting...")
	os.Exit(1)
}
