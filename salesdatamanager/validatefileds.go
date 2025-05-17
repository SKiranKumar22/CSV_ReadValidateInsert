package salesdatamanager

import (
	"fmt"
	"net/mail"
	"strconv"
	"time"
)

func validateCSVRecords(records []CSVRecord) ([]CSVRecord, error) {
	var validRecords []CSVRecord

	for i, r := range records {
		// Required fields check
		if isEmpty(r.OrderID, r.CustomerID, r.ProductID) {
			return nil, fmt.Errorf("row %d: missing required ID field(s)", i+1)
		}

		// Date validation
		if err := validateDate(r.DateOfSale); err != nil {
			return nil, fmt.Errorf("row %d: %v", i+1, err)
		}

		// Numeric fields
		if err := validateNumericFields(i+1, r.QuantitySold, r.UnitPrice, r.Discount, r.ShippingCost); err != nil {
			return nil, err
		}

		// Email validation
		if err := validateEmail(r.CustomerEmail); err != nil {
			return nil, fmt.Errorf("row %d: %v", i+1, err)
		}

		validRecords = append(validRecords, r)
	}

	return validRecords, nil
}

func isEmpty(fields ...string) bool {
	for _, field := range fields {
		if field == "" {
			return true
		}
	}
	return false
}

func validateDate(dateStr string) error {
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return fmt.Errorf("invalid date format: %s", dateStr)
	}
	return nil
}

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}
	return nil
}

func validateNumericFields(row int, fields ...string) error {
	for idx, field := range fields {
		if _, err := strconv.ParseFloat(field, 64); err != nil {
			return fmt.Errorf("row %d: invalid numeric value in field %d (%s)", row, idx+1, field)
		}
	}
	return nil
}
