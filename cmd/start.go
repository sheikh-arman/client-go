package cmd

import (
	"github.com/sheikh-arman/api-server/handler"
	"github.com/spf13/cobra"
)

var (
	//Port is flag to store the default port for http server.
	Port     int
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "start the server on a default port",
		Long: `start the server on a default port ,
				but port can be specify using the port flag`,
		Run: func(cmd *cobra.Command, args []string) {
			handler.StartServer(Port)
		},
	}
)

func init() {
	startCmd.PersistentFlags().IntVarP(&Port, "port", "p", 5050, "default port for http server")
	rootCmd.AddCommand(startCmd)
}
