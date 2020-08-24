package main

import (
	"github.com/ltpquang/tokeny/pkg/keyvalue"
	"github.com/ltpquang/tokeny/pkg/password"
	"github.com/ltpquang/tokeny/pkg/tokeny"
	"github.com/ltpquang/tokeny/pkg/tokenycli"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	kvStore, err := keyvalue.NewSQLStore("./tokeny.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	pwdManager := password.NewManager(kvStore)

	tokenRepo := tokeny.NewRepository(kvStore)

	cliSvc := tokenycli.NewService(pwdManager, tokenRepo)

	app := cli.NewApp()
	app.EnableBashCompletion = true

	cliSvc.Register(app)

	if err = app.Run(os.Args); err != nil {
		log.SetFlags(0)
		log.Println(err.Error())
	}
}
