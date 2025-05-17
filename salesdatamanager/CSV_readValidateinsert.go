package salesdatamanager

import (
	"csvreader/dbconnect"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm/clause"
)

func CSV_readValidateinsert(pPath string) error {
	log.Println("CSV_readValidateinsert +")

	// variable declaration
	var lCSVRecordArr []CSVRecord

	// open the file
	lFile, lErr := os.Open(pPath)
	if lErr != nil {
		log.Println("Error while open CSV file from path :>", lErr.Error())
		return lErr
	}

	// after reading close the file
	defer lFile.Close()

	lReader := csv.NewReader(lFile)

	// Skip header to store the value from DB.
	if lHeaderValues, lErr := lReader.Read(); lErr != nil {
		log.Println("Error while read the header from CSV :>", lErr.Error())
		return lErr
	} else {
		// print the header values for references
		log.Println("Header values :>", lHeaderValues)
	}

	for {

		lRow, lErr := lReader.Read()

		// To get to known end of the file
		if lErr == io.EOF {
			log.Println("EOF is the error returned by Read when no more input is available :>", io.EOF)
			break
		}

		// Error while read the file
		if lErr != nil {
			log.Println("Error while read the record from the file :>", lErr.Error())
			return lErr
		}

		// read the values from the file rows and store the records in array.
		lRecord := CSVRecord{
			OrderID:         lRow[0],
			ProductID:       lRow[1],
			CustomerID:      lRow[2],
			ProductName:     lRow[3],
			Category:        lRow[4],
			Region:          lRow[5],
			DateOfSale:      lRow[6],
			QuantitySold:    lRow[7],
			UnitPrice:       lRow[8],
			Discount:        lRow[9],
			ShippingCost:    lRow[10],
			PaymentMethod:   lRow[11],
			CustomerName:    lRow[12],
			CustomerEmail:   lRow[13],
			CustomerAddress: lRow[14],
		}

		lCSVRecordArr = append(lCSVRecordArr, lRecord)
	}

	log.Println("Before Validate CSV data :>", len(lCSVRecordArr))

	lValidateCSVData, lErr := validateCSVRecords(lCSVRecordArr)
	if lErr != nil {
		log.Println("Error while validate CSV Records from the file :>", lErr.Error())
		return lErr
	}

	log.Println("After Validate CSV data :>", len(lValidateCSVData))

	lErr = BulkInsertCSVData(lValidateCSVData)
	if lErr != nil {
		log.Println("Error while bulk insert CSV data from file :>", lErr.Error())
		return lErr
	}

	log.Println("CSV_readValidateinsert -")
	return nil
}

func BulkInsertCSVData(records []CSVRecord) error {

	var (
		customersMap = make(map[string]Customer)
		productsMap  = make(map[string]Product)
		regionsMap   = make(map[string]Region)
		orders       []Order
		orderItems   []OrderItem
	)

	// Deduplicate and prepare
	for _, r := range records {

		// --- Customers ---
		if _, exists := customersMap[r.CustomerID]; !exists {
			customersMap[r.CustomerID] = Customer{
				ID:      r.CustomerID,
				Name:    r.CustomerName,
				Email:   r.CustomerEmail,
				Address: r.CustomerAddress,
			}
		}

		// --- Products ---
		if _, exists := productsMap[r.ProductID]; !exists {
			productsMap[r.ProductID] = Product{
				ID:       r.ProductID,
				Name:     r.ProductName,
				Category: r.Category,
			}
		}

		// --- Regions ---
		if _, exists := regionsMap[r.Region]; !exists {
			regionsMap[r.Region] = Region{Name: r.Region}
		}
	}

	// --- Insert customers/products/regions if not exist ---
	if err := upsertCustomers(customersMap); err != nil {
		return err
	}
	if err := upsertProducts(productsMap); err != nil {
		return err
	}
	if err := upsertRegions(regionsMap); err != nil {
		return err
	}

	// Map region name to ID after insertion
	regionIDs := make(map[string]int)
	for _, r := range regionsMap {
		var region Region
		if err := dbconnect.G_DB_Conn.GORM_GO1_DB.Where("name = ?", r.Name).First(&region).Error; err != nil {
			return err
		}
		regionIDs[r.Name] = region.ID
	}

	// Build orders and orderItems
	for _, r := range records {
		date, _ := time.Parse("2006-01-02", r.DateOfSale)
		qty, _ := strconv.Atoi(r.QuantitySold)
		unitPrice, _ := strconv.ParseFloat(r.UnitPrice, 64)
		discount, _ := strconv.ParseFloat(r.Discount, 64)
		shipping, _ := strconv.ParseFloat(r.ShippingCost, 64)

		orders = append(orders, Order{
			ID:            r.OrderID,
			Date:          date,
			CustomerID:    r.CustomerID,
			RegionID:      regionIDs[r.Region],
			PaymentMethod: r.PaymentMethod,
			ShippingCost:  shipping,
		})

		orderItems = append(orderItems, OrderItem{
			OrderID:   r.OrderID,
			ProductID: r.ProductID,
			Quantity:  qty,
			UnitPrice: unitPrice,
			Discount:  discount,
		})
	}

	// --- Bulk insert Orders & OrderItems ---
	if err := dbconnect.G_DB_Conn.GORM_GO1_DB.Create(&orders).Error; err != nil {
		return fmt.Errorf("bulk insert orders failed: %v", err)
	}
	if err := dbconnect.G_DB_Conn.GORM_GO1_DB.Create(&orderItems).Error; err != nil {
		return fmt.Errorf("bulk insert order items failed: %v", err)
	}

	return nil
}

func upsertCustomers(pCustomers map[string]Customer) error {

	for _, c := range pCustomers {
		dbconnect.G_DB_Conn.GORM_GO1_DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&c)
	}
	return nil

}

func upsertProducts(pProducts map[string]Product) error {

	for _, p := range pProducts {
		dbconnect.G_DB_Conn.GORM_GO1_DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&p)
	}

	return nil
}

func upsertRegions(pRegions map[string]Region) error {

	for _, r := range pRegions {
		dbconnect.G_DB_Conn.GORM_GO1_DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoNothing: true,
		}).Create(&r)
	}

	return nil
}
