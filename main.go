package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/pheelee/dyndns/pkg/config"
	"github.com/pheelee/dyndns/pkg/logger"
	"github.com/pheelee/dyndns/pkg/server"
)

type arguments struct {
	cfgPath string
}

func (a *arguments) valid() bool {
	return a.cfgPath != ""
}

func main() {
	args := arguments{}

	flag.StringVar(&args.cfgPath, "config", "", "Path to config file")
	flag.Parse()
	if !args.valid() {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	cfg := config.Config{}
	logger.Info("Reading config file " + args.cfgPath)
	cfg.Load(args.cfgPath)

	if cfg.HTTPReqAuth.Username == "" || cfg.HTTPReqAuth.Password == "" {
		logger.Error("Username and/or password not specified please set env vars HTTPREQ_USERNAME and HTTPREQ_PASSWORD or specify them in the config file")
		return
	}

	// verify we have an active dns server
	if cfg.DNSServer.Active == nil {
		logger.Error("No DNS Server specified in config file")
		return
	}
	logger.Info(fmt.Sprintf("Setup API on Port %d", cfg.Global.Listen))
	server.SetupServer(&cfg)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Global.Listen), nil)

}
