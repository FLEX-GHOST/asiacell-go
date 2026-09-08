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

	fmt.Println("=== 1. Spin the Wheel Status ===")
	sw, err := client.GetSpinWheelStatus(ctx)
	if err != nil {
		fmt.Printf("GetSpinWheelStatus error: %v\n", err)
	} else {
		fmt.Printf("Title:          %s\n", sw.Data.Title)
		fmt.Printf("Subtitle:       %s\n", sw.Data.SubTitle)
		fmt.Printf("Can Play Now:   %v\n", sw.Data.IsPlayable)
		fmt.Printf("Day Index:      %d\n", sw.Data.DayIndex)
		fmt.Printf("Available Slides: %d\n", len(sw.Data.Slides))
	}

	fmt.Println("\n=== 2. Check Spin Wheel Eligibility ===")
	check, err := client.CheckSpinWheel(ctx)
	if err != nil {
		fmt.Printf("CheckSpinWheel error: %v\n", err)
	} else {
		fmt.Printf("Next Action:    %s\n", check.NextAction)
		fmt.Printf("Message:        %s\n", check.Message)
	}

	fmt.Println("\n=== 3. Shukran Service (Emergency Credit & Data) ===")
	shukran, err := client.GetShukranInfo(ctx)
	if err != nil {
		fmt.Printf("GetShukranInfo error: %v\n", err)
	} else {
		for i, h := range shukran.Headers {
			fmt.Printf("Header [%d]: Title: %s | Credit: %s | Data: %s\n",
				i+1, h.Title, h.Credit, h.Data)
		}
		fmt.Printf("Total Bodies: %d\n", len(shukran.Bodies))
	}

	fmt.Println("\n=== 4. Roaming Information ===")
	roaming, err := client.GetRoamingInfo(ctx)
	if err != nil {
		fmt.Printf("GetRoamingInfo error: %v\n", err)
	} else {
		fmt.Printf("Screen Title:   %s\n", roaming.ScreenTitle)
		fmt.Printf("Available Tags: %d\n", len(roaming.Tags))
		for i, t := range roaming.Tags {
			if i < 5 {
				fmt.Printf("  Tag [%d]: %s (%s)\n", i+1, t.Title, t.Tag)
			}
		}
	}

	fmt.Println("\n=== 5. Promotions & Rewards ===")
	promos, err := client.GetPromotions(ctx)
	if err != nil {
		fmt.Printf("GetPromotions error: %v\n", err)
	} else {
		fmt.Printf("Active Promotions: %d\n", len(promos))
		for i, p := range promos {
			fmt.Printf("  Promo [%d]: %s (ID: %d)\n    Action: %s (%s)\n",
				i+1, p.TitleDesc, p.ID, p.ActionButton.Title, p.ActionButton.Action)
		}
	}
}
