package root

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func rootFlagErrorHandler(cmd *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}
	return err
}
