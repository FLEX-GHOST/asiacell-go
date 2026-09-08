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

	fmt.Println("=== 1. Unlimited 4G Packages ===")
	bundles, err := client.GetUnlimited4GBundles(ctx)
	if err != nil {
		fmt.Printf("GetUnlimited4GBundles error: %v\n", err)
	} else {
		for i, b := range bundles {
			fmt.Printf("[%d] %s\n    Price: %s | Validity: %s | Volume: %s\n",
				i+1, b.Title, b.Price, b.Validity, b.Volume)
		}
	}

	fmt.Println("\n=== 2. Special Offers ===")
	offers, err := client.GetSpecialOffers(ctx)
	if err != nil {
		fmt.Printf("GetSpecialOffers error: %v\n", err)
	} else {
		for i, o := range offers {
			fmt.Printf("[%d] %s\n    Price: %s | Volume: %s | Validity: %s | ID: %d\n",
				i+1, o.Title, o.Price, o.Volume, o.Validity, o.ID)
		}
	}

	fmt.Println("\n=== 3. Addon Categories ===")
	cats, err := client.GetAddonCategories(ctx)
	if err != nil {
		fmt.Printf("GetAddonCategories error: %v\n", err)
	} else {
		for i, c := range cats {
			fmt.Printf("[%d] %s (GroupID: %d, Count: %d)\n", i+1, c.Title, c.GroupID, len(c.Items))
		}
	}

	fmt.Println("\n=== 4. Addon Package Summary via API ===")
	targetID := 2017
	summary, err := client.GetAddonSummary(ctx, targetID)
	if err != nil {
		fmt.Printf("GetAddonSummary error: %v\n", err)
	} else {
		fmt.Printf("Addon ID: %d\n", summary.ID)
		for i, opt := range summary.Options {
			fmt.Printf("Option [%d]: Title: %s | Price: %s | Validity: %s | Volume: %s\n",
				i+1, opt.Title, opt.Price, opt.Validity, opt.Volumed)
		}
		fmt.Println("\nTo subscribe directly via API, you call:")
		fmt.Printf("client.SubscribeAddon(ctx, %d)\n", targetID)
	}
}
