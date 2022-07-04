package config

import (
	"context"
	"encoding/json"
	"fmt"
	"nms/api/v1/common"
	inv "nms/inventory/pkg/inventory"
	"nms/lib/repo"

	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runHandler(cmd *cobra.Command, args []string) {
	page, err := cmd.Flags().GetInt32("page")
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	size, err := cmd.Flags().GetInt32("size")
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	rep, err := repo.GetRepo(ctx)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	invs, err := inv.ListInventories(ctx, rep.DB(), &common.Pagination{
		Page: page,
		Size: size,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	if len(invs) <= 0 {
		fmt.Println("No result")
		os.Exit(0)
	}
	jsonStr, err := json.MarshalIndent(invs, "", " ")
	os.Stdout.Write(jsonStr)

}

func NewCmdList() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List inventories",
		Run:   runHandler,
	}
	cmd.Flags().Int32P("page", "p", 1, "A unique device id")
	cmd.Flags().Int32P("size", "s", 20, "A path of device, might be a IP address, modbus address etc..")

	return cmd
}
