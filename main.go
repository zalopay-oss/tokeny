package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli/v2"
	"github.com/zalopay-oss/tokeny/pkg/keyvalue"
	"github.com/zalopay-oss/tokeny/pkg/password"
	"github.com/zalopay-oss/tokeny/pkg/session"
	"github.com/zalopay-oss/tokeny/pkg/tokeny"
	"github.com/zalopay-oss/tokeny/pkg/tokenycli"
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
