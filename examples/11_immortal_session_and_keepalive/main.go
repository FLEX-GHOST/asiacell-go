package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func main() {
	fmt.Println("=======================================================================")
	fmt.Println(" Asiacell Immortal Session & Background Keep-Alive Daemon Example")
	fmt.Println(" الحفاظ الدائم على الجلسة ونبض الخلفية الدوري ضد انتهاء الصلاحية")
	fmt.Println("=======================================================================")

	storage := asiacell.NewFileSessionStorage("session.json")
	client, err := asiacell.NewClient(
		asiacell.WithTimeout(25*time.Second),
		asiacell.WithLanguage("ar"),
		asiacell.WithSessionStorage(storage),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Load session from disk
	if err := client.LoadFromStorage(ctx, storage); err != nil {
		fmt.Printf("No saved session found (%v). Please run 'examples/01_auth_and_session' first.\n", err)
		return
	}

	fmt.Println("\n1. Session Restored Successfully:")
	fmt.Printf("   • Phone Number       : %s\n", client.Phone())
	fmt.Printf("   • Device ID (UUID)   : %s\n", client.DeviceID())
	fmt.Printf("   • Biometric Secret   : %s\n", client.BiometricSecret())
	if exp := client.TokenExpiration(); !exp.IsZero() {
		fmt.Printf("   • Access Token Expiry: %s (in %v - proactive renewal triggers at < 2h)\n", exp.Format("2006-01-02 15:04:05"), time.Until(exp).Round(time.Minute))
	}

	// 2. Test 3-Layer Authentication:
	// If Access Token is expired, client automatically tries:
	// - Layer 2: POST /api/v1/validate (Refresh Token)
	// - Layer 3: POST /api/v1/biometrics/do-login (Biometric Silent Login)
	fmt.Println("\n2. Querying live user profile (verifying 3-layer auth resilience)...")
	profile, err := client.GetProfile(ctx)
	if err != nil {
		fmt.Printf("   Profile query error: %v\n", err)
		fmt.Println("   Attempting manual session refresh...")
		if refErr := client.RefreshSession(ctx); refErr != nil {
			fmt.Printf("   RefreshSession failed: %v\n", refErr)
			return
		}
		profile, err = client.GetProfile(ctx)
		if err != nil {
			fmt.Printf("   Re-attempt failed: %v\n", err)
			return
		}
	}
	fmt.Printf("   ✓ Profile active! Name: %s | Balance: %s\n", profile.Name, profile.Balance)

	// 3. Register CDR Expiration callback
	client.SetOnCDRExpired(func() {
		fmt.Println("\n[ALERT] CDR session has expired due to inactivity or server timeout!")
		fmt.Println("Triggering CDR re-activation via SMS OTP...")
		go activateCDR(client)
	})

	// 4. Start Background Keep-Alive Goroutine (Pulse every 15 minutes in prod, 10s for demo)
	keepAliveInterval := 10 * time.Second
	fmt.Printf("\n3. Starting Background Keep-Alive Goroutine (pulsing every %v)...\n", keepAliveInterval)
	fmt.Println("   • Pulse 1: Verifies main Profile (keeps auth token active)")
	fmt.Println("   • Pulse 2: Fetches CDR detail (keeps CDR ledger session warm & prevents inactivity timeout)")
	fmt.Println("   Press Ctrl+C to stop cleanly.")

	pulseChan := client.StartKeepAlive(ctx, keepAliveInterval)

	// Handle graceful shutdown on Ctrl+C (SIGINT / SIGTERM)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sigCh:
			fmt.Println("\nReceived shutdown signal, terminating keep-alive goroutine cleanly...")
			cancel()
			time.Sleep(200 * time.Millisecond)
			fmt.Println("Shutdown complete. Session data preserved.")
			return

		case pulse, ok := <-pulseChan:
			if !ok {
				fmt.Println("Keep-alive channel closed.")
				return
			}
			fmt.Printf("[%s] Pulse completed: ProfileActive=%v | CDRWarm=%v | CDRExpired=%v\n",
				pulse.Timestamp.Format("15:04:05"), pulse.ProfileOK, pulse.CDROK, pulse.CDRExpired)
			if pulse.Err != nil {
				fmt.Printf("    Notice: pulse detail: %v\n", pulse.Err)
			}
		}
	}
}

// activateCDR demonstrates how to activate CDR 2FA safely using PID extraction and GenericSMSConfirmationDTO
func activateCDR(client *asiacell.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("Requesting CDR activation OTP...")
	pid, err := client.SendCDROTP(ctx)
	if err != nil {
		fmt.Printf("SendCDROTP failed: %v\n", err)
		return
	}
	fmt.Printf("CDR OTP sent! Extracted PID: %s\n", pid)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the CDR SMS OTP code: ")
	codeStr, _ := reader.ReadString('\n')
	code := strings.TrimSpace(codeStr)

	// Confirm directly via POST /api/v1/cdr/confirm with GenericSMSConfirmationDTO {"PID": pid, "passcode": code}
	// Never calls /api/v1/smsvalidation to protect the OTP from burning!
	if err := client.ConfirmCDROTP(ctx, pid, code); err != nil {
		fmt.Printf("ConfirmCDROTP failed: %v\n", err)
		return
	}
	fmt.Println("✓ CDR Ledger successfully re-activated! Continuous monitoring resumed.")
}
