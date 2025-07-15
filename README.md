# H3 Geo Application

A real-time geolocation application using Uber's H3 geospatial indexing to efficiently find nearby users. The project consists of a Go backend API, an Angular frontend, and a MongoDB database. Data generation is handled by a Go-based generator.

---

## 🌟 Features

- **H3 Geospatial Indexing:** Efficient location-based queries using Uber's H3 hexagonal grid.
- **Real-time User Discovery:** Find users within a 5km radius (H3 resolution 8).
- **Interactive Map:** Angular frontend with map visualization.
- **RESTful API:** Go backend with MongoDB for data storage.
- **Scalable Data Generation:** Generate and insert millions of users for testing.

---

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Angular App   │    │   Go API        │    │   MongoDB       │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│   (Database)    │
│   - Map UI      │    │   - H3 Indexing │    │   - User Data   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 📁 Project Structure

```
h3-geo-project/
├── geo-api/         # Go backend API
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── geo-map-app/     # Angular frontend
│   ├── src/
│   ├── package.json
│   └── angular.json
├── data-generator/  # Go data generator
│   ├── main.go
│   ├── go.mod
│   └── go.sum
└── README.md
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+ (for backend and data generator)
- Node.js 18+ (for frontend)
- MongoDB (local or remote instance)

---

### 1. Start MongoDB

You can run MongoDB locally (default: `mongodb://localhost:27017`).  
Create a database named `geo_data`.

---

### 2. Generate and Insert Sample Data

The data generator inserts millions of fake users directly into MongoDB.

```bash
cd data-generator
go mod download
go run main.go
```

- This will connect to MongoDB at `mongodb://localhost:27017` and insert data into the `geo_data.users` collection.
- The generator creates up to 200 million users (configurable in `main.go`).

---

### 3. Run the Go API

```bash
cd geo-api
go mod download
go run main.go
```

- The API will be available at [http://localhost:8080](http://localhost:8080) by default.
- The backend connects to MongoDB at `mongodb://localhost:27017` and serves the Angular app (if built) from the `static/` directory.

---

### 4. Run the Angular Frontend

```bash
cd geo-map-app
npm install
npm run start
```

- The frontend will be available at [http://localhost:4200](http://localhost:4200) (default Angular port).
---

## 🗺️ API Endpoints

### Find Nearby Users

`POST /api/nearby`

**Request Body:**
```json
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

---

### Health Check

`GET /api/health`

**Response:**
```json
{
  "status": "healthy"
}
```

---

## 🔧 Configuration

- The Go backend and data generator connect to MongoDB at `mongodb://localhost:27017` and use the `geo_data` database by default.
- The backend server listens on port `8080` (can be changed via the `PORT` environment variable).

---

## 📝 Notes

- The data generator is optimized for high-volume inserts and will create an index on `h3_id` after data insertion.
- All geospatial queries use H3 resolution 8 for ~5km coverage.
- The backend serves static files from the `static/` directory if present (for production builds of the Angular app).

---

## ✨ Contributing

Feel free to open issues or submit pull requests for improvements!
