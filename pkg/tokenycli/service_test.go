package tokenycli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"github.com/zalopay-oss/tokeny/pkg/totp"
)

func TestWriteToken(t *testing.T) {
	token := totp.Token{Value: "123456", TimeoutSec: 12}

	t.Run("raw", func(t *testing.T) {
		var output bytes.Buffer

		writeToken(&output, token, "example", true)

		assert.Equal(t, "123456", output.String())
	})

	t.Run("normal", func(t *testing.T) {
		var output bytes.Buffer

		writeToken(&output, token, "example", false)

		assert.Equal(t, "Here is your token for 'example', valid within the next 12 seconds\n123456\n", output.String())
	})
}

func TestGetCommandRawFlag(t *testing.T) {
	svc := &service{}
	commands := svc.getNormalCommands()

	var getCommand *cli.Command
	for _, command := range commands {
		if command.Name == "get" {
			getCommand = command
			break
		}
	}

	if assert.NotNil(t, getCommand) {
		assert.Contains(t, getCommand.Names(), "get")
		assert.Contains(t, getCommand.Flags[1].Names(), "raw")
		assert.Contains(t, getCommand.Flags[1].Names(), "r")
	}
}
