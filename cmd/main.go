package main

import (
	"github.com/ltpquang/tokeny/pkg/tokenycli"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cliSvc := tokenycli.NewService()

	app := cli.NewApp()
	app.EnableBashCompletion = true

	cliSvc.Register(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
