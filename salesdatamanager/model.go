package salesdatamanager

import "time"

type Customer struct {
	ID      string `gorm:"primaryKey"`
	Name    string
	Email   string
	Address string
}

type Product struct {
	ID       string `gorm:"primaryKey"`
	Name     string
	Category string
}

type Region struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"unique"`
}

type Order struct {
	ID            string `gorm:"primaryKey"`
	Date          time.Time
	CustomerID    string
	Customer      Customer
	PaymentMethod string
	ShippingCost  float64
	RegionID      int
	Region        Region
}

type OrderItem struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice float64
	Discount  float64
	Order     Order
	Product   Product
}

type CSVRecord struct {
	OrderID         string
	ProductID       string
	CustomerID      string
	ProductName     string
	Category        string
	Region          string
	DateOfSale      string
	QuantitySold    string
	UnitPrice       string
	Discount        string
	ShippingCost    string
	PaymentMethod   string
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string
}
