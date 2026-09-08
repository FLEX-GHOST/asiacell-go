package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func loadClient() (*asiacell.Client, error) {
	client, err := asiacell.NewClient(asiacell.WithTimeout(20 * time.Second))
	if err != nil {
		return nil, err
	}

	if err := client.LoadSessionFromFile("session.json"); err == nil {
		return client, nil
	}

	return nil, fmt.Errorf("no saved session found (run 01_auth_and_session first)")
}

func main() {
	client, err := loadClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== 1. Asiacell Covered Cities & Governorates ===")
	cities, err := client.GetCities(ctx)
	if err != nil {
		fmt.Printf("GetCities error: %v\n", err)
		return
	}

	for i, c := range cities {
		fmt.Printf("[%2d] %-15s (City ID: %d)\n", i+1, c.Name, c.ID)
	}

	targetCityID := 5
	cityName := "Baghdad"
	for _, c := range cities {
		if c.ID == targetCityID {
			cityName = c.Name
			break
		}
	}

	fmt.Printf("\n=== 2. Certified Asiacell Shops in %s (City ID: %d) ===\n", cityName, targetCityID)
	shops, err := client.GetCityShops(ctx, targetCityID)
	if err != nil {
		fmt.Printf("GetCityShops error: %v\n", err)
		return
	}

	fmt.Printf("Total Certified Branches Found: %d\n\n", len(shops))
	displayCount := 5
	if len(shops) < displayCount {
		displayCount = len(shops)
	}

	for i := 0; i < displayCount; i++ {
		s := shops[i]
		fmt.Printf("Branch [%d]: %s\n", i+1, s.Name)
		fmt.Printf("  Address:      %s\n", s.Address)
		fmt.Printf("  Phone:        %s\n", s.Phone)
		fmt.Printf("  Working Days: %s\n", s.WorkingDays)
		fmt.Printf("  Hours:        %s - %s\n", s.StartHour, s.FinishHour)
		fmt.Printf("  Coordinates:  Lat: %f, Lng: %f\n", s.Lat, s.Lng)
		fmt.Println()
	}
}
