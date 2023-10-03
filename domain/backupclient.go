package domain

import "io"

type BackupClient interface {
	Backup(logWriter io.Writer) error
	Close() error
}
