package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"

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
	SambaTaskIntervalDays int    `env:"SAMBA_TASK_INTERVAL_DAYS"`
	SambaTaskDoAt         string `env:"SAMBA_TASK_DO_AT"`
	LogfileBase           string `env:"LOG_FILE"`
	ServerPort            string
}

var Config = Config_t{
	SambaHostSSHPort: "22",
	LogfileBase:      "/var/log/backup-manager",
}

func LoadConfig() error {
	serverport := flag.String("port", "8080", "server port")
	flag.Parse()
	Config.ServerPort = *serverport

	err := godotenv.Load(".env")
	if err != nil {
		slog.Warn("Error loading .env file")
	}

	t := reflect.TypeOf(Config)
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldType := t.Field(i).Type
		tag, ok := t.Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}
		v, ok := os.LookupEnv(tag)
		if !ok {
			return fmt.Errorf("%s is not set", tag)
		}
		switch fieldType.Kind() {
		case reflect.String:
			reflect.ValueOf(&Config).Elem().FieldByName(fieldName).SetString(v)
		case reflect.Int:
			i, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("%s is not int", tag)
			}
			reflect.ValueOf(&Config).Elem().FieldByName(fieldName).SetInt(int64(i))
		}
	}
	return nil
}
