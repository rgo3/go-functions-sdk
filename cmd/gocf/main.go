package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dergoegge/go-functions-sdk/internal/pkg/cli"
	"github.com/google/logger"
)

func main() {
	logger.Init("gocf logger", true, false, ioutil.Discard)
	logger.SetFlags(log.Ltime)

	app := cli.Init()
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
