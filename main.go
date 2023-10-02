package main

import (
	"os"

	"github.com/walnuts1018/backup-manager/config"
	"github.com/walnuts1018/backup-manager/targets/samba"
	"golang.org/x/exp/slog"
)

func main() {
	config.LoadConfig()
	smbclient, err := samba.NewClient()
	if err != nil {
		slog.Error("failed to create samba client", "error", err)
		os.Exit(1)
	}
	defer smbclient.Close()

	err = smbclient.Backup()
	if err != nil {
		slog.Error("failed to backup", "error", err)
		os.Exit(1)
	}
}
