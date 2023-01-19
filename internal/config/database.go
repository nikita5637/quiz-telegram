package config

import "fmt"

const (
	// username:password@protocol(ip:port)/dbname
	mysqlDSNFormat = "%s:%s@tcp(%s:%d)/%s?parseTime=true"

	// DriverMySQL ...
	DriverMySQL = "mysql"
)

// DatabaseConfig ...
type DatabaseConfig struct {
	DatabaseAddress  string `toml:"database_address"`
	DatabaseName     string `toml:"database_name"`
	DatabaseUserName string `toml:"database_username"`
	DatabasePort     uint16 `toml:"database_port"`
	Driver           string `toml:"driver"`
}

// GetDatabaseDSN ...
func GetDatabaseDSN() string {
	databasePassword := GetSecretValue(DatabasePassword)

	switch globalConfig.Driver {
	case DriverMySQL:
		return fmt.Sprintf(mysqlDSNFormat,
			globalConfig.DatabaseUserName,
			databasePassword,
			globalConfig.DatabaseAddress,
			globalConfig.DatabasePort,
			globalConfig.DatabaseName,
		)
	}

	return ""
}
