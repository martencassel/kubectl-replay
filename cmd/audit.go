package cmd

import (
	"github.com/martencassel/kubectl-replay/replay"
	"github.com/spf13/cobra"
)

func NewAuditCmd() *cobra.Command {
	var file string
	var speed int

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Replay audit logs as kubectl commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return replay.StreamAudit(file, speed)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to audit log file")
	cmd.Flags().IntVar(&speed, "replay-speed", 1, "Replay speed multiplier (e.g. 10 for 10x faster)")
	return cmd
}
