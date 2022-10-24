package cmd

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/SecurityBrewery/catalystctl/backup"
	"github.com/SecurityBrewery/catalystctl/client"
)

type RestoreCmd struct {
	SourceFile *os.File `arg:"" required:"" type:"existingfile" help:"File to restore."`
	URL        *url.URL `help:"URL of the Catalyst server."`
	Token      string   `help:"Token to access the Catalyst server."`
	Force      bool     `help:"Force restore even if the version is incompatible."`
}

func (c *RestoreCmd) Run() error {
	catalystClient, err := client.New(c.URL, c.Token)
	if err != nil {
		return err
	}

	catalystVersion, err := catalystClient.Version()
	if err != nil {
		return err
	}

	backupFile, err := backup.NewReader(c.SourceFile)
	if err != nil {
		return err
	}

	if backupFile.Version() != catalystVersion && !c.Force {
		log.Printf("Backup version %s does not match Catalyst version %s. Use --force to override.", backupFile.Version(), catalystVersion)
		return nil
	}

	return c.restoreBackup(catalystClient, backupFile)
}

func (c *RestoreCmd) restoreBackup(client *client.CatalystClient, backupFile *backup.Reader) error {
	body, contentType, err := backupFile.AsPayload()
	if err != nil {
		return err
	}

	// upload file
	header := http.Header{"Content-Type": []string{contentType}}
	resp, err := client.Post("/api/backup/restore", header, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
