/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	mq "nms/messaging"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type SubOptions struct {
	Kind       string
	Topic      string
	RoutingKey string
}

// NewCmdSubscribe create new get cobra command
func NewCmdSub() *cobra.Command {
	o := &SubOptions{}
	var pubCmd = &cobra.Command{
		Use:   "sub",
		Short: "receive message",
		Long: heredoc.Doc(`
			Receive messages with specific topic, please refer --help for details. 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := mq.NewClient(o.Kind) // new client, using amqp protocol as default
			if err != nil {
				cmd.PrintErr(err)
				return
			}
			defer c.Close() // close connection

			cmd.Println("Starting to receive from", o.Kind, "topic", o.Topic)
			cmd.Println("Press Ctrl + C to stop receiving...")
			forever := make(chan bool)
			c.Subscribe(
				o.Topic,
				func(out io.Writer) chan<- string {
					msgs := make(chan string)
					go func() {
						for msg := range msgs {
							fmt.Fprintln(out, msg)
						}
					}()
					return msgs
				}(cmd.OutOrStdout()),
			)
			<-forever
		},
	}

	pubCmd.Flags().StringVarP(&o.Kind, "kind", "k", mq.AMQPKind, "MQ protocol (amqp/mqtt)")
	pubCmd.Flags().StringVarP(&o.Topic, "topic", "t", "nmstest", "Topic")

	//getCmd.MarkFlagRequired("topic")

	return pubCmd
}
