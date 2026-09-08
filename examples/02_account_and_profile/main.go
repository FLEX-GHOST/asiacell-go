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

	fmt.Println("=== 1. Account Overview ===")
	overview, err := client.GetProfile(ctx)
	if err != nil {
		fmt.Printf("GetProfile error: %v\n", err)
		return
	}

	fmt.Printf("Name:             %s\n", overview.Name)
	fmt.Printf("Phone Number:     %s\n", overview.PhoneNumber)
	fmt.Printf("Balance:          %s\n", overview.Balance)
	fmt.Printf("Validity:         %s\n", overview.Validity)
	fmt.Printf("Remaining Data:   %s\n", overview.RemainingData)
	fmt.Printf("Remaining Calls:  %s\n", overview.RemainingCalls)
	fmt.Printf("Remaining SMS:    %s\n", overview.RemainingSMS)

	fmt.Println("\n=== 2. Active Bundles ===")
	bundles, err := client.GetActiveBundles(ctx)
	if err != nil {
		fmt.Printf("GetActiveBundles error: %v\n", err)
	} else if len(bundles) == 0 {
		fmt.Println("No active bundles found.")
	} else {
		for i, b := range bundles {
			fmt.Printf("[%d] %s: %.0f %s (Expires: %s)\n", i+1, b.Title, b.RemainingVolume, b.Unit, b.ExpireDate)
		}
	}

	fmt.Println("\n=== 3. Profile Details ===")
	details, err := client.GetProfileDetails(ctx)
	if err != nil {
		fmt.Printf("GetProfileDetails error: %v\n", err)
	} else {
		fmt.Printf("Full Name: %s %s\n", details.FirstName, details.LastName)
		fmt.Printf("Phone:     %s\n", details.PhoneNumber)
		fmt.Printf("Email:     %s\n", details.Email)
		fmt.Printf("Photo:     %s\n", details.PhotoURL)
	}

	fmt.Println("\n=== 4. Notifications ===")
	notifs, err := client.GetNotifications(ctx)
	if err != nil {
		fmt.Printf("GetNotifications error: %v\n", err)
	} else if len(notifs) == 0 {
		fmt.Println("No notifications available.")
	} else {
		for i, n := range notifs {
			fmt.Printf("[%d] %s: %s (%s)\n", i+1, n.Title, n.Body, n.CreatedAt)
		}
	}
}
