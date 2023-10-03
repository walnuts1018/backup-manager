package config

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config_t struct {
	SambaHostSSHURL       string `env:"SAMBA_HOST_SSH_URL"`
	SambaHostSSHPort      string `env:"SAMBA_HOST_SSH_PORT"`
	SambaHostSSHUser      string `env:"SAMBA_HOST_SSH_USER"`
	SambaHostSSHKeyPath   string `env:"SAMBA_HOST_SSH_KEY_PATH"`
	SambaHostSSHPublicKey string `env:"SAMBA_HOST_SSH_PUBLIC_KEY"`
	SambaSrcDir           string `env:"SAMBA_SRC_DIR"`
	SambaDstDir           string `env:"SAMBA_DST_DIR"`
	Logfile               string `env:"LOG_FILE"`
}

var Config = Config_t{
	SambaHostSSHPort: "22",
	Logfile:          "/var/log/backup-manager.log",
}

func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	t := reflect.TypeOf(Config)
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		tag, ok := t.Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}
		v, ok := os.LookupEnv(tag)
		if !ok {
			return fmt.Errorf("%s is not set", tag)
		}
		reflect.ValueOf(&Config).Elem().FieldByName(fieldName).SetString(v)
	}
	return nil
}
