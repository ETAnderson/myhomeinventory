 # MyHomeInventory

A lightweight, single-page web application for managing a simple home inventory — groceries, food, and household items — with a MySQL backend and a Go server.

---

## Features

- Add new items to your inventory
- View all items in a simple table
- Increment (`+`) and decrement (`−`) item quantities
- Track item usage over time
- Lightweight, fast, and no heavy frameworks

---

## Requirements

- [Go 1.20+](https://golang.org/dl/)
- [MySQL Server](https://dev.mysql.com/downloads/mysql/)
- [MySQL Workbench (optional)](https://dev.mysql.com/downloads/workbench/)
- A modern web browser (Chrome, Firefox, Edge)

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/myhomeinventory.git
cd myhomeinventory
```
2. Install Go Dependencies
```bash
go mod tidy
```
3. Configure Environment Variables
Copy the example environment file:

```bash
cp env.example .env
```

```dotenv
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=inventory_db
APP_HOST=localhost
APP_PORT=8080
```
⚠️ Note: Ensure your MySQL server is running and the inventory_db database exists.

4. Run the Server

```bash
go run ./cmd/inventory
```
5. Open the Application
Visit:

```arduino
this is defined by the APP_HOST and APP_PORT environment variables like http://localhost:8080/
```
You should see the Home Inventory web interface.

Project Structure
```text
myhomeinventory/
├── cmd/
│   └── inventory/        # Main entry point
├── internal/
│   └── inventory/        # Database logic, models, helpers
├── server/               # HTTP handlers and router
├── static/               # Frontend (HTML/CSS/JS)
│   ├── app.js
│   ├── styles.css
├── templates/            # HTML templates
│   └── index.html
├── .env                   # Environment variables (not committed)
├── go.mod                 # Go module file
└── README.me
```
Configuration
```text
All application configuration is handled via the .env file:

Database credentials (user, password, host, port, name)

Application host and port
```
Dependencies
```text
github.com/go-sql-driver/mysql — MySQL driver for Go

github.com/joho/godotenv — Environment variable loader
```
Installed automatically via:

```bash
 go mod tidy
```
License
```
This project is licensed under the MIT License.
```
