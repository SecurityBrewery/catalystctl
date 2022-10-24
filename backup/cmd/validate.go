package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/SecurityBrewery/catalystctl/backup"
)

type ValidateCmd struct {
	SourceFile *os.File `arg:"" required:"" help:"File to validate."`
}

func (c *ValidateCmd) Run() error {
	backupFile, err := backup.NewReader(c.SourceFile)
	if err != nil {
		return err
	}

	log.Printf("Comment / Version: %s", backupFile.Version())
	switch backupFile.Version() {
	case "v0.10.0":
		checkArangoCollections(backupFile)
	default:
		log.Printf("Unknown backup version: %s", backupFile.Version())
	}

	return nil
}

func checkArangoCollections(zr *backup.Reader) {
	requiredCollections := map[string]bool{
		"automations": false,
		"dashboards":  false,
		"jobs":        false,
		"logs":        false,
		"migrations":  false,
		"playbooks":   false,
		"related":     false,
		"settings":    false,
		"templates":   false,
		"tickets":     false,
		"tickettypes": false,
		"userdata":    false,
		"users":       false,
	}

	for _, zf := range zr.Files() {
		for k := range requiredCollections {
			if strings.HasPrefix(zf.Name, "arango/"+k) {
				requiredCollections[k] = true
			}
		}
	}

	// check if all required files are present
	for k, v := range requiredCollections {
		if !v {
			log.Printf("Missing collection: %s", k)
		}
	}
}
