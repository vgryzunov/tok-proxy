package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tok-proxy",
	Short: "TokProxy",
	Long:  "Token Proxy Program",
	Run:   Server,
}

const HttpPortFlag = "http"
const OriginUrlFlag = "origin_url"

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	standardFields := log.Fields{
		"hostname": "staging-1",
		"appname":  "tok-proxy",
	}

	log.WithFields(standardFields).Trace("Started init()")

	log.Info("Configuring serverCmd...")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("tp") // TP_ - prefix for environment variables
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config.yaml")

	err2 := viper.ReadInConfig()
	if err2 != nil {
		log.Fatal(err2)
	}

	elementsMap := viper.GetStringMap("server")

	portHttp := elementsMap["http"].(map[string]interface{})["port"]
	log.Printf("HTTP port from %s: %d", viper.ConfigFileUsed(), portHttp)

	originUrl := elementsMap["origin_url"]
	if originUrl != nil {

	}
	log.Println("Origin URL from config file %s", originUrl)

	flags := rootCmd.Flags()

	flags.String(HttpPortFlag, "8080", "HTTP port")
	err := viper.BindPFlag(HttpPortFlag, flags.Lookup("http"))
	if err != nil {
		log.Fatalln(err)
	}

	flags.String(OriginUrlFlag, "http://httpbin.org/", "Origin URL")
	err = viper.BindPFlag(OriginUrlFlag, flags.Lookup("http"))
	if err != nil {
		log.Fatalln(err)
	}

}
