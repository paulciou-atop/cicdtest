package conifg

import (
	"nms/config/pkg/cmd/get/config"

	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <command>",
		Short: "get config or session results ",
	}
	cmd.AddCommand(config.NewCmdConfig())
	return cmd
}
