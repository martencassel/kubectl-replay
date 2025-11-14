package main

import (
	"fmt"
	"os"

	cmd "github.com/martencassel/kubectl-replay/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kubectl-replay",
		Short: "Replay Kubernetes audit logs and events as kubectl commands",
		Long:  "A kubectl krew plugin that tails Kubernetes audit logs and events, translating them into reproducible kubectl commands with context.",
	}
	// Add subcommands
	rootCmd.AddCommand(cmd.NewAuditCmd())
	rootCmd.AddCommand(cmd.NewEventsCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
