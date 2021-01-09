package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ophum/humstack-redeployment/pkg/api"
	v0 "github.com/ophum/humstack-redeployment/pkg/api/v0"
	store "github.com/ophum/humstack/pkg/store/leveldb"
)

var (
	listenAddress string
	listenPort    int64
	isDebug       bool
)

func init() {
	flag.StringVar(&listenAddress, "listen-address", "localhost", "listen address")
	flag.Int64Var(&listenPort, "listen-port", 8090, "listen port")
	flag.BoolVar(&isDebug, "debug", false, "debug mode true/false")
	flag.Parse()
}

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	notifier := make(chan string, 100)
	s, err := store.NewLevelDBStore("./database", notifier, isDebug)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	go func() {
		for range notifier {
		}
	}()

	rdh := v0.NewRedeploymentHandler(s)

	v0 := r.Group("/api/v0")
	{
		rdi := api.NewRedeploymentHandler(v0, rdh)

		rdi.RegisterHandlers()
	}

	if err := r.Run(fmt.Sprintf("%s:%d", listenAddress, listenPort)); err != nil {
		log.Fatal(err)
	}
}
