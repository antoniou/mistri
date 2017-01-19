package main

import (
	"os"

	"github.com/antoniou/zero2Pipe/client"
)

//go:generate go-bindata -o domain/template.go -pkg domain templates/...

func main() {
	client.New().Run(os.Args)
}
