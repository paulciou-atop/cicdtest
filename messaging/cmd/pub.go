/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	mq "nms/messaging"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type PubOptions struct {
	Kind       string
	Topic      string
	Message    string
	RoutingKey string
}

// NewCmdPublish create new get cobra command
func NewCmdPub() *cobra.Command {
	o := &PubOptions{}
	var pubCmd = &cobra.Command{
		Use:   "pub",
		Short: "Send message",
		Long: heredoc.Doc(`
			Send message with specific topic, please refer --help for details. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := mq.NewClient(o.Kind) // new client, using amqp protocol as default
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			defer c.Close() // close connection
			err = c.Publish(o.Topic, o.Message)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
		},
	}

	pubCmd.Flags().StringVarP(&o.Kind, "kind", "k", mq.AMQPKind, "MQ protocol (amqp/mqtt)")
	pubCmd.Flags().StringVarP(&o.Topic, "topic", "t", "nmstest", "Topic")
	pubCmd.Flags().StringVarP(&o.Message, "message", "m", "Hello World!", "Message")

	//getCmd.MarkFlagRequired("topic")
	//getCmd.MarkFlagRequired("message")

	return pubCmd
}
