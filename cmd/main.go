package main

import (
	"fmt"
	"github.com/ltpquang/tokeny/pkg/keyvalue"
	"github.com/ltpquang/tokeny/pkg/password"
	"github.com/ltpquang/tokeny/pkg/session"
	"github.com/ltpquang/tokeny/pkg/tokeny"
	"github.com/ltpquang/tokeny/pkg/tokenycli"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dbDir := fmt.Sprintf("%s/.tokeny", usr.HomeDir)
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	kvStore, err := keyvalue.NewSQLStore(fmt.Sprintf("%s/d.db", dbDir))
	if err != nil {
		log.Fatal(err)
	}

	pwdManager := password.NewManager(kvStore)

	sessionManager := session.NewManager(kvStore)

	tokenRepo := tokeny.NewRepository(kvStore)

	cliSvc := tokenycli.NewService(pwdManager, sessionManager, tokenRepo)

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = "Another TOTP generator"

	cliSvc.Register(app)

	if err = app.Run(os.Args); err != nil {
		log.SetFlags(0)
		log.Println(err.Error())
	}
}
