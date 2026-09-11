package main

import (
	"context"
	"fmt"
	"os"
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
		fmt.Printf("Session file '%s' not found or invalid: %v\n", sessionFile, err)
		fmt.Println("Please run 'examples/01_auth_and_session' first to login and create a session.")
		return
	}

	wallet := "07722025730"
	if envWallet := os.Getenv("ASIACELL_MASTER_WALLET"); envWallet != "" {
		wallet = envWallet
	}
	client.SetMasterWallet(wallet)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	fmt.Println("===============================================================")
	fmt.Println(" Asiacell CDR Automated Incoming Transfer Verification (Go SDK)")
	fmt.Println(" التحقق الآلي من تحويلات الرصيد عبر سجل كشف الحساب (CDR)")
	fmt.Println("===============================================================")
	fmt.Printf("Master Wallet: %s\n\n", wallet)

	fmt.Println("1. Fetching live CDR transfer history from Asiacell ledger...")
	records, err := client.GetCDRTransferHistory(ctx, 1, 20)
	if err != nil {
		fmt.Printf("Warning: Failed to fetch CDR directly (%v)\n", err)
		fmt.Println("If CDR requires initial OTP activation, use client.SendCDROTP() and client.ConfirmCDROTP().")
	} else if len(records) > 0 {
		fmt.Printf("Found %d CDR records in ledger:\n", len(records))
		for idx, rec := range records {
			fmt.Printf("  [%02d] Phone: %-14s | Amount: %-10s | Date: %s | Title: %s\n",
				idx+1, rec.SubTitle, rec.Amount, rec.Description, rec.Title)
		}
	} else {
		fmt.Println("No CDR records returned from server.")
	}

	fmt.Println("\n2. Demonstrating automated payment verification (VerifyIncomingTransfer)...")
	testSender := "07744298878"
	expectedAmount := 1000.0

	fmt.Printf("Checking if '%s' transferred at least %.0f IQD...\n", testSender, expectedAmount)
	verified, matchedRec, err := client.VerifyIncomingTransfer(ctx, testSender, expectedAmount)
	if err != nil {
		fmt.Printf("Verification encountered error: %v\n", err)
	} else if verified && matchedRec != nil {
		fmt.Println(">> [VERIFIED SUCCESS / تم التأكيد بنجاح]")
		fmt.Printf("   Sender Phone : %s\n", matchedRec.MSISDN)
		fmt.Printf("   Amount Paid  : %s IQD\n", matchedRec.Amount)
		fmt.Printf("   Transfer Time: %s\n", matchedRec.CreatedAt)
	} else {
		fmt.Println(">> [NOT FOUND / لم يتم العثور على تحويل مطابق]")
		fmt.Printf("   No incoming transfer found from %s with minimum amount %.0f IQD.\n", testSender, expectedAmount)
	}

	fmt.Println("\n3. Testing non-existent transaction...")
	fakeVerified, _, _ := client.VerifyIncomingTransfer(ctx, "07700000000", 50000.0)
	fmt.Printf("   Fake sender check result: %v (Expected: false)\n", fakeVerified)

	fmt.Println("\n===============================================================")
	fmt.Println("Production bot architecture pattern:")
	fmt.Println("  1. Bot asks user to transfer credit to the master wallet.")
	fmt.Println("  2. User sends transfer and enters their phone number.")
	fmt.Println("  3. Bot calls client.VerifyIncomingTransfer(ctx, userPhone, amount).")
	fmt.Println("  4. If verified, order is approved instantly without any admin intervention!")
	fmt.Println("===============================================================")
}
