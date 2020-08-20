package tokenycli

import "github.com/urfave/cli/v2"

type Service interface {
	Register(app *cli.App)
}
