MyHomeInventory
A lightweight, single-page web application for managing a simple home inventory — groceries, food, and household items — with a MySQL backend and a Go server.

📚 Features
Add new items to your inventory

View all items in a simple table

Increment (+) and decrement (−) item quantities

Track item usage over time

Lightweight, fast, and no heavy frameworks

🚀 Requirements
Go (1.20+)

MySQL Server

MySQL Workbench (optional)

A modern web browser (Chrome, Firefox, Edge)

🛠️ Setup Instructions
Clone this repository

bash
Copy
Edit
git clone https://github.com/yourusername/myhomeinventory.git
cd myhomeinventory
Install dependencies

go
Copy
Edit
go mod tidy
Configure your environment

Create a .env file in the project root:

ini
Copy
Edit
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=inventory_db
APP_HOST=localhost
APP_PORT=8080
⚠️ Note: Make sure your MySQL server is running and the database inventory_db exists.

Create the table

If it doesn't exist already, create the item_inventory table:

pgsql
Copy
Edit
CREATE TABLE item_inventory (
    id INT AUTO_INCREMENT PRIMARY KEY,
    itemName VARCHAR(255) NOT NULL,
    itemQTY INT NOT NULL,
    minimumQTY INT NOT NULL,
    itemUsedToDate INT NOT NULL DEFAULT 0,
    itemType VARCHAR(100),
    createDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    lastModifiedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
Run the server

arduino
Copy
Edit
go run ./cmd/inventory
Open the application

Open your browser and navigate to:

arduino
Copy
Edit
http://localhost:8080
You should see the Home Inventory web interface.

📁 Project Structure
csharp
Copy
Edit
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
└── README.md
🔧 Configuration
All configuration is handled via the .env file:

Database user/password/host/port

Application host/port

📦 Dependencies
github.com/go-sql-driver/mysql — MySQL driver for Go

github.com/joho/godotenv — .env loader

Installed automatically via go mod tidy.

📝 License
MIT License. See LICENSE file for details.
