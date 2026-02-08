# Go Kasir API

A Point-of-Sale (POS) REST API built with Go and PostgreSQL. This API provides endpoints for managing products, categories, and transactions in a cashier system.

## Features

- 🛍️ **Product Management** - CRUD operations for products
- 📁 **Category Management** - Organize products by categories
- 💰 **Transaction Processing** - Process sales with automatic stock updates
- 📊 **Sales Reports** - Generate revenue and best-selling product reports
- 🔒 **Transaction Safety** - Atomic transactions with automatic rollback on errors
- ⚡ **Optimized Queries** - Batch operations to minimize database round-trips

## Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Configuration:** Viper (supports .env files)
- **Architecture:** Clean Architecture (Handlers → Services → Repositories)

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Git

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd go-kasir-api
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_CONN=postgresql://username:password@localhost:5432/database_name?sslmode=disable
```

### 4. Set up the database

Run the database migration script located in `database/init.sql`:

```bash
psql -U username -d database_name -f database/init.sql
```

### 5. Run the application

```bash
go run main.go
```

The server will start on `http://0.0.0.0:8080`

## API Endpoints

### Health Check

#### `GET /health`
Check if the API is running.

**Response:**
```json
{
  "status": "ok"
}
```

---

### Products

#### `GET /api/products`
Get all products with optional name filtering.

**Query Parameters:**
- `name` (optional) - Filter products by name (case-insensitive)

**Response:**
```json
[
  {
    "id": 1,
    "name": "Coca Cola",
    "price": 5000,
    "stock": 100,
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "Beverages",
      "description": "Drinks and beverages"
    }
  }
]
```

#### `POST /api/products`
Create a new product.

**Request Body:**
```json
{
  "name": "Pepsi",
  "price": 5000,
  "stock": 50,
  "category_id": 1
}
```

**Response:**
```json
{
  "id": 2,
  "name": "Pepsi",
  "price": 5000,
  "stock": 50,
  "category_id": 1
}
```

#### `GET /api/products/{id}`
Get a single product by ID.

**Response:**
```json
{
  "id": 1,
  "name": "Coca Cola",
  "price": 5000,
  "stock": 100,
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "Beverages",
    "description": "Drinks and beverages"
  }
}
```

#### `PUT /api/products/{id}`
Update a product.

**Request Body:**
```json
{
  "name": "Coca Cola 500ml",
  "price": 5500,
  "stock": 120,
  "category_id": 1
}
```

#### `DELETE /api/products/{id}`
Delete a product.

**Response:**
```json
{
  "message": "Product deleted"
}
```

---

### Categories

#### `GET /api/categories`
Get all categories.

**Response:**
```json
[
  {
    "id": 1,
    "name": "Beverages",
    "description": "Drinks and beverages"
  }
]
```

#### `POST /api/categories`
Create a new category.

**Request Body:**
```json
{
  "name": "Snacks",
  "description": "Chips and snacks"
}
```

#### `GET /api/categories/{id}`
Get a single category by ID.

#### `PUT /api/categories/{id}`
Update a category.

**Request Body:**
```json
{
  "name": "Beverages",
  "description": "Cold and hot drinks"
}
```

#### `DELETE /api/categories/{id}`
Delete a category.

---

### Transactions

#### `POST /api/transactions`
Create a new transaction (checkout).

**Request Body:**
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```

**Response:**
```json
{
  "id": 1,
  "total_amount": 15000,
  "created_at": "2026-02-08T14:30:00Z",
  "details": [
    {
      "id": 1,
      "transaction_id": 1,
      "product_id": 1,
      "product_name": "Coca Cola",
      "quantity": 2,
      "subtotal": 10000
    },
    {
      "id": 2,
      "transaction_id": 1,
      "product_id": 3,
      "product_name": "Chips",
      "quantity": 1,
      "subtotal": 5000
    }
  ]
}
```

**Features:**
- ✅ Validates product existence
- ✅ Checks stock availability
- ✅ Updates stock automatically
- ✅ Atomic transaction (all-or-nothing)
- ✅ Automatic rollback on errors

---

### Reports

#### `GET /api/transactions/reports/today`
Get today's sales report.

**Response:**
```json
{
  "total_revenue": 150000,
  "total_transaction": 12,
  "best_selling_product": {
    "product_name": "Coca Cola",
    "quantity_sold": 25
  }
}
```

#### `GET /api/transactions/reports`
Get sales report for a date range.

**Query Parameters:**
- `start_date` (required) - Start date (YYYY-MM-DD)
- `end_date` (required) - End date (YYYY-MM-DD)

**Example:**
```
GET /api/transactions/reports?start_date=2026-02-01&end_date=2026-02-08
```

**Response:**
```json
{
  "total_revenue": 500000,
  "total_transaction": 45,
  "best_selling_product": {
    "product_name": "Coca Cola",
    "quantity_sold": 85
  }
}
```

## Database Schema

### Products
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Categories
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transactions
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    total_amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transaction Details
```sql
CREATE TABLE transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    subtotal INTEGER NOT NULL
);
```

## Project Structure

```
go-kasir-api/
├── database/           # Database initialization scripts
├── handlers/           # HTTP request handlers
├── models/            # Data models/structs
├── repositories/      # Database operations
├── services/          # Business logic
├── main.go           # Application entry point
├── .env              # Environment variables
├── go.mod            # Go module definition
└── README.md         # This file
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET/PUT/DELETE
- `201 Created` - Successful POST
- `400 Bad Request` - Invalid request body
- `404 Not Found` - Resource not found
- `405 Method Not Allowed` - Invalid HTTP method
- `500 Internal Server Error` - Server error

Error responses include a message:
```json
{
  "error": "Product not found"
}
```

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o kasir-api main.go
```

### Docker Deployment (Optional)
```bash
docker build -t go-kasir-api .
docker run -p 8080:8080 --env-file .env go-kasir-api
```

## Performance Optimizations

- **Batch Queries**: Transaction creation uses batch SELECT and INSERT operations
- **Connection Pooling**: PostgreSQL connection pooling for better performance
- **Indexed Queries**: Database queries utilize indexes on foreign keys and timestamps

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.