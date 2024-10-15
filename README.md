# Golang Data Management System

## Overview
This project is a Golang-based system that imports data from an Excel file, stores it into a MySQL database, and caches the data in Redis for efficient retrieval. It includes CRUD functionality to manage the imported data, allowing users to view, edit, update, and delete records while maintaining synchronization between MySQL and Redis.

## Project Flow

### 1. Importing Excel Data
- **Upload the Excel File**: Users can upload an Excel file using an API development tool like Postman.
- **Data Parsing**: The uploaded Excel file is parsed and validated to ensure adherence to required column headers and data types.
- **Asynchronous Processing**: The data parsing and structuring are processed asynchronously to ensure scalability, allowing the system to handle large datasets without blocking operations.

### 2. Storing Data
- **MySQL Database**: 
  - Connects to a MySQL database and creates a table to store the parsed data.
  - Implements functions to insert the structured data into the MySQL table.
- **Redis Cache**:
  - Connects to Redis to cache the imported data for quick retrieval.
  - Sets an expiration for cached data (5 minutes) to ensure data freshness and reduce memory usage.

### 3. Viewing Imported Data
- **Fetch from Redis**: 
  - Provides an API endpoint to view the imported data. The endpoint first attempts to retrieve data from Redis.
  - If the data is not found in Redis, it fetches it from the MySQL database and updates the cache.
- **Display Data**: Data is returned in a readable format for users to easily analyze the imported information.

### 4. Editing Records
- **Update Specific Records**: Allows users to edit individual records using an API endpoint.
- **Synchronize Changes**:
  - Updates the edited record in both MySQL and Redis to maintain data consistency.
  - Redis cache is updated immediately to reflect the changes made in the database.

### 5. Error Handling and Validation
- **File Validation**: Ensures the uploaded Excel file meets the expected format and structure before proceeding with parsing.
- **Graceful Error Handling**: Implements mechanisms to handle errors during file upload, data parsing, database operations, and cache interactions to ensure the system remains stable.

### 6. Optimization and Scalability
- **Database and Cache Optimization**: Uses efficient database queries and caching mechanisms to improve the speed and performance of data retrieval.
- **Scalable Architecture**: Designed to scale horizontally to handle increased data volume and API traffic, ensuring the application can manage growth without significant performance degradation.

## Prerequisites
- **Golang**: v1.18 or above.
- **MySQL**: v8.0 or above.
- **Redis**: v6.0 or above.
- **Golang Packages**: 
  - `github.com/gin-gonic/gin`: For building the web server.
  - `github.com/go-sql-driver/mysql`: MySQL driver for Golang.
  - `github.com/go-redis/redis`: Redis client for Golang.
  - `github.com/tealeg/xlsx`: For parsing Excel files.

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/GautamViki/crm-backend.git
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Configure the database and Redis connection in `config.yaml`.
4. Start the server:
   ```bash
   go run main.go
   ```
5. Use Postman or another API tool to interact with the endpoints for uploading files, viewing data, and editing records.

## API Endpoints
- **Upload Excel Data**: `POST /upload`
- **View Imported Data**: `GET /customers`
- **Edit Record**: `PUT /customers/:id`
- **Delete Record**: `DELETE /customers/:id`
- **View Cached Data**: `GET /customers/cache`
