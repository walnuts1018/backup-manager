package domain

type BackupClient interface {
	Backup() error
	Close() error
}
