package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/walnuts1018/backup-manager/config"
	"github.com/walnuts1018/backup-manager/targets/samba"
	"golang.org/x/exp/slog"
)

func main() {
	config.LoadConfig()

	if _, err := os.Stat(config.Config.Logfile); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(config.Config.Logfile), 0700)
	}
	logfile, err := os.OpenFile(config.Config.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("failed to open logfile", "error", err)
		os.Exit(1)
	}
	defer logfile.Close()

	logWriter := io.MultiWriter(os.Stdout, logfile)

	smbclient, err := samba.NewClient()
	if err != nil {
		slog.Error("failed to create samba client", "error", err)
		os.Exit(1)
	}
	defer smbclient.Close()

	err = smbclient.Backup(logWriter)
	if err != nil {
		slog.Error("failed to backup", "error", err)
		os.Exit(1)
	}
}
