package main

import (
	"csvreader/dbconnect"
	"csvreader/networkconfig"
	"csvreader/salesdatamanager"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Func is used to update the next sehedule run time
func scheduleDailyTask(hour, min int, task func()) {
	log.Println("scheduleDailyTask +")
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())

			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			duration := next.Sub(now)
			log.Printf("Next run scheduled at: %s (in %s)", next.Format(time.RFC1123), duration)

			time.Sleep(duration)
			task()
		}
	}()
	log.Println("scheduleDailyTask -")
}

func main() {

	var logDir = "./log"
	// lPath := "./csvreader/historical_sales_100_records.csv"

	// Create filename with lTimestamp
	lTimestamp := time.Now().UnixNano()
	lTimestampString := strconv.Itoa(int(lTimestamp))
	lLogFileName := filepath.Join(logDir, "logfile_"+lTimestampString+".log")

	// Open the log lFile (create if not exists)
	lFile, lErr := os.OpenFile(lLogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatal("Failed to open log file:", lErr)
	}
	defer lFile.Close()

	// Set log output to file
	log.SetOutput(lFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Hello world")
	if lErr = dbconnect.OpenDBConnections(); lErr != nil {
		log.Println("Error while processing OpenDBConnections :>", lErr.Error())

		// Close the db connection even error occurs
		dbconnect.CloseDBConnections()
	} else {

		// When all the process are over close the Global DB connection
		defer dbconnect.CloseDBConnections()

		// Schedule daily refresh at 5.5 AM
		scheduleDailyTask(5, 5, func() {
			log.Println("Scheduled RefreshData triggered")
			salesdatamanager.RefreshData()
		})

		log.Println("DB Connection Successfull")

		router := mux.NewRouter()

		// Register all your API routes separately
		RegisterRoutes(router)

		// Wrap router with CORS middleware
		handler := networkconfig.CorsMiddleware(router)

		// RefreshData()

		log.Println("Starting server on :8080")

		// Process with out API
		// if lErr = salesdatamanager.CSV_readValidateinsert(lPath); lErr != nil {
		// 	log.Println("Error while process the CSV_readValidateinsert :>", lErr.Error())
		// }

		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatal(err)
		}
	}

}

// Function to register API routes on a router
func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/refresh", salesdatamanager.RefreshDataHandler).Methods("GET")
	router.HandleFunc("/revenue", salesdatamanager.RevenueHandler).Methods("GET")
}
