package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tok-proxy",
	Short: "TokProxy",
	Long:  `Token Proxy Program`,
}

const HttpPortFlag = "http"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetOutput(os.Stdout)

	log.Println("Started init()")

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "start",
		Long:  "start token proxy server",
		Run: func(cmd *cobra.Command, agrs []string) {
			log.Printf("inside server command")
			log.Printf(cmd.Short)

			log.Print("Config file used: ", viper.ConfigFileUsed())
			port := viper.GetInt(HttpPortFlag)
			log.Printf("Using HTTP Port: %d", port)
		},
	}

	log.Println("Configuring serverCmd...")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("tp") // TP_ - prefix for environment variables
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")

	err2 := viper.ReadInConfig()
	if err2 != nil {
		log.Fatal(err2)
	}

	elementsMap := viper.GetStringMap("server")
	for k, vMap := range elementsMap {
		log.Print("Key: ", k)
		log.Println(" Value: ", vMap)
	}

	portHttp := elementsMap["http"].(map[string]interface{})["port"]
	log.Printf("HTTP port from %s: %d", viper.ConfigFileUsed(), portHttp)

	portHttps := elementsMap["https"].(map[string]interface{})["port"]
	log.Printf("HTTPS port from %s: %d", viper.ConfigFileUsed(), portHttps)

	rootCmd.AddCommand(serverCmd)

	//serverCmd.Flags().StringVar(&HttpPortNumber, HttpPortFlag, "8080", "HTTP port")
	serverCmd.Flags().String(HttpPortFlag, "8080", "HTTP port")
	err := viper.BindPFlag(HttpPortFlag, serverCmd.Flags().Lookup("http"))
	if err != nil {
		log.Fatalln(err)
	}

}
