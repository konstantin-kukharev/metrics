package settings

import (
	"flag"
	"os"
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Config struct {
	Address         string
	StoreInterval   int    // интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск
	FileStoragePath string // путь до файла, куда сохраняются текущие значения
	Restore         bool   // загружать или нет ранее сохранённые значения из указанного файла при старте
	DatabaseDNS     string // адрес базы данных, если не указана, то используется по умолчанию
}

func NewConfig() *Config {
	c := &Config{
		Address:         internal.DefaultServerAddr,
		StoreInterval:   internal.DefaultServerStoreInterval,
		FileStoragePath: internal.DefaultFileStoragePath,
		Restore:         internal.DefaultRestore,
	}

	return c
}

func (c *Config) WithFlag() {
	flag.StringVar(&c.Address, "a", internal.DefaultServerAddr, "server address")
	flag.IntVar(&c.StoreInterval, "i", internal.DefaultServerStoreInterval, "interval to store data on FS")
	flag.StringVar(&c.FileStoragePath, "f", internal.DefaultFileStoragePath, "file path to store data")
	flag.BoolVar(&c.Restore, "r", internal.DefaultRestore, "is restore file data")
	flag.StringVar(&c.DatabaseDNS, "d", "", "database dns path")
	flag.Parse()
}

func (c *Config) WithEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.Address = envRunAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		if val, err := strconv.Atoi(envStoreInterval); err != nil {
			c.StoreInterval = val
		}
	}
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		c.FileStoragePath = envFilePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if val, err := strconv.ParseBool(envRestore); err != nil {
			c.Restore = val
		}
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		c.DatabaseDNS = envDB
	}
}
