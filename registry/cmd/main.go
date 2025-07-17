package main

import (
	_ "embed"
	"flag"

	"github.com/burgerdev/evil-registry/registry"
)

var addr = flag.String("addr", ":80", "address to serve on")

func main() {
	flag.Parse()
	registry.Run(addr)
}
