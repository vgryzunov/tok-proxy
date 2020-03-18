package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)


func Server (cmd *cobra.Command, agrs []string) {
	log.Printf("inside server command")
	log.Printf(cmd.Short)

	log.Print("Config file used: ", viper.ConfigFileUsed())
	port := viper.GetInt(HttpPortFlag)
	log.Printf("Using HTTP Port: %d", port)
}