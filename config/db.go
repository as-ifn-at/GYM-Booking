package config

type DBConfigOptions struct {
	Host         string
	Username     string
	Password     string
	DBName       string
	AppName      string
	DBToUse      string
	IsLogEnabled bool
	RetryCount   int
	ErrorCode    ErrorCode
}

type Error int
type ErrorCode interface {
	GetError(code string) Error
}
