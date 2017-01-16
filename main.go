package main

import (
	"os"

	"github.com/antoniou/zero2Pipe/client"
)

//go:generate go-bindata templates/...

func main() {
	client.New().Run(os.Args)
}
