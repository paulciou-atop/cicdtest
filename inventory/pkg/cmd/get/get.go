package config

import (
	"context"
	"encoding/json"
	inv "nms/inventory/pkg/inventory"
	"nms/lib/repo"

	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runHandler(cmd *cobra.Command, args []string) {
	id := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	rep, err := repo.GetRepo(ctx)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	invs, err := inv.GetInventory(rep.DB(), id)
	countInvs := len(invs)
	if countInvs != 1 {
		logrus.Errorf("should get one result but %d", countInvs)
		os.Exit(1)
	}

	jsonStr, err := json.MarshalIndent(invs[0], "", " ")
	os.Stdout.Write(jsonStr)

}

func NewCmdList() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "get [id]",
		Short:   "get inventory",
		Example: "inventory get 60f2d4cc-6f22-4702-bacc-1cd26bade276",
		Args:    cobra.MinimumNArgs(1),
		Run:     runHandler,
	}

	return cmd
}
