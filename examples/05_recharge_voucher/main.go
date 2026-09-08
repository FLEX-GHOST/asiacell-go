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

	fmt.Println("=== Asiacell Voucher Recharge Example ===")

	fmt.Println("1. Fetching Recharge History...")
	recharges, err := client.GetRechargeHistory(ctx)
	if err != nil {
		fmt.Printf("GetRechargeHistory error: %v\n", err)
	} else if len(recharges) == 0 {
		fmt.Println("No past recharge records found in current cache.")
	} else {
		for i, r := range recharges {
			fmt.Printf("[%d] Voucher: %s | Date: %s | Amount: %s\n",
				i+1, r.Voucher, r.CreatedAt, r.Amount)
		}
	}

	fmt.Println("\n2. To recharge a 14-digit voucher card, call:")
	fmt.Println("   resp, err := client.RechargeVoucher(ctx, phone, \"12345678901234\", asiacell.RechargeTypeNormal)")
	fmt.Println("   // resp.Success -> true/false")
	fmt.Println("   // resp.Message -> Asiacell network response message")
}
