package dbconnect

import "log"

// open db connections
func OpenDBConnections() error {
	log.Println("OpenDBConnections +")

	// Load DB Toml
	lConfig, lErr := LoadDBToml()
	if lErr != nil {
		log.Println("Error while load the DB TOML :>", lErr.Error())
		return lErr
	}

	// Open Connection
	G_DB_Conn.G_GO1_DB, lErr = InitDB(lConfig.MYSql)
	if lErr != nil {
		log.Println("Error while the InitDB lConfig.MYSql:>", lErr.Error())
		return lErr
	}

	G_DB_Conn.GORM_GO1_DB, lErr = GORM_Connection(G_DB_Conn.G_GO1_DB)
	if lErr != nil {
		log.Println("Error while the InitDB GORM G_DB_Conn.G_GO1_DB:>", lErr.Error())
		return lErr
	}

	log.Println("OpenDBConnections -")
	return nil
}

// Close DB connection
func CloseDBConnections() {
	log.Println("CloseDBConnections +")
	var lErr error

	if lErr = G_DB_Conn.G_GO1_DB.Close(); lErr != nil {
		log.Println("Error while CloseDBConnections G_DB_Conn.G_GO1_DB.Close() :>", lErr.Error())
	}

	log.Println("CloseDBConnections -")
}
