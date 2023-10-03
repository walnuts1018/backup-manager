package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/walnuts1018/backup-manager/config"
	"github.com/walnuts1018/backup-manager/domain"
	"github.com/walnuts1018/backup-manager/exporter"
	"github.com/walnuts1018/backup-manager/targets/samba"
	"github.com/walnuts1018/backup-manager/timeJST"
	"golang.org/x/exp/slog"
)

type backupTask struct {
	name         string
	backupClient domain.BackupClient
	every        int
	at           string
}

func main() {
	config.LoadConfig()

	go exporter.Export()

	s := gocron.NewScheduler(timeJST.JST)

	smbclient, err := samba.NewClient()
	if err != nil {
		slog.Error("failed to create samba client", "error", err)
		os.Exit(1)
	}
	defer smbclient.Close()

	tasks := []backupTask{
		{
			name:         "samba backup",
			backupClient: smbclient,
			every:        config.Config.SambaTaskIntervalDays,
			at:           config.Config.SambaTaskDoAt,
		},
	}

	for _, task := range tasks {
		job, err := s.Every(task.every).Days().At(task.at).WaitForSchedule().Do(func() {
			doBackup(task.backupClient)
		})
		if err != nil {
			slog.Error("failed to set scheduler", "error", err)
			continue
		}
		job.RegisterEventListeners(
			gocron.WhenJobReturnsError(func(jobname string, err error) {
				exporter.SetJobs(jobname, false)
			}),
			gocron.WhenJobReturnsNoError(func(jobname string) {
				exporter.SetJobs(jobname, true)
			}),
		)
		slog.Info("set scheduler", "job", task.name, "every", fmt.Sprintf("%v days", task.every), "at", task.at)
	}

	go func() {
		for {
			jobs := s.Jobs()
			actives := 0
			for _, job := range jobs {
				if job.IsRunning() {
					actives++
				}
			}
			exporter.SetRunningTasks(actives)
			time.Sleep(1 * time.Minute)
		}
	}()

	slog.Info("start scheduler")
	s.StartBlocking()
}

func doBackup(backupClient domain.BackupClient) {
	now := timeJST.Now()
	logfilepath := config.Config.LogfileBase + "-" + now.Format("2006-01-02_15h04m05s") + ".log"
	if _, err := os.Stat(logfilepath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(logfilepath), 0700)
	}
	logfile, err := os.OpenFile(logfilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("failed to open logfile", "error", err)
		os.Exit(1)
	}
	defer logfile.Close()

	logWriter := io.MultiWriter(os.Stdout, logfile)

	err = backupClient.Backup(logWriter)
	if err != nil {
		slog.Error("failed to backup", "error", err)
		os.Exit(1)
	}
}
