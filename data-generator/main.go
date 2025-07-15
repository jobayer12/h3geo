package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/uber/h3-go/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"os"
)

type User struct {
	Name  string  `bson:"name"`
	Email string  `bson:"email"`
	Lat   float64 `bson:"lat"`
	Long  float64 `bson:"long"`
	H3ID  string  `bson:"h3_id"`
}

var (
	preGeneratedNames  []string
	preGeneratedEmails []string
	mongoClient        *mongo.Client
	mongoCollection    *mongo.Collection
	insertedCount      int64
)

func init() {
	// Pre-generate 10k fake names and emails
	preGeneratedNames = make([]string, 10000)
	preGeneratedEmails = make([]string, 10000)
	for i := 0; i < 10000; i++ {
		preGeneratedNames[i] = gofakeit.Name()
		preGeneratedEmails[i] = gofakeit.Email()
	}
}

func initMongoDB() {
	var err error
	_ = godotenv.Load("../.env")
	mongoURI := os.Getenv("DATABASE_CONNECTION_URI")
	fmt.Println("mongoURI", mongoURI)
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("MONGO_INITDB_DATABASE")
	if dbName == "" {
		dbName = "geo_data"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	mongoCollection = mongoClient.Database(dbName).Collection("users")
}

func main() {
	// Load .env if present
	_ = godotenv.Load()
	mongoURI := os.Getenv("DATABASE_CONNECTION_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("MONGO_INITDB_DATABASE")
	if dbName == "" {
		dbName = "geo_data"
	}
	// Global coordinates
	minLat := -90.0
	maxLat := 90.0
	minLng := -180.0
	maxLng := 180.0

	initMongoDB()

	// Configuration
	totalRecords := 200_000_000 // 200 million
	numWorkers := runtime.NumCPU() * 4
	batchSize := 10_000
	recordsPerWorker := totalRecords / numWorkers

	fmt.Printf("Starting to insert %d records using %d workers...\n", totalRecords, numWorkers)
	start := time.Now()

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go mongoWorker(i, recordsPerWorker, batchSize, minLat, maxLat, minLng, maxLng, &wg)
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("✅ Completed inserting %d records in %v\n", totalRecords, duration)
	fmt.Printf("⚡ Insert rate: %.0f records/second\n", float64(totalRecords)/duration.Seconds())

	// Create index on h3_id after all data is inserted
	createH3IDIndex()
}

func mongoWorker(id, recordCount, batchSize int, minLat, maxLat, minLng, maxLng float64, wg *sync.WaitGroup) {
	defer wg.Done()

	localRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
	var batch []interface{}
	ctx := context.Background()

	for i := 0; i < recordCount; i++ {
		lat := minLat + localRand.Float64()*(maxLat-minLat)
		lng := minLng + localRand.Float64()*(maxLng-minLng)

		h3ID, err := h3.LatLngToCell(h3.LatLng{Lat: lat, Lng: lng}, 8)
		if err != nil {
			continue
		}

		user := User{
			Name:  preGeneratedNames[localRand.Intn(len(preGeneratedNames))],
			Email: preGeneratedEmails[localRand.Intn(len(preGeneratedEmails))],
			Lat:   lat,
			Long:  lng,
			H3ID:  h3ID.String(),
		}
		batch = append(batch, user)

		if len(batch) >= batchSize {
			insertBatch(batch, ctx, id)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		insertBatch(batch, ctx, id)
	}
}

func insertBatch(batch []interface{}, ctx context.Context, id int) {
	insertCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := mongoCollection.InsertMany(insertCtx, batch)
	if err != nil {
		fmt.Printf("❌ Worker %d: Insert error: %v\n", id, err)
		return
	}

	currentTotal := atomic.AddInt64(&insertedCount, int64(len(batch)))
	if currentTotal%1_000_000 == 0 {
		percent := float64(currentTotal) / 200_000_000 * 100
		fmt.Printf("📦 Inserted %d records (%.2f%% complete)\n", currentTotal, percent)
	}
}

func createH3IDIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    map[string]interface{}{"h3_id": 1},
		Options: nil,
	}

	fmt.Println("⏳ Creating index on h3_id...")
	_, err := mongoCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("❌ Failed to create index on h3_id: %v\n", err)
		return
	}
	fmt.Println("✅ Successfully created index on h3_id.")
}
