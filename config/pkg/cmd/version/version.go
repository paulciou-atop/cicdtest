package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdVersion(version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(Format(version, buildDate))
		},
	}

	return cmd
}

func Format(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		dateStr = fmt.Sprintf(" (%s)", buildDate)
	}

	return fmt.Sprintf("nconfig version %s%s\n", version, dateStr)
}
