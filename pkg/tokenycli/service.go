package tokenycli

import (
	"fmt"
	"io"
	"os"

	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/zalopay-oss/tokeny/pkg/password"
	"github.com/zalopay-oss/tokeny/pkg/session"
	"github.com/zalopay-oss/tokeny/pkg/tokeny"
	"github.com/zalopay-oss/tokeny/pkg/totp"
)

var (
	ppidStr = fmt.Sprintf("%d", os.Getppid())
)

type service struct {
	pwdManager     password.Manager
	sessionManager session.Manager
	tokenRepo      tokeny.Repository
}

func NewService(pwdManager password.Manager, sessionManager session.Manager, tokenRepo tokeny.Repository) *service {
	return &service{
		pwdManager:     pwdManager,
		sessionManager: sessionManager,
		tokenRepo:      tokenRepo,
	}
}

func (s *service) Register(app *cli.App) error {
	userProfileAvailable, err := s.pwdManager.IsRegistered()
	if err != nil {
		return err
	}

	if !userProfileAvailable {
		app.Commands = s.getSetupCommand()
	} else {
		app.Commands = s.getNormalCommands()
	}

	return nil
}

func (s *service) getSetupCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "setup",
			Usage:  "setup master password",
			Action: s.setup,
		},
	}
}

func (s *service) getNormalCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "add",
			Usage: "add new entry",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "alias",
					Aliases:  []string{"a"},
					Required: true,
					Usage:    "entry name/alias, must be identical to each other",
				},
				&cli.StringFlag{
					Name:     "secret",
					Aliases:  []string{"s"},
					Required: true,
					Usage:    "secret of the entry",
				},
			},
			Action: s.sessionWrapper(s.add),
		},
		{
			Name:  "get",
			Usage: "get OTP",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:     "copy",
					Aliases:  []string{"c"},
					Required: false,
					Usage:    "copy generated token to clipboard",
				},
				&cli.BoolFlag{
					Name:     "raw",
					Aliases:  []string{"r"},
					Required: false,
					Usage:    "print only the generated token",
				},
			},
			Action: s.sessionWrapper(s.get),
		},
		{
			Name:   "delete",
			Usage:  "delete selected entry",
			Action: s.sessionWrapper(s.delete),
		},
		{
			Name:   "list",
			Usage:  "list all entries",
			Action: s.sessionWrapper(s.list),
		},
	}
}

func (s *service) setup(c *cli.Context) error {
	registered, err := s.pwdManager.IsRegistered()
	if err != nil {
		return err
	}
	if registered {
		println("You have registered already.")
		return nil
	}
	return s.doRegister()
}

func (s *service) doRegister() error {
	prompt := promptui.Prompt{
		Label: "Password",
		Mask:  ' ',
	}

	pwd, err := prompt.Run()

	if err != nil {
		return err
	}

	prompt = promptui.Prompt{
		Label: "Re-type password",
		Mask:  ' ',
	}

	rePwd, err := prompt.Run()

	if err != nil {
		return err
	}

	err = s.pwdManager.Register(pwd, rePwd)
	if err != nil {
		if errors.Is(err, password.ErrPasswordsMismatch) {
			println("Passwords do not match, please try again.")
			return nil
		}
		return err
	}

	println("Registered.")
	return nil
}

func (s *service) sessionWrapper(actionFunc cli.ActionFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		if valid, err := s.ensureSession(); err != nil || !valid {
			return err
		}
		return actionFunc(c)
	}
}

func (s *service) add(c *cli.Context) error {
	alias := c.String("alias")
	secret := c.String("secret")
	err := s.tokenRepo.Add(alias, secret)
	if err != nil {
		if errors.Is(err, tokeny.ErrEntryExistedBefore) {
			println("Alias has been used before, please choose another.")
			return nil
		}
		return err
	}
	println("Entry has been added successfully.")
	return nil
}

func (s *service) get(c *cli.Context) error {
	var alias string
	if c.NArg() > 0 {
		alias = c.Args().Get(0)
	} else {
		var err error
		alias, err = s.tokenRepo.LastValidEntry()
		if errors.Is(err, tokeny.ErrNoEntryFound) {
			println("Please specify entry to generate token: tokeny get <alias>")
			return nil
		}
		if err != nil {
			return err
		}
	}
	t, err := s.tokenRepo.Generate(alias)
	if err != nil {
		if errors.Is(err, tokeny.ErrNoEntryFound) {
			println("Invalid entry, please choose another.")
			return nil
		}
		return err
	}
	if err := writeToken(os.Stdout, t, alias, c.Bool("raw")); err != nil {
		return err
	}
	if c.Bool("copy") {
		err := clipboard.WriteAll(t.Value)
		if err != nil {
			if c.Bool("raw") {
				return err
			}
			println("Cannot copy to clipboard.")
		} else if !c.Bool("raw") {
			println("Copied to clipboard.")
		}
	}
	return nil
}

func writeToken(writer io.Writer, token totp.Token, alias string, raw bool) error {
	if raw {
		_, err := fmt.Fprint(writer, token.Value)
		return err
	}

	secString := "second"
	if token.TimeoutSec > 1 {
		secString += "s"
	}
	_, err := fmt.Fprintf(
		writer,
		"Here is your token for '%s', valid within the next %d %s\n%s\n",
		alias,
		token.TimeoutSec,
		secString,
		token.Value,
	)
	return err
}

func (s *service) delete(c *cli.Context) error {
	if c.NArg() == 0 {
		println("Please specify entry to be deleted.")
		return nil
	}
	alias := c.Args().Get(0)

	err := s.tokenRepo.Delete(alias)
	if err != nil {
		if errors.Is(err, tokeny.ErrNoEntryFound) {
			println("Invalid entry, please choose another.")
			return nil
		}
		return err
	}
	println("Deleted.")
	return nil
}

func (s *service) list(c *cli.Context) error {
	aliases, err := s.tokenRepo.List()
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		println("No entry.")
		return nil
	}

	if len(aliases) == 1 {
		println("Here is your entry:")
	} else {
		println("Here are your entries:")
	}

	for _, alias := range aliases {
		println(alias)
	}

	return nil
}

func (s *service) ensureSession() (bool, error) {
	valid, err := s.sessionManager.IsSessionValid(ppidStr)
	if err != nil {
		return false, err
	}

	if valid {
		return true, nil
	}

	err = s.doLogin()
	if err != nil {
		if errors.Is(err, password.ErrWrongPassword) {
			println("Wrong password, please try again.")
			return false, nil
		}
		return false, err
	}

	err = s.sessionManager.NewSession(ppidStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) doLogin() error {
	prompt := promptui.Prompt{
		Label: "Password",
		Mask:  ' ',
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	err = s.pwdManager.Login(result)
	if err != nil {
		return err
	}
	return nil
}
