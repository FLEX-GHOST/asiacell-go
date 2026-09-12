package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func main() {
	// 1. Initialize client with storage and persistent DeviceID
	storage := asiacell.NewFileSessionStorage("session.json")
	client, err := asiacell.NewClient(
		asiacell.WithTimeout(20*time.Second),
		asiacell.WithLanguage("ar"),
		asiacell.WithSessionStorage(storage),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	fmt.Println("=== Asiacell Auth & Persistent Session Example ===")
	fmt.Printf("Device ID (Immutable UUID v4): %s\n\n", client.DeviceID())

	// Check if an existing session is available
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := client.LoadFromStorage(ctx, storage); err == nil && client.BiometricSecret() != "" {
		fmt.Println("Existing session loaded with Biometric Secret!")
		fmt.Println("Testing 3-Layer immortal authentication with live profile query...")

		profile, err := client.GetProfile(ctx)
		if err == nil {
			fmt.Printf("Session Active! Account Phone: %s | Balance: %s\n", profile.PhoneNumber, profile.Balance)
			fmt.Println("Session was restored without needing any SMS OTP!")
			return
		}
		fmt.Printf("Profile request failed, initiating re-login: %v\n", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter phone number (e.g. 07701234567): ")
	phoneInput, _ := reader.ReadString('\n')
	phone := strings.TrimSpace(phoneInput)
	if phone == "" {
		fmt.Println("Phone number is required")
		return
	}

	// 2. Request SMS OTP (with automatic Captcha resolution if requested by server)
	pid, err := client.Login(ctx, phone)
	if err != nil {
		if errors.Is(err, asiacell.ErrCaptchaRequired) {
			fmt.Println("Captcha required by server.")
			capt, captErr := client.GetCaptcha(ctx)
			if captErr != nil || capt == nil {
				fmt.Printf("Failed to get captcha: %v\n", captErr)
				return
			}
			fmt.Printf("Captcha Image URL: %s\n", capt.ImageURL())
			fmt.Print("Enter captcha solution: ")
			captInput, _ := reader.ReadString('\n')
			captCode := strings.TrimSpace(captInput)

			pid, err = client.LoginWithCaptcha(ctx, phone, captCode)
			if err != nil {
				fmt.Printf("Login with captcha failed: %v\n", err)
				return
			}
		} else {
			fmt.Printf("Login failed: %v\n", err)
			return
		}
	}

	fmt.Println("SMS OTP code sent to your phone.")
	fmt.Print("Enter the 6-digit SMS code: ")
	otpInput, _ := reader.ReadString('\n')
	otp := strings.TrimSpace(otpInput)

	// 3. Verify SMS OTP (automatically registers biometrics and saves BiometricSecret)
	resp, err := client.VerifySMS(ctx, pid, otp)
	if err != nil {
		fmt.Printf("SMS Verification failed: %v\n", err)
		return
	}

	fmt.Println("\nLogin successful!")
	fmt.Printf("User Phone     : %s (ID: %s)\n", resp.Username, resp.UserID)
	fmt.Printf("Device ID      : %s\n", client.DeviceID())
	fmt.Printf("Biometric Secret: %s\n", client.BiometricSecret())

	// 4. Export session and verify persistence
	exported, err := client.ExportSession()
	if err != nil {
		fmt.Printf("Failed to export session: %v\n", err)
		return
	}
	fmt.Println("\nSession exported successfully:")
	fmt.Printf("  • Access Token    : %s...\n", exported.AccessToken[:min(len(exported.AccessToken), 15)])
	fmt.Printf("  • Refresh Token   : %s...\n", exported.RefreshToken[:min(len(exported.RefreshToken), 15)])
	fmt.Printf("  • Biometric Secret: %s\n", exported.BiometricSecret)
	fmt.Printf("  • Device ID       : %s\n", exported.DeviceID)

	// 5. Test restoring session in a fresh client instance
	newClient, err := asiacell.NewClient()
	if err != nil {
		fmt.Printf("Failed to create fresh client: %v\n", err)
		return
	}
	if err := newClient.ImportSession(exported); err != nil {
		fmt.Printf("Failed to import session: %v\n", err)
		return
	}
	fmt.Println("\nTesting profile fetch on fresh client using imported session...")
	profile, err := newClient.GetProfile(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch profile: %v\n", err)
		return
	}
	fmt.Printf("Success! Account Phone: %s | Balance: %s\n", profile.PhoneNumber, profile.Balance)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
