package config

import (
	"context"
	"encoding/json"
	"nms/config/internal/services"
	"nms/config/pkg/config"
	"nms/lib/repo"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runHandler(cmd *cobra.Command, args []string) {
	devid, err := cmd.Flags().GetString("device-id")
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	devpath, err := cmd.Flags().GetString("device-path")
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	proto, err := cmd.Flags().GetString("protocol")
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
	services.InitServices(rep)
	c := config.NewConfig(rep)
	res, err := c.Upload(ctx, config.Device{ID: devid, Path: devpath}, proto, []string{})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	jsonStr, err := json.MarshalIndent(res.Payload, "", " ")
	os.Stdout.Write(jsonStr)

}

func NewCmdConfig() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "config ",
		Short: "Get configuration from device",
		Run:   runHandler,
	}
	cmd.Flags().StringP("device-id", "d", "", "A unique device id")
	cmd.Flags().StringP("device-path", "a", "", "A path of device, might be a IP address, modbus address etc..")
	cmd.Flags().StringP("protocol", "p", "", "A config protocol ")
	cmd.MarkFlagRequired("device-id")
	cmd.MarkFlagRequired("device-path")
	cmd.MarkFlagRequired("protocol")

	return cmd
}
