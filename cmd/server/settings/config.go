package settings

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Config struct {
	Address         string
	StoreInterval   int    // интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск
	FileStoragePath string // путь до файла, куда сохраняются текущие значения
	Restore         bool   // загружать или нет ранее сохранённые значения из указанного файла при старте
	DatabaseDNS     string // адрес базы данных, если не указана, то используется по умолчанию
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetStoreInterval() time.Duration {
	return time.Duration(c.StoreInterval * int(time.Second))
}

func (c *Config) GetFileStoragePath() string {
	return c.FileStoragePath
}

func (c *Config) GetRestore() bool {
	return c.Restore
}

func (c *Config) GetDatabaseDNS() string {
	return c.DatabaseDNS
}

func NewConfig() *Config {
	c := &Config{
		Address:         internal.DefaultServerAddr,
		StoreInterval:   internal.DefaultServerStoreInterval,
		FileStoragePath: internal.DefaultFileStoragePath,
		Restore:         internal.DefaultRestore,
		DatabaseDNS:     internal.DefaultDatabaseDNS,
	}

	return c
}

func (c *Config) WithFlag() *Config {
	flag.StringVar(&c.Address, "a", internal.DefaultServerAddr, "server address")
	flag.IntVar(&c.StoreInterval, "i", internal.DefaultServerStoreInterval, "interval to store data on FS")
	flag.StringVar(&c.FileStoragePath, "f", internal.DefaultFileStoragePath, "file path to store data")
	flag.BoolVar(&c.Restore, "r", internal.DefaultRestore, "file path to store data")
	flag.StringVar(&c.DatabaseDNS, "d", internal.DefaultDatabaseDNS, "database dns path")
	flag.Parse()

	return c
}

func (c *Config) WithEnv() *Config {
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

	return c
}
