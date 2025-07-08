package main

import (
	"encoding/csv"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/uber/h3-go/v4"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID    string
	Name  string
	Email string
	Lat   float64
	Long  float64
	H3ID  string
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())

	// Bangladesh boundaries (approximate)
	minLat := 20.670883
	maxLat := 26.446526
	minLng := 88.028336
	maxLng := 92.672668

	// Create CSV file
	file, err := os.Create("users_data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"name", "email", "lat", "long", "h3_id"}
	writer.Write(header)

	// Generate 2 million records
	fmt.Println("Generating 2 million user records...")
	for i := 0; i < 2000000; i++ {
		// Generate random coordinates within Bangladesh
		lat := minLat + rand.Float64()*(maxLat-minLat)
		lng := minLng + rand.Float64()*(maxLng-minLng)

		// Generate H3 ID at resolution 7
		h3ID, err := h3.LatLngToCell(h3.LatLng{
			Lat: lat,
			Lng: lng,
		}, 8)
		if err != nil {
			fmt.Printf("Error generating H3 ID for lat: %f, lng: %f: %v\n", lat, lng, err)
			continue
		}

		user := User{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Lat:   lat,
			Long:  lng,
			H3ID:  h3ID.String(),
		}

		// Write to CSV
		record := []string{
			user.Name,
			user.Email,
			strconv.FormatFloat(user.Lat, 'f', 6, 64),
			strconv.FormatFloat(user.Long, 'f', 6, 64),
			user.H3ID,
		}
		writer.Write(record)

		// Progress indicator
		if i%100000 == 0 {
			fmt.Printf("Generated %d records\n", i)
		}
	}

	fmt.Println("Data generation completed! File saved as users_data.csv")
}