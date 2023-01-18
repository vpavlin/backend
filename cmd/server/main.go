package main

import (
	"log"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/w3"
	"github.com/metaconflux/backend/internal/api/v1alpha"
	cache "github.com/metaconflux/backend/internal/cache/ipfs"
	"github.com/metaconflux/backend/internal/resolver/memory"

	"github.com/metaconflux/backend/internal/transformers"
	"github.com/metaconflux/backend/internal/transformers/core/v1alpha/contract"
	"github.com/metaconflux/backend/internal/transformers/core/v1alpha/ipfs"
)

func main() {
	e := echo.New()
	e.Use(
		middleware.Logger(), // Log everything to stdout
	)
	g := e.Group("/api/v1alpha")

	tm, _ := transformers.NewTransformerManager()

	url := "http://localhost:5001"
	shell := shell.NewShellWithClient(url, &http.Client{})
	ipfsT := ipfs.NewTransformer(shell)

	clients := make(map[uint64]*w3.Client)

	var err error
	clients[80001] = w3.MustDial("https://polygon-testnet.public.blastapi.io")

	defer clients[80001].Close()

	constractT := contract.NewTransformer(clients)

	err = tm.Register(ipfs.GVK, ipfsT.WithSpec)
	if err != nil {
		log.Fatal(err)
	}

	err = tm.Register(contract.GVK, constractT.WithSpec)
	if err != nil {
		log.Fatal(err)
	}

	r := memory.NewResolver()
	c := cache.NewIPFSCache(url, shell)
	a := v1alpha.NewAPI(c, r, tm)
	a.Register(g)

	log.Fatal(e.Start("localhost:8081"))
}
