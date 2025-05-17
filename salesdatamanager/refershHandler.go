package salesdatamanager

import (
	"csvreader/dbconnect"
	"log"
	"net/http"
)

func RefreshDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RefreshData API called")

	lErr := RefreshData()
	if lErr != nil {
		log.Printf("Error during RefreshData: %v", lErr)
		http.Error(w, "Refresh failed: "+lErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Data refreshed successfully"}`))
}

func RefreshData() error {
	log.Println("RefreshData +")

	// path declarations
	lPath := "./csvreader/historical_sales_100_records.csv"

	// Truncate or clear tables
	lErr := ClearExistingData()
	if lErr != nil {
		log.Println("Failed to clear existing data:", lErr.Error())
		return lErr
	}

	log.Println("Existing data cleared.")

	// Call existing CSV reader and inserter
	if lErr = CSV_readValidateinsert(lPath); lErr != nil {
		log.Println("Error while process the CSV_readValidateinsert :>", lErr.Error())
		return lErr
	}

	log.Println("RefreshData -")
	return nil
}

func ClearExistingData() error {

	// Delete table name list
	lTablesArr := []string{"order_items", "orders", "products", "customers", "regions"}

	for _, lTable := range lTablesArr {
		// Construct the delete query
		lQuery := "DELETE FROM " + lTable

		if lErr := dbconnect.G_DB_Conn.GORM_GO1_DB.Exec(lQuery).Error; lErr != nil {
			log.Println("Error while delete the table records from the database :>", lErr.Error())
			return lErr
		}

		log.Println("Cleared table: :>", lTable)
	}
	return nil
}
