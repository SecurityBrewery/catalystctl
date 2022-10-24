package main

import (
	"fmt"
	"github.com/SecurityBrewery/catalystctl/backup/cmd"
	"github.com/SecurityBrewery/catalystctl/client"
	"github.com/SecurityBrewery/catalystctl/config"
	"github.com/SecurityBrewery/catalystctl/fake"
	"github.com/alecthomas/kong"
	"net/url"
)

type CLI struct {
	SetCredentials SetCredentialsCmd `cmd:"" help:"Set the credentials to use for authentication."`

	Backup struct {
		Create   cmd.CreateCmd   `kong:"cmd,help='Create a backup.'"`
		Validate cmd.ValidateCmd `kong:"cmd,help='Validate a backup.'"`
		Upgrade  cmd.UpgradeCmd  `kong:"cmd,help='Upgrade a backup.'"`
		Restore  cmd.RestoreCmd  `kong:"cmd,help='Restore a backup.'"`
	} `cmd:"" help:"Backup commands."`

	Generate struct {
		FakeData fake.Cmd `cmd:"" help:"Generate fake data."`
	} `cmd:"" help:"Generate commands."`

	Ping PingCmd `cmd:"" help:"Ping the server."`
}

func main() {
	cli := &CLI{}
	ctx := kong.Parse(cli)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

type SetCredentialsCmd struct {
	URL   *url.URL `arg:"" required:"" help:"URL of the Catalyst server."`
	Token string   `arg:"" required:"" help:"Token to use for authentication."`
}

func (c *SetCredentialsCmd) Run() error {
	return config.CreateFile(c.URL.String(), c.Token)
}

type PingCmd struct {
	URL   *url.URL `help:"URL of the Catalyst server."`
	Token string   `help:"Token to access the Catalyst server."`
}

func (c *PingCmd) Run() error {
	catalystClient, err := client.New(c.URL, c.Token)
	if err != nil {
		return err
	}

	_, err = catalystClient.Version()
	if err != nil {
		return err
	}
	fmt.Println("Ping successful.")
	return nil
}
