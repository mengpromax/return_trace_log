package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "return_trace_log",
	Short: "return_trace_log is a return trace log tool for go language",
	Long:  "return_trace_log is a return trace log tool for go language",
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
	}
}
