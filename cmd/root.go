package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tok-proxy",
	Short: "TokProxy",
	Long:  "Token Proxy Program",
	Run:   Server,
}

const HttpPortFlag = "http"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	standardFields := log.Fields{
		"hostname": "staging-1",
		"appname":  "tok-proxy",
	}

	log.WithFields(standardFields).Trace("Started init()")

	log.WithFields(standardFields).Info("Configuring serverCmd...")

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

	//serverCmd.Flags().StringVar(&HttpPortNumber, HttpPortFlag, "8080", "HTTP port")
	rootCmd.Flags().String(HttpPortFlag, "8080", "HTTP port")
	err := viper.BindPFlag(HttpPortFlag, rootCmd.Flags().Lookup("http"))
	if err != nil {
		log.Fatalln(err)
	}

}
