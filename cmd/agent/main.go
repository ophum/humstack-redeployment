package main

import (
	"flag"
	"log"
	"os"

	"github.com/ophum/humstack-redeployment/pkg/agent"
	"github.com/ophum/humstack-redeployment/pkg/client"
	hsClient "github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	config agent.Config = agent.Config{}
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config path")
	flag.Parse()

	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal("error open config file")
	}
	defer configFile.Close()

	if err := yaml.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatal("failed decode file")
	}
}
func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	client := client.NewRedeploymentClient("http", config.APIServerAddress, config.APIServerPort)
	humstackClient := hsClient.NewClients(config.HumstackAPIServerAddress, config.HumstackAPIServerPort)
	agent := agent.NewRedeploymentAgent(client, humstackClient, &config, logger.With(zap.Namespace("RedeploymentAgent")))

	agent.Run()
}
