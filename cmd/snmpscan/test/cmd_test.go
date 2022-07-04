package cmd_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	scanV1 "nms/api/v1/snmpscan"
	snmp "nms/cmd/snmpscan"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

var rootCmd = &cobra.Command{
	Use:   "test",
	Short: "test command with snmp sub commands",
	Long: heredoc.Doc(`
	test command with snmp sub commands.

	`),
	Run: runHelp,
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func init() {
	rootCmd.AddCommand(snmp.NewCmdGet())
	//rootCmd.AddCommand(snmp.NewScanCmd())
	//rootCmd.AddCommand(snmp.NewCmdWalkAll())

	// start server manually
}

func Test_ExecuteSNMPCommand(t *testing.T) {
	b := new(bytes.Buffer)

	subCommands := rootCmd.Commands()
	for i := range subCommands {
		subCommands[i].SetOut(b)
	}

	rootCmd.SetOut(b)
	rootCmd.SetArgs([]string{"get", "-t", "127.0.0.1", "--service-addr", "127.0.0.1:40051"})
	rootCmd.Execute()

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
		return
	}

	pbV, err := structpb.NewValue("get test")
	if err != nil {
		t.Fatal(err)
		return
	}
	pdu := scanV1.PDU{Value: pbV, Name: "Test"}
	jsonret, err := json.MarshalIndent(&pdu, "", "  ")
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, string(out), string(jsonret), "return is not expected, please make sure test server is running")
}
