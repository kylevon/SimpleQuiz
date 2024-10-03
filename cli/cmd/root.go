package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "quiz",
    Short: "A simple quiz CLI",
    Long:  `A CLI that interacts with a quiz API`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
    }
}
