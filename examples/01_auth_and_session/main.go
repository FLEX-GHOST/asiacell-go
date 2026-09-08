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
	client, err := asiacell.NewClient(
		asiacell.WithTimeout(20 * time.Second),
		asiacell.WithLanguage("ar"),
	)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Asiacell Auth & Session Example ===")
	fmt.Print("Enter phone number (e.g. 07701234567): ")
	phoneInput, _ := reader.ReadString('\n')
	phone := strings.TrimSpace(phoneInput)
	if phone == "" {
		fmt.Println("Phone number is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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

	resp, err := client.VerifySMS(ctx, pid, otp)
	if err != nil {
		fmt.Printf("SMS Verification failed: %v\n", err)
		return
	}

	fmt.Println("Login successful!")
	fmt.Printf("User: %s (ID: %s)\n", resp.Username, resp.UserID)

	err = client.SaveSessionToFile("session.json")
	if err != nil {
		fmt.Printf("Failed to save session: %v\n", err)
		return
	}
	fmt.Println("Session exported and saved to session.json")

	newClient, err := asiacell.NewClient()
	if err != nil {
		fmt.Printf("Failed to create new client: %v\n", err)
		return
	}
	err = newClient.LoadSessionFromFile("session.json")
	if err != nil {
		fmt.Printf("Failed to load session: %v\n", err)
		return
	}
	fmt.Println("Session loaded successfully into new client.")

	profile, err := newClient.GetProfile(ctx)
	if err != nil {
		fmt.Printf("Failed to fetch profile: %v\n", err)
		return
	}
	fmt.Printf("Account Phone: %s | Balance: %s\n", profile.PhoneNumber, profile.Balance)
}
