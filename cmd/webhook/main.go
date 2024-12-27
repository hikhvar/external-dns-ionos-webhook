package main

import (
	"fmt"

	"github.com/hikhvar/external-dns-inwx-webhook/cmd/webhook/init/configuration"
	"github.com/hikhvar/external-dns-inwx-webhook/cmd/webhook/init/dnsprovider"
	"github.com/hikhvar/external-dns-inwx-webhook/cmd/webhook/init/logging"
	"github.com/hikhvar/external-dns-inwx-webhook/cmd/webhook/init/server"
	"github.com/hikhvar/external-dns-inwx-webhook/pkg/webhook"
	log "github.com/sirupsen/logrus"
)

const banner = `
 external-dns-inwx-webhook
 version: %s (%s)

`

var (
	Version = "local"
	Gitsha  = "?"
)

func main() {
	fmt.Printf(banner, Version, Gitsha)
	logging.Init()
	config := configuration.Init()
	provider, err := dnsprovider.Init(config)
	if err != nil {
		log.Fatalf("Failed to initialize DNS provider: %v", err)
	}
	srv := server.Init(config, webhook.New(provider))
	server.ShutdownGracefully(srv)
}
