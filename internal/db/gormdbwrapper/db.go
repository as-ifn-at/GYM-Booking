package gormdbwrapper

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/as-ifn-at/REST/config"
	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// sql impl here to be done
const (
	maxIdleConnection     = 250
	maxOpenConnection     = 500
	maxConnectionLifeTime = 5 * time.Second

	mysqlDB    = "mysql"
	postgresDB = "postgres"
)

func MariaDNS(username, password, host, dbname string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true&tls=false&interpolateParams=true",
		username, password, host, dbname,
	)
}

func PostgresDSN(username, password, host, port, dbname string, sslmode string) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   dbname,
	}

	// Add query params
	q := u.Query()
	if sslmode == "" {
		sslmode = "disable" // default for local dev
	}
	q.Set("sslmode", sslmode)
	u.RawQuery = q.Encode()

	return u.String()
}

type DBWrapper struct {
	log         zerolog.Logger
	config      config.DBConfigOptions
	db          *gorm.DB
	ErrorCode   config.ErrorCode
	cacheConfig config.CacheConfig
	// retry
}

func NewDBWrapper(log zerolog.Logger, config config.DBConfigOptions,
	cacheConfig config.CacheConfig) (*DBWrapper, error) {

	if config.DBName == "" || config.Host == "" {
		return nil, fmt.Errorf("invalid config for host: %v", config.Host)
	}

	DBWrapper := &DBWrapper{
		log:         log,
		config:      config,
		cacheConfig: cacheConfig,
		db:          &gorm.DB{},
	}
	err := DBWrapper.init()
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	return DBWrapper, nil
}

func (g *DBWrapper) init() error {
	var (
		database *gorm.DB
		err      error
	)

	switch g.config.DBName {
	case postgresDB:
		// database, err  = g.NewMariaDBConnectionPool()
		// if err != nil {
		// 	err = fmt.Errorf("mariadb: %w", err)
		// }
		panic("unimplemented")
	case mysqlDB:
		database, err = g.NewMariaDBConnectionPool()
		if err != nil {
			err = fmt.Errorf("mariadb: %w", err)
		}
	default:
		return fmt.Errorf("no db selected")
	}

	if err != nil {
		fmt.Println("error connecting to DB")
		return err
	}

	g.db = database
	return err
}

func (g DBWrapper) NewMariaDBConnectionPool() (*gorm.DB, error) {
	dsn := MariaDNS(g.config.Username, g.config.Password, g.config.Host, g.config.DBName)
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold: time.Nanosecond * 1_000_000_000,
			Colorful:      true,
		},
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		CreateBatchSize:        10,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get SQL DB instance: %w", err)
	}

	sqlDB, err := database.DB()
	rawQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", g.config.DBName)
	dbCreateErr := database.Raw(rawQuery).Error
	if dbCreateErr != nil {
		return nil, dbCreateErr
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	sqlDB.SetMaxIdleConns(maxIdleConnection)
	sqlDB.SetMaxOpenConns(maxOpenConnection)
	sqlDB.SetConnMaxIdleTime(maxConnectionLifeTime)

	return database, nil
}

// func (g DBWrapper) NewPostgresConnectionPool() (*gorm.DB, error) {
// 	return database, nil
// }
