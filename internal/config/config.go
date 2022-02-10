package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

//goland:noinspection SqlResolve
var Conf struct {
	TgToken string `env:"TG_TOKEN,notEmpty"  envExpand:"true"`

	PgUser          string        `env:"DB_USER,notEmpty"  envExpand:"true"`
	PgPass          string        `env:"DB_PASSWORD,notEmpty"  envExpand:"true"`
	PgHost          string        `env:"DB_ADDR,notEmpty"  envExpand:"true"`
	PgPort          int           `env:"DB_PORT,notEmpty"  envExpand:"true"`
	PgDb            string        `env:"DB_DATABASE,notEmpty"  envExpand:"true"`
	PgParams        string        `env:"DB_PARAMS,notEmpty" envDefault:"sslmode=disable&application_name=gotestbot"  envExpand:"true"`
	PgMaxOpenConn   int           `env:"DB_MAX_OPEN_CONN" envDefault:"10"`
	PgMaxIdleConn   int           `env:"DB_MAX_IDLE_CONN" envDefault:"0"`
	PgMaxLifeTime   time.Duration `env:"DB_MAX_LIFE_TIME" envDefault:"30m"`
	PgMaxIdleTime   time.Duration `env:"DB_MAX_IDLE_TIME" envDefault:"1m"`
	PoolConnTimeout time.Duration `env:"POOL_CONNECTION_TIMEOUT" envDefault:"1m"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"logstash"`
	Dry       bool   `env:"DRY" envDefault:"false"`
}

type SftpConfig struct {
	ssh.ClientConfig
	Addr         string
	DirPath      string
	SftpTmpExt   string
	SftpFinalExt string
}

type S3Config struct {
	S3Url    string
	S3Id     string
	S3Secret string
	S3Region string
	S3Bucket string
	S3Prefix string
}

func InitConfig() {
	conf := &Conf
	if err := env.Parse(conf); err != nil {
		log.Fatal().Err(err).Msg("Unable to init config")
	}
}

func GetPgDsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		url.QueryEscape(Conf.PgUser),
		url.QueryEscape(Conf.PgPass),
		Conf.PgHost,
		Conf.PgPort,
		Conf.PgDb,
		Conf.PgParams)
}

func DsnMaskPass(dsn string) string {
	at := strings.Index(dsn, "@")
	beforeAt := dsn[:at]
	colon := strings.LastIndex(beforeAt, ":")
	beforeColon := dsn[:colon+1]
	afterAt := dsn[at:]
	return beforeColon + "********" + afterAt
}
