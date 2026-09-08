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

	fmt.Println("=== 1. Services Management & Balance Protection ===")
	mgmt, err := client.GetServicesManagement(ctx)
	if err != nil {
		fmt.Printf("GetServicesManagement error: %v\n", err)
		return
	}

	fmt.Printf("Internet Protection: %s\n", mgmt.InternetControl.Title)
	fmt.Printf("USSD Code:           %s\n", mgmt.InternetControl.USSDCode)
	fmt.Printf("Customer Care:       %s\n", mgmt.InternetControl.CustomerCare)
	fmt.Printf("WhatsApp Support:    %s\n\n", mgmt.InternetControl.WhatsAppCare)

	fmt.Println("=== 2. Cancellation Guides & USSD Codes ===")
	for i, g := range mgmt.CancellationGuides {
		fmt.Printf("[%d] %s\n", i+1, g.Title)
		fmt.Printf("    ID:     %s\n", g.ID)
		fmt.Printf("    Code:   %s\n", g.Code)
		fmt.Printf("    Method: %s\n", g.Method)
		fmt.Printf("    Desc:   %s\n", g.Description)
		for _, a := range g.Actions {
			fmt.Printf("    Action: [%s] %s -> %s\n", a.Type, a.Title, a.Value)
		}
		fmt.Println()
	}

	fmt.Println("=== 3. Executing a Service Action ===")
	res, err := client.ExecuteServiceAction(ctx, "stop_payg_data", 0)
	if err != nil {
		fmt.Printf("ExecuteServiceAction error: %v\n", err)
	} else {
		fmt.Printf("Action Title:       %s\n", res.Title)
		fmt.Printf("Action Type:        %s\n", res.ActionType)
		fmt.Printf("Dial/Send Value:    %s\n", res.Value)
		fmt.Printf("Full Instruction:   %s\n", res.Instruction)
	}

	fmt.Println("\n=== 4. Digital Services & Platforms ===")
	digital, err := client.GetDigitalServices(ctx)
	if err != nil {
		fmt.Printf("GetDigitalServices error: %v\n", err)
	} else {
		for i, d := range digital.Services {
			fmt.Printf("[%d] %s\n    URL:  %s\n    Desc: %s\n", i+1, d.Title, d.URL, d.Description)
		}
	}
}
