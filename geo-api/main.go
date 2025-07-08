package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/uber/h3-go/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID    string  `json:"id" bson:"id"`
	Name  string  `json:"name" bson:"name"`
	Email string  `json:"email" bson:"email"`
	Lat   float64 `json:"lat" bson:"lat"`
	Long  float64 `json:"long" bson:"long"`
	H3ID  string  `json:"h3_id" bson:"h3_id"`
}

type NearbyRequest struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type NearbyResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

var client *mongo.Client
var collection *mongo.Collection

func main() {
	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(""))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		log.Fatal(err)
	}

	// Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get collection
	collection = client.Database("h3geo").Collection("users")

	// Create index on h3_id for faster queries
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"h3_id": 1,
		},
	}
	collection.Indexes().CreateOne(context.TODO(), indexModel)

	// Setup routes
	r := mux.NewRouter()
	
	// API routes
	r.HandleFunc("/api/nearby", getNearbyUsers).Methods("POST")
	r.HandleFunc("/api/health", healthCheck).Methods("GET")

	// Serve static files (Angular app)
	r.PathPrefix("/").Handler(http.HandlerFunc(serveStaticFiles))

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// serveStaticFiles serves the Angular app static files
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Don't serve static files for API routes
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Get the file path
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Construct the full file path - Angular files are in static/browser/
	filePath := filepath.Join("static", "browser", path)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If file doesn't exist, serve index.html for SPA routing
		filePath = filepath.Join("static", "browser", "index.html")
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}


func getNearbyUsers(w http.ResponseWriter, r *http.Request) {
	var req NearbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert lat/long to H3 ID at resolution 8
	h3ID, err := h3.LatLngToCell(h3.LatLng{
		Lat: req.Lat,
		Lng: req.Long,
	}, 8) // Changed from 7 to 8
	if err != nil {
		http.Error(w, "Invalid coordinates", http.StatusBadRequest)
		return
	}

	// Get k-ring neighbors for 5km coverage
	// For resolution 8, approximately k=6-7 rings needed for 5km
	neighbors, err := h3.GridDisk(h3ID, 5) // Changed from 1 to 7
	if err != nil {
		http.Error(w, "Error calculating neighbors", http.StatusInternalServerError)
		return
	}

	// Build query for MongoDB
	var h3IDStrings []string
	for _, neighbor := range neighbors {
		h3IDStrings = append(h3IDStrings, neighbor.String())
	}

	filter := bson.M{"h3_id": bson.M{"$in": h3IDStrings}}

	// Find users in nearby hexagons
	cursor, err := collection.Find(context.TODO(), filter, options.Find())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		http.Error(w, "Error decoding results", http.StatusInternalServerError)
		return
	}

	response := NearbyResponse{
		Users: users,
		Total: len(users),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
