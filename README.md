# 📊 Sales Data Management and Analytics Backend Project

A backend system built with Go (Golang) to manage, process, and analyze historical sales data efficiently. Designed with scalability, maintainability, and data integrity in mind, the system supports automated data refresh, RESTful APIs for analytics, and robust CSV ingestion.

## 🧰 Tech Stack

- **Go (Golang)** — Backend development
- **MySQL** — Relational database
- **GORM** — ORM for Go
- **Gorilla Mux** — HTTP routing
- **TOML** — Configuration management

## 📌 Key Features

### ✅ ETL & CSV Loader
- Reads and validates large CSV files
- Performs bulk insert into normalized MySQL tables
- Handles errors and logs failures with line numbers

### ✅ RESTful API
- Endpoints to refresh data and fetch revenue analytics
- Supports query parameters like date ranges and region filters
- Built using Gorilla mux

### ✅ Data Validation
- Strict checks for input formats (e.g., date: `YYYY-MM-DD`)
- Ensures referential integrity before data insert

### ✅ Data Refresh Mechanism
- API-triggered and scheduled (cron) data refresh
- Clear logs for success and failure scenarios

### ✅ Database Design
- Normalized schema with proper constraints and foreign keys
- Entities: Sales, Orders, Customers, Products, Regions

### ✅ Config Management
- Uses `config.toml` to separate environment configs (DB, scheduling, etc.)
- Simplifies deployment and environment control

### ✅ Logging & Error Handling
- Centralized logging with detailed error messages
- Fail-safe mechanisms to avoid crashing on data issues

### 🚀 Refresh Sales Data
GET http://localhost:8080/refresh

###📈 Get Revenue by Region and Date Range
GET http://localhost:8080/revenue?start_date=2024-01-01&end_date=2024-12-31

🚦 How to Run

1️⃣ Clone the Repo
git clone https://github.com/SKiranKumar22/CSV_ReadValidateInsert
cd CSV_ReadValidateInsert

2️⃣ Configure .toml
[database] # Database name
host = "localhost"
port = 3306
user = "root"
password = "yourpassword"
name = "sales_db"

3️⃣ Run the Server
go run main.go

🧑‍💻 Author
Kiran Kumar . S
MCA | Full Stack Software Developer
