package dbconnect

import (
	"csvreader/common"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadDBToml() (config, error) {
	lConfig, lErr := common.LoadConfig[config]("./toml/dbconfig.toml")
	if lErr != nil {
		log.Println("Failed to load DB config :>", lErr)
		return lConfig, lErr
	}
	return lConfig, lErr
}

func InitDB(lDBConfig dbConfig) (*sql.DB, error) {
	var lDSN string // Data Source Name
	var lDBConnection *sql.DB

	switch lDBConfig.Driver {
	case "mysql":
		lDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			lDBConfig.User, lDBConfig.Password, lDBConfig.Host, lDBConfig.Port, lDBConfig.Name)
	default:
		return lDBConnection, fmt.Errorf("unsupported driver: %s", lDBConfig.Driver)
	}

	lDBConnection, lErr := sql.Open(lDBConfig.Driver, lDSN)
	if lErr != nil {
		log.Println("Error while sql.Open :>", lErr.Error())
		return lDBConnection, lErr
	}

	if lErr := lDBConnection.Ping(); lErr != nil {
		log.Println("Error while lDBConnection.Ping() :>", lErr.Error())
		return lDBConnection, lErr
	}

	return lDBConnection, nil
}

// Use existing sql.DB in GORM
func GORM_Connection(pDBConn *sql.DB) (*gorm.DB, error) {

	lGORM_DB, lErr := gorm.Open(mysql.New(mysql.Config{
		Conn: pDBConn,
	}), &gorm.Config{})
	if lErr != nil {
		log.Println("GORM_Connection with existing sql.DB error:>", lErr.Error())
		return nil, lErr
	}
	return lGORM_DB, nil
}
