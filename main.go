package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/syndtr/goleveldb/leveldb"
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
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/ld.db", dbDir), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	kvStore := keyvalue.NewLevelDBStore(db)

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
