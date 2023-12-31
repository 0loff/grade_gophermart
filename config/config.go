package config

import (
	"flag"
	"os"
)

type Config struct {
	Host                string
	SigningKey          string
	LogLevel            string
	DatabaseDSN         string
	AccrualSystemAddres string
}

type ConfigBuilder struct {
	config Config
}

func (cb ConfigBuilder) SetHost(host string) ConfigBuilder {
	cb.config.Host = host
	return cb
}

func (cb ConfigBuilder) SetSigningKey(signingKey string) ConfigBuilder {
	cb.config.SigningKey = signingKey
	return cb
}

func (cb ConfigBuilder) SetLogLevel(logLevel string) ConfigBuilder {
	cb.config.LogLevel = logLevel
	return cb
}

func (cb ConfigBuilder) SetDatabaseDSN(databaseDSN string) ConfigBuilder {
	cb.config.DatabaseDSN = databaseDSN
	return cb
}

func (cb ConfigBuilder) SetAccrualSystemAddres(accrualSystemAddres string) ConfigBuilder {
	cb.config.AccrualSystemAddres = accrualSystemAddres
	return cb
}

func (cb ConfigBuilder) Build() Config {
	return cb.config
}

func NewConfigBuilder() Config {
	var host string
	flag.StringVar(&host, "a", "localhost:8080", "server host")

	var signingKey string
	flag.StringVar(&signingKey, "sk", "signingkey", "signing key")

	var logLevel string
	flag.StringVar(&logLevel, "l", "info", "log level")

	var databaseDSN string
	flag.StringVar(&databaseDSN, "d", "", "Database DSN config string")
	// flag.StringVar(&databaseDSN, "d", "host=localhost port=5432 user=postgres password=root dbname=gophermart sslmode=disable", "Database DSN config string")

	var accrualSystemAddres string
	flag.StringVar(&accrualSystemAddres, "r", "", "accrual system addres")
	// flag.StringVar(&accrualSystemAddres, "r", "localhost:8081", "accrual system addres")

	flag.Parse()

	if envHost := os.Getenv("RUN_ADDRESS"); envHost != "" {
		host = envHost
	}

	if envSigningKey := os.Getenv("SIGNING_KEY"); envSigningKey != "" {
		signingKey = envSigningKey
	}

	if envLoglevel := os.Getenv("LOG_LEVEL"); envLoglevel != "" {
		logLevel = envLoglevel
	}

	if envStorageFile := os.Getenv("DATABASE_URI"); envStorageFile != "" {
		databaseDSN = envStorageFile
	}

	if envAccrualSystemAddres := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddres != "" {
		accrualSystemAddres = envAccrualSystemAddres
	}

	return new(ConfigBuilder).
		SetHost(host).
		SetSigningKey(signingKey).
		SetLogLevel(logLevel).
		SetDatabaseDSN(databaseDSN).
		SetAccrualSystemAddres(accrualSystemAddres).
		Build()
}
