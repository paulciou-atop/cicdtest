package root

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	getCmd "nms/serviceswatcher/pkg/cmd/get"
	listCmd "nms/serviceswatcher/pkg/cmd/list"
	registerCmd "nms/serviceswatcher/pkg/cmd/register"
	runCmd "nms/serviceswatcher/pkg/cmd/run"
)

var rootCmd = &cobra.Command{
	Use:   "servicewatcher",
	Short: "servicewatcher provide running services' information",
	Long: heredoc.Doc(`
		When running one or more services ( scanner service, snmp service, gwd service, etc.) 
		on a machine you can use this service to specify to list each service information.
		Each service can query this service to figure out IP address and port for each service. 

	`),
	Run: runHelp,
}

func runHelp(cmd *cobra.Command, args []string) {
	host, err := os.Hostname()
	if err == nil {
		fmt.Println("Host :", host)
	}

	cmd.Help()
}

// Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error("Default CLI execute err: ", err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.AddCommand(listCmd.NewCmdList())
	rootCmd.AddCommand(runCmd.NewRunCmd())
	rootCmd.AddCommand(getCmd.NewCmdGet())
	rootCmd.AddCommand(registerCmd.NewCmdRegister())
}

// NewDefaultCommand creates the `snmp` command with default arguments
func NewDefaultCommand() *cobra.Command {

	return nil
}

// NewDefaultCommandWithArgs
//func NewDefaultCommandWithArgs(o Options)
