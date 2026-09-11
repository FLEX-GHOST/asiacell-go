package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func main() {
	client, err := asiacell.NewClient(asiacell.WithTimeout(25 * time.Second))
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	sessionFile := "session.json"
	if err := client.LoadSessionFromFile(sessionFile); err != nil {
		fmt.Printf("Session file '%s' not found: %v\n", sessionFile, err)
		fmt.Println("Please run 'examples/01_auth_and_session' first to create a session.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	fmt.Println("===============================================================")
	fmt.Println(" Asiacell Advanced Features (Vanity, Gifting, Tickets, E-Cards)")
	fmt.Println("===============================================================")

	// 1. VIP / Vanity Numbers Marketplace
	fmt.Println("\n1. Browsing Vanity VIP Number Classes...")
	classes, err := client.GetVanityClasses(ctx)
	if err != nil {
		fmt.Printf("  Could not fetch vanity classes: %v\n", err)
	} else if classes != nil && len(classes.Data) > 0 {
		for _, c := range classes.Data {
			fmt.Printf("  • Tier: %-15s | Price: %s\n", c.Title, c.Price)
		}
	} else {
		fmt.Println("  No vanity tiers returned.")
	}

	fmt.Println("\n2. Searching Available VIP Numbers (Matching: 770)...")
	searchRes, err := client.SearchVanityNumbers(ctx, "770", "", 1, 5)
	if err != nil {
		fmt.Printf("  Search error: %v\n", err)
	} else if searchRes != nil && searchRes.Data != nil && len(searchRes.Data.List) > 0 {
		fmt.Printf("  Found %d matching numbers:\n", searchRes.Data.Total)
		for idx, item := range searchRes.Data.List {
			fmt.Printf("  [%02d] MSISDN: %-15s | Class: %-10s | Price: %s\n",
				idx+1, item.MSISDN, item.ClassName, item.Price)
		}
	} else {
		fmt.Println("  No matching VIP numbers currently listed.")
	}

	// 2. Compensation Checking
	fmt.Println("\n3. Checking Compensation Eligibility for this Line...")
	comp, err := client.CheckCompensation(ctx)
	if err != nil {
		fmt.Printf("  Compensation check error: %v\n", err)
	} else if comp != nil && len(comp.Data) > 0 {
		for _, item := range comp.Data {
			fmt.Printf("  • %s: %s (Eligible: %v | Claimed: %v)\n",
				item.Title, item.Benefit, item.Eligible, item.Claimed)
		}
	} else {
		fmt.Println("  No active compensation notices for this account.")
	}

	// 3. Support Resolution Center
	fmt.Println("\n4. Checking Open Support Tickets...")
	tickets, err := client.GetTickets(ctx)
	if err != nil {
		fmt.Printf("  Tickets query error: %v\n", err)
	} else if len(tickets) > 0 {
		for _, t := range tickets {
			fmt.Printf("  • Ticket #%s [%s]: %s (%s)\n",
				t.TicketNumber, t.Status, t.Subject, t.CreatedAt)
		}
	} else {
		fmt.Println("  No active or past support tickets found.")
	}

	// 4. Yooz MGM Referral Program
	fmt.Println("\n5. Checking Yooz MGM Youth Referral Code...")
	yooz, err := client.GetYoozMGM(ctx)
	if err != nil {
		fmt.Printf("  Yooz MGM query error: %v\n", err)
	} else if yooz != nil && yooz.Data != nil {
		fmt.Printf("  Referral Code: %s | Total Invites: %d | Total Earned: %s\n",
			yooz.Data.ReferralCode, yooz.Data.TotalInvites, yooz.Data.TotalEarned)
	} else {
		fmt.Println("  Yooz MGM not available for this tariff line.")
	}

	// 5. Digital E-Voucher Gift Cards (Gaming & Apps)
	fmt.Println("\n6. Browsing Digital E-Vouchers Catalog...")
	vouchers, err := client.GetEVoucherPackages(ctx)
	if err != nil {
		fmt.Printf("  E-Vouchers query error: %v\n", err)
	} else if vouchers != nil && len(vouchers.Data) > 0 {
		fmt.Printf("  Available E-Cards (%d items):\n", len(vouchers.Data))
		for idx, v := range vouchers.Data {
			if idx >= 5 {
				break
			}
			fmt.Printf("  [%02d] Card: %-25s | Category: %-15s | Price: %s\n",
				idx+1, v.Title, v.Category, v.Price)
		}
	} else {
		fmt.Println("  No e-voucher packages returned.")
	}

	fmt.Println("\n===============================================================")
	fmt.Println("All advanced endpoints successfully integrated into Asiacell Go SDK.")
	fmt.Println("===============================================================")
}
