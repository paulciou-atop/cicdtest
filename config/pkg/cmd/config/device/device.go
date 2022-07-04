package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nms/config/internal/services"
	"nms/config/pkg/config"
	"nms/config/pkg/session"
	"nms/lib/pgutils"
	"nms/lib/repo"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/MakeNowJust/heredoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func isFileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func readYaml(file string) ([]map[string]interface{}, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	err = yaml.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func readJson(file string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var concreteJson map[string]interface{}
	err = json.Unmarshal(byteValue, &concreteJson)
	if err != nil {
		return nil, err
	}
	return concreteJson, nil
}

func printSessionResult(c config.IConfig, db pgutils.IClient, sid string, devId string) (table.Row, error) {

	configSession, err := session.GetConfigSession(db, sid)
	if err != nil {
		return nil, err
	}
	sessionState := session.UnMarshalSessionState(&configSession)

	return table.Row{sessionState.Id, sessionState.State, devId}, nil

}

var TIMEOUT = config.CONFIG_TIMEOUT

func runHandler(cmd *cobra.Command, args []string) {
	start := time.Now()
	exit := func(err error) {
		elapsed := time.Since(start)
		fmt.Printf("Execution took %s \n", elapsed)
		if err == nil {

			os.Exit(0)
		}
		logrus.Errorf("exit reason: %v", err)
		os.Exit(1)
	}
	tempFileName, err := cmd.Flags().GetString("template-file")
	if err != nil {
		logrus.Errorf("Get flag template-file fail: %v", err)
		exit(err)
	}
	confFileName, err := cmd.Flags().GetString("config-file")
	if err != nil {
		logrus.Errorf("Get flag config-file fail: %v", err)
		exit(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	// Get template
	if !isFileExist(tempFileName) {
		err = fmt.Errorf("Template file %s do not exist", tempFileName)
		fmt.Println(err)
		exit(err)
	}

	template, err := readJson(tempFileName)
	if err != nil {
		logrus.Errorf("read template file %s fail: %v", tempFileName, err)
		exit(err)
	}

	// for every devices
	devsConfig, err := readYaml(confFileName)
	if err != nil {
		logrus.Errorf("read config file %s fail: %v", confFileName, err)
		exit(err)
	}
	var batchArguments []batchArg
	fmt.Printf("Prepare %d configuration...\n", len(devsConfig))

	for i, dev := range devsConfig {
		var d config.Device
		var configMetrics []*config.ConfigMetric
		for k, v := range dev {
			switch k {
			case "device_id":
				id, ok := v.(string)
				if !ok {
					exit(fmt.Errorf("index[%d] device_id should be string but %v", i, v))
				}
				d.ID = id
			case "device_path":
				path, ok := v.(string)
				if !ok {
					exit(fmt.Errorf("index[%d] device_path should be string but %v", i, v))
				}
				d.Path = path
			default:
				p, err := Convert(v, nil)
				if err != nil {
					exit(fmt.Errorf("index[%d] config key=%s should be map but %t", i, k, v))
				}
				payload, ok := p.(map[string]interface{})
				if !ok {
					exit(fmt.Errorf("index[%d] config key=%s covert to map fail", i, k))
				}
				t, ok := template[k]
				if !ok {
					logrus.Infof("template file didn't has kind %s or bad format", k)
					continue
				}
				temp, ok := t.(map[string]interface{})
				if !ok {
					logrus.Infof("template file kind %s has bad format", k)
					continue
				}
				dst := mergeMap(temp, payload)
				m := config.NewConfigMetric("", k, dst)
				configMetrics = append(configMetrics, m)
			}

		}
		batchArguments = append(batchArguments, batchArg{
			dev: d, metrics: configMetrics,
		})
	}

	rep, err := repo.GetRepo(ctx)
	if err != nil {
		exit(err)
	}
	services.InitServices(rep)
	config.InitTable(rep.DB())
	session.InitDatabaseTables(rep.DB())
	c := config.NewConfig(rep)

	r := batchconfig(ctx, c, batchArguments)

	var sessions = map[string]string{}
	for ret := range r {
		if ret.e != nil {
			fmt.Printf("config has error: %v \n", ret.e)
		}
		if ret.sessionid != "" {
			sessions[ret.sessionid] = ret.deviceID
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"session", "state", "device_id"})
	for sid, did := range sessions {
		row, err := printSessionResult(c, rep.DB(), sid, did)
		if err != nil {
			logrus.Errorf("session %s has error: %v", sid, err)
		} else {
			t.AppendRow(row)
		}
	}
	t.Render()
	exit(nil)
}

func NewCmdDevice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "configure device",
		Long: heredoc.Doc(`
		configure specific device with template configuration file and configuration file. service combines template
		configuration file and config file to generate each configuration settings and send to the device.`),
		Run: runHandler,
	}
	cmd.Flags().StringP("template-file", "t", "", "Specify the configuration template file")
	cmd.Flags().StringP("config-file", "c", "", "A http server listen port")
	cmd.MarkFlagRequired("template-file")
	cmd.MarkFlagRequired("config-file")
	return cmd
}
