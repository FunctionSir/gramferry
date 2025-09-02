package main

import "github.com/spf13/cobra"

var UDP, TCP string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gf",
		Short: "Gram Ferry",
	}
	rootCmd.PersistentFlags().StringVarP(&TCP, "tcp", "t", "", "Specific TCP Addr:Port")
	rootCmd.PersistentFlags().StringVarP(&UDP, "udp", "u", "", "Specific UDP Addr:Port")
	_ = rootCmd.MarkFlagRequired("tcp")
	_ = rootCmd.MarkFlagRequired("udp")

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "As the server side",
		Run:   cmdServer,
	}

	var clientCmd = &cobra.Command{
		Use:   "client",
		Short: "As the client side",
		Run:   cmdClient,
	}

	rootCmd.AddCommand(serverCmd, clientCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
