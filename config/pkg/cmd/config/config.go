package conifg

import (
	"nms/config/pkg/cmd/config/device"

	"github.com/spf13/cobra"
)

func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "configure device",
	}
	cmd.AddCommand(device.NewCmdDevice())
	return cmd
}
