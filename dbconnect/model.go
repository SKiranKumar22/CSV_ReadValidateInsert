package dbconnect

import (
	"database/sql"

	"gorm.io/gorm"
)

// global scope
type DBConnections struct {
	G_GO1_DB    *sql.DB
	GORM_GO1_DB *gorm.DB
}

var G_DB_Conn DBConnections

// Local Scope
type dbConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Name     string `toml:"name"`
	Driver   string `toml:"driver"` // fixed here
}

type config struct {
	MYSql dbConfig `toml:"MYSql"`
}
