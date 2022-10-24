package cmd

import (
	"io"
	"net/url"
	"os"

	"github.com/SecurityBrewery/catalystctl/client"
)

type CreateCmd struct {
	URL        *url.URL `help:"URL of the Catalyst server."`
	Token      string   `help:"Token to access the Catalyst server."`
	TargetFile string   `arg:"" required:"" help:"File to write the backup to."`
}

func (c *CreateCmd) Run() error {
	catalystClient, err := client.New(c.URL, c.Token)
	if err != nil {
		return err
	}

	resp, err := catalystClient.Get("/api/backup/create")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// write file
	f, err := os.Create(c.TargetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	return err
}
