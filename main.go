package main

import (
	"os"
	"strings"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed descroot.txt
var RootDesc string

var UDP, TCP string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gf",
		Short: "Gram Ferry",
		Long:  strings.TrimSpace(RootDesc),
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "As the server side",
		Run:   cmdServer,
	}
	serverCmd.Flags().StringVarP(&TCP, "tcp", "t", "", "Specific TCP Addr:Port")
	serverCmd.Flags().StringVarP(&UDP, "udp", "u", "", "Specific UDP Addr:Port")
	_ = serverCmd.MarkFlagRequired("tcp")
	_ = serverCmd.MarkFlagRequired("udp")

	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "As the client side",
		Run:   cmdClient,
	}
	clientCmd.Flags().StringVarP(&TCP, "tcp", "t", "", "Specific TCP Addr:Port")
	clientCmd.Flags().StringVarP(&UDP, "udp", "u", "", "Specific UDP Addr:Port")
	_ = clientCmd.MarkFlagRequired("tcp")
	_ = clientCmd.MarkFlagRequired("udp")

	rootCmd.AddCommand(serverCmd, clientCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
