package main

import (
	"os"

	"github.com/antoniou/zero2Pipe/client"
)

func main() {
	client.New().Run(os.Args)
}
