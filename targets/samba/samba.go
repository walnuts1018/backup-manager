package samba

import (
	"fmt"
	"io"
	"os"

	"github.com/walnuts1018/backup-manager/config"
	"github.com/walnuts1018/backup-manager/domain"
	"github.com/walnuts1018/backup-manager/timeJST"
	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slog"
)

type sambaClient struct {
	config *ssh.ClientConfig
	client *ssh.Client
}

func NewClient() (domain.BackupClient, error) {
	_, _, hostKey, _, _, _ := ssh.ParseKnownHosts([]byte(config.Config.SambaHostSSHPublicKey))
	key, err := os.ReadFile(config.Config.SambaHostSSHKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %s", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key failed: %s", err)
	}

	sambaconfig := &ssh.ClientConfig{
		User: config.Config.SambaHostSSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return &sambaClient{
		config: sambaconfig,
	}, nil
}

func (c *sambaClient) Backup(logWriter io.Writer) error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", config.Config.SambaHostSSHURL, config.Config.SambaHostSSHPort), c.config)
	if err != nil {
		return fmt.Errorf("dial failed: %s", err)
	}
	c.client = client

	slog.Info("start backup")
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	srcdir := config.Config.SambaSrcDir
	dstdir := config.Config.SambaDstDir
	datedir := dstdir + "/" + timeJST.Now().Format("backup-2006-01-02_15h04m05s")
	command := fmt.Sprintf(`
		LATESTBKUPDIR=$(ls %v | tail -n 1)
		mkdir %v
		rsync -avh --link-dest="%v/$LATESTBKUPDIR" %v "%v"
	`, dstdir, datedir, dstdir, srcdir, datedir)

	session.Stdout = logWriter
	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("failed to run command: %s", err)
	}

	return nil
}

func (c *sambaClient) Close() error {
	return c.client.Close()
}
