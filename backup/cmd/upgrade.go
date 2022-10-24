package cmd

import (
	"log"
	"os"

	"github.com/SecurityBrewery/catalystctl/backup"
)

type UpgradeCmd struct {
	AssumeVersion string   `help:"Assume the backup is of this version."`
	SourceFile    *os.File `arg:"" required:"" type:"existingfile" help:"File to upgrade."`
	TargetFile    string   `arg:"" required:"" help:"File to write the upgraded backup to."`
}

func (c *UpgradeCmd) Run() error {
	backupFile, err := backup.NewReader(c.SourceFile)
	if err != nil {
		return err
	}

	log.Printf("Comment / Version: %s", backupFile.Version())

	backupVersion := backupFile.Version()
	if c.AssumeVersion != "" {
		backupVersion = c.AssumeVersion
	}

	newBackupFile, err := backup.NewFromReader(c.TargetFile, backupFile)
	if err != nil {
		return err
	}
	defer newBackupFile.Close()

	switch backupVersion {
	case "v0.10.0":
		if err := upgradeToV0o100(newBackupFile); err != nil {
			return err
		}
	default:
		log.Printf("Unknown backup version: %s", backupVersion)
	}

	return nil
}

func upgradeToV0o100(w *backup.Writer) error {
	log.Println("Upgrading backup to v0.10.0")

	if err := addMinioFolder(w); err != nil {
		return err
	}

	return w.SetVersion("v0.10.0")
}

func addMinioFolder(w *backup.Writer) error {
	log.Printf("Adding minio folder to backup")

	return w.CreateFolder("minio")
}
