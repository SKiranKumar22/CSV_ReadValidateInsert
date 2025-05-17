package salesdatamanager

import (
	"csvreader/dbconnect"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func RevenueHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("RevenueHandler +")

	start := r.URL.Query().Get("start_date")
	end := r.URL.Query().Get("end_date")

	if start == "" || end == "" {
		log.Println("Missing start_date or end_date", http.StatusBadRequest)
		http.Error(w, "Missing start_date or end_date", http.StatusBadRequest)
		return
	}

	// Define the expected date layout
	const layout = "2006-01-02"

	// Parse and validate start_date
	lStartDate, lErr := time.Parse(layout, start)
	if lErr != nil {
		http.Error(w, "Invalid start_date format. Expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Parse and validate end_date
	lEndDate, lErr := time.Parse(layout, end)
	if lErr != nil {
		log.Println("Invalid end_date format. Expected YYYY-MM-DD", http.StatusBadRequest)
		http.Error(w, "Invalid end_date format. Expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Check startDate is before or equal endDate
	if lStartDate.After(lEndDate) {
		log.Println("start_date cannot be after end_date", http.StatusBadRequest)
		http.Error(w, "start_date cannot be after end_date", http.StatusBadRequest)
		return
	}

	// Your existing GORM query here, using startDate and endDate:
	var lRevenue float64

	lErr = dbconnect.G_DB_Conn.GORM_GO1_DB.
		Table("order_items").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.order_date BETWEEN ? AND ?", lStartDate, lEndDate).
		Select("COALESCE(SUM(order_items.quantity * (order_items.unit_price - order_items.discount)),'0')").
		Scan(&lRevenue).Error

	if lErr != nil {
		log.Println("Failed to calculate revenue  :> " + lErr.Error())
		http.Error(w, "Failed to calculate revenue :> "+lErr.Error(), http.StatusInternalServerError)
		return
	}

	lResp := map[string]interface{}{
		"start_date": start,
		"end_date":   end,
		"revenue":    lRevenue,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lResp)

	log.Println("RevenueHandler -")
}
