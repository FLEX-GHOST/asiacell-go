package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func loadClient() (*asiacell.Client, string, error) {
	client, err := asiacell.NewClient(asiacell.WithTimeout(20 * time.Second))
	if err != nil {
		return nil, "", err
	}

	if err := client.LoadSessionFromFile("session.json"); err != nil {
		return nil, "", fmt.Errorf("no saved session found (run 01_auth_and_session first)")
	}

	wallet := "07701234567"
	client.SetMasterWallet(wallet)
	return client, wallet, nil
}

func main() {
	client, wallet, err := loadClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== Asiacell Credit Transfer & Verification ===")
	fmt.Printf("Master Wallet: %s\n\n", wallet)

	fmt.Println("1. Recording mock incoming transfer in client ledger...")
	sampleSender := "07712345678"
	sampleAmount := 5000.0
	client.RecordIncomingTransfer(asiacell.TransactionRecord{
		MSISDN:    asiacell.FlexString(sampleSender),
		Amount:    asiacell.FlexString("5000"),
		CreatedAt: asiacell.FlexString(time.Now().Format("2006-01-02 15:04")),
	})
	fmt.Println("Recorded successfully.")

	fmt.Println("\n2. Verifying incoming transfer via VerifyIncomingTransfer...")
	found, rec, err := client.VerifyIncomingTransfer(ctx, sampleSender, sampleAmount)
	if err != nil {
		fmt.Printf("Verify error: %v\n", err)
	} else if found && rec != nil {
		fmt.Printf("Transfer Verified: True\n    Sender: %s | Amount: %s IQD | Date: %s\n",
			rec.MSISDN, rec.Amount, rec.CreatedAt)
	} else {
		fmt.Println("Transfer not found.")
	}

	fmt.Println("\n3. Testing non-existent sender...")
	foundFake, _, _ := client.VerifyIncomingTransfer(ctx, "07799999999", 5000)
	fmt.Printf("Fake sender verification result: %v (expected false)\n", foundFake)

	fmt.Println("\n4. Credit Transfer API Methods Available:")
	fmt.Println("  • client.StartCreditTransfer(ctx, receiverPhone, amount) -> returns PID, sends SMS OTP")
	fmt.Println("  • client.ConfirmCreditTransfer(ctx, pid, passcode)      -> confirms transfer with OTP")
	fmt.Println("  • client.TransferToWallet(ctx, amount)                  -> transfers directly to master wallet")
	fmt.Println("  • client.VerifyTransferTo(ctx, targetPhone, minAmount)  -> verifies outgoing transfer")
	fmt.Println("  • client.VerifyIncomingTransfer(ctx, sender, minAmount) -> verifies incoming transfer")
}
