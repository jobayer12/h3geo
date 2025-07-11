# H3 Geo Application

A real-time geolocation application that uses H3 geospatial indexing to find nearby users efficiently. Built with Go backend, Angular frontend, and MongoDB database.

## 🌟 Features

- **H3 Geospatial Indexing**: Uses Uber's H3 hexagonal grid system for efficient location-based queries
- **Real-time User Discovery**: Find users within a 5km radius using H3 resolution 8
- **Interactive Map**: Angular frontend with Leaflet.js for map visualization
- **RESTful API**: Go backend with MongoDB for data storage
- **Docker Support**: Containerized application for easy deployment
- **Scalable Architecture**: Designed for horizontal scaling

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Angular App   │    │   Go API        │    │   MongoDB       │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│   (Database)    │
│   - Leaflet     │    │   - H3 Indexing │    │   - User Data   │
│   - Map UI      │    │   - REST API    │    │   - Geospatial  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📝 Full Project Setup Guide

Follow these steps to set up, generate data, import it into MongoDB, and run both the backend (Go API) and frontend (Angular app):

---

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/h3-geo-project.git
cd h3-geo-project
```

---

### 2. Start with Docker Compose (Recommended)

This will spin up MongoDB, the Go API, and the Angular frontend in containers.

```bash
docker-compose up -d
```

- **Frontend:** http://localhost:8080  
- **API Health Check:** http://localhost:8080/api/health

---

### 3. Generate Sample Data

The project includes a data generator service to create sample user data.

**With Docker Compose:**

```bash
docker-compose --profile init-data up data-generator
```

This will run the data generator and populate the database.

**Manually (outside Docker):**

```bash
cd data-generator
go mod download
go run main.go
```

This will generate a `users_data.csv` file (if not already present) and/or insert data into MongoDB (depending on the generator's implementation).

---

### 4. Import Data into MongoDB (if needed)

If you have a `users_data.csv` and want to import it manually:

```bash
# Make sure MongoDB is running (see below for how to start it)
mongoimport --uri="mongodb://localhost:27017/h3geo" --collection=users --type=csv --headerline --file=users_data.csv
```

- Adjust the URI, database, and collection as needed.

---

### 5. Manual Setup (If Not Using Docker Compose)

#### a. Start MongoDB

```bash
docker run -d -p 27017:27017 --name mongo mongo:7.0
```

#### b. Run the Go API

```bash
cd geo-api
go mod download
go run main.go
```

- The API will be available at http://localhost:8080

#### c. Run the Angular Frontend

```bash
cd geo-map-app
npm install
npm start
```

- The frontend will be available at http://localhost:8080 (or another port if specified).

---

### 6. API Endpoints

- **Find Nearby Users:**  
  `POST /api/nearby`  
  Body:
  ```json
  {
    "lat": 40.7128,
    "long": -74.0060
  }
  ```
- **Health Check:**  
  `GET /api/health`

---

### 7. Environment Variables

Set these as needed (in your environment or Docker Compose):

| Variable                   | Description                | Default                   |
|----------------------------|----------------------------|---------------------------|
| DATABASE_CONNECTION_URI    | MongoDB connection string  | mongodb://mongo:27017     |
| MONGO_INITDB_DATABASE      | MongoDB database name      | h3geo                     |

---

### 8. Project Structure Overview

```
h3-geo-project/
├── geo-api/         # Go backend
├── geo-map-app/     # Angular frontend
├── data-generator/  # Data generator
├── docker-compose.yml
└── README.md
```

---

### Summary Table

| Step                | Command/Action                                                                 |
|---------------------|-------------------------------------------------------------------------------|
| Clone repo          | `git clone ...`                                                                |
| Start all (Docker)  | `docker-compose up -d`                                                         |
| Generate data       | `docker-compose --profile init-data up data-generator` or run generator manually|
| Import CSV (manual) | `mongoimport ...`                                                              |
| Start MongoDB       | `docker run -d -p 27017:27017 --name mongo mongo:7.0`                          |
| Run Go API          | `cd geo-api && go run main.go`                                                 |
| Run Angular app     | `cd geo-map-app && npm install && npm start`                                   |

---

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 18+ (for local development)
- MongoDB (for local development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/h3-geo-project.git
   cd h3-geo-project
   ```

2. **Start the application with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Populate the database with sample data**
   ```bash
   docker-compose --profile init-data up data-generator
   ```

4. **Access the application**
   - Frontend: http://localhost:8080
   - API Health Check: http://localhost:8080/api/health

### Manual Setup (Development)

1. **Start MongoDB**
   ```bash
   docker run -d -p 27017:27017 --name mongo mongo:7.0
   ```

2. **Build and run the Go API**
   ```bash
   cd geo-api
   go mod download
   go run main.go
   ```

3. **Build and run the Angular app**
   ```bash
   cd geo-map-app
   npm install
   npm start
   ```

## 📁 Project Structure

```
h3-geo-project/
├── geo-api/                 # Go backend application
│   ├── main.go             # Main application entry point
│   ├── Dockerfile          # Docker configuration for API
│   ├── go.mod              # Go module dependencies
│   └── go.sum              # Go module checksums
├── geo-map-app/            # Angular frontend application
│   ├── src/                # Angular source code
│   ├── package.json        # Node.js dependencies
│   └── angular.json        # Angular configuration
├── data-generator/         # Data population service
│   ├── main.go            # Data generator script
│   ├── users_data.csv     # Sample user data
│   └── Dockerfile         # Docker configuration
├── docker-compose.yml      # Local development setup
├── koyeb.yaml             # Koyeb deployment configuration
└── README.md              # This file
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_CONNECTION_URI` | MongoDB connection string | `mongodb://mongo:27017` | Yes |
| `MONGO_INITDB_DATABASE` | MongoDB database name | `h3geo` | No |

### H3 Configuration

- **Resolution**: 8 (hexagons ~36km²)
- **Search Radius**: 5 rings (~5km coverage)
- **Index**: H3 ID field for efficient queries

## 🗺️ API Endpoints

### Find Nearby Users
```http
POST /api/nearby
Content-Type: application/json

{
  "lat": 40.7128,
  "long": -74.0060
}
```

**Response:**
```json
{
  "users": [
    {
      "id": "user123",
      "name": "John Doe",
      "email": "john@example.com",
      "lat": 40.7128,
      "long": -74.0060,
      "h3_id": "8828308281fffff"
    }
  ],
  "total": 1
}
```

### Health Check
```http
GET /api/health
```

**Response:**
```json
{
  "status": "healthy"
}
```
