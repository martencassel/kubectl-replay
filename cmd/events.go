package cmd

import (
	"fmt"

	"github.com/martencassel/kubectl-replay/replay"
	"github.com/spf13/cobra"
)

func NewEventsCmd() *cobra.Command {
	var fromEventLog bool
	var replaySpeed string
	var kubeconfig string

	cmd := &cobra.Command{
		Use:   "events",
		Short: "Replay Kubernetes events as kubectl commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse replay speed
			speedMultiplier, err := replay.ParseReplaySpeed(replaySpeed)
			if err != nil {
				return fmt.Errorf("invalid replay speed: %w", err)
			}

			fmt.Printf("Replaying events (fromEventLog=%v) at speed %s...\n", fromEventLog, replaySpeed)

			return replay.StreamLiveEvents(kubeconfig, speedMultiplier)
		},
	}

	cmd.Flags().BoolVar(&fromEventLog, "from-event-log", false, "Replay from event log instead of live cluster")
	cmd.Flags().StringVar(&replaySpeed, "replay-speed", "1x", "Replay speed multiplier (e.g. 10x)")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file (defaults to $KUBECONFIG or ~/.kube/config)")

	return cmd
}
