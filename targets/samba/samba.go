package samba

import (
	"bytes"
	"fmt"
	"os"

	"github.com/walnuts1018/backup-manager/config"
	"github.com/walnuts1018/backup-manager/domain"
	"golang.org/x/crypto/ssh"
)

type sambaClient struct {
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
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", config.Config.SambaHostSSHURL, config.Config.SambaHostSSHPort), sambaconfig)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %s", err)
	}
	return &sambaClient{client}, nil
}

func (c *sambaClient) Backup() error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	srcdir := "/mnt/share/"
	dstdir := "/mnt/HDD1TB/smb-backup"
	command := fmt.Sprintf(`
		LATESTBKUPDIR=$(ls %v | tail -n 1)
		DATEDIR=%v/$(date +%%Y-%%m-%%d-%%H-%%M-%%S)
		mkdir $DATEDIR
		rsync -avh --link-dest="%v/$LATESTBKUPDIR" %v "$DATEDIR"
	`, dstdir, dstdir, dstdir, srcdir)

	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(command)
	fmt.Println(b.String())
	if err != nil {
		return fmt.Errorf("failed to run command: %s", err)
	}
	b.Reset()

	return nil
}

func (c *sambaClient) Close() error {
	return c.client.Close()
}
