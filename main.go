package main

import (
	"context"
	"log"

	"github.com/RichSweetWafer/EWallet-API/api"
	"github.com/RichSweetWafer/EWallet-API/config"
	"github.com/RichSweetWafer/EWallet-API/wallets"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	wallets := wallets.NewMongoWallets(cfg.Database)
	server := api.NewServer(cfg.HTTPServer, wallets)
	server.Start(ctx)
}
