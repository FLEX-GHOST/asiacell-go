package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func loadClient() (*asiacell.Client, string, error) {
	client, err := asiacell.NewClient(asiacell.WithTimeout(25 * time.Second))
	if err != nil {
		return nil, "", err
	}

	wallet := "07701234567"
	if err := client.LoadSessionFromFile("session.json"); err == nil {
		return client, wallet, nil
	}

	return client, wallet, nil
}

func main() {
	client, wallet, err := loadClient()
	if err != nil {
		fmt.Printf("Error initializing client: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========================================")
		fmt.Println("    Asiacell Go API - Interactive CLI   ")
		fmt.Println("========================================")
		fmt.Println("1. Account Overview & Profile")
		fmt.Println("2. Active Bundles")
		fmt.Println("3. Unlimited 4G Packages")
		fmt.Println("4. Special Offers & Addons")
		fmt.Println("5. Credit Transfer & Verify Transfer")
		fmt.Println("6. Recharge Voucher")
		fmt.Println("7. Services Management (Internet Protection & Codes)")
		fmt.Println("8. Branches & Cities Coverage")
		fmt.Println("9. Spin the Wheel & Rewards")
		fmt.Println("0. Exit")
		fmt.Println("----------------------------------------")
		fmt.Print("Choose an option [0-9]: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		switch choice {
		case "1":
			prof, err := client.GetProfile(ctx)
			if err != nil {
				fmt.Printf("GetProfile error: %v\n", err)
			} else {
				fmt.Printf("\n[Account Details]\nName:     %s\nPhone:    %s\nBalance:  %s\nValidity: %s\nData:     %s\nCalls:    %s\nSMS:      %s\n",
					prof.Name, prof.PhoneNumber, prof.Balance, prof.Validity, prof.RemainingData, prof.RemainingCalls, prof.RemainingSMS)
			}

		case "2":
			bundles, err := client.GetActiveBundles(ctx)
			if err != nil {
				fmt.Printf("GetActiveBundles error: %v\n", err)
			} else if len(bundles) == 0 {
				fmt.Println("\nNo active bundles found.")
			} else {
				fmt.Println("\n[Active Bundles]")
				for i, b := range bundles {
					fmt.Printf("[%d] %s: %.0f %s (Expires: %s)\n", i+1, b.Title, b.RemainingVolume, b.Unit, b.ExpireDate)
				}
			}

		case "3":
			fourG, err := client.GetUnlimited4GBundles(ctx)
			if err != nil {
				fmt.Printf("GetUnlimited4GBundles error: %v\n", err)
			} else {
				fmt.Println("\n[4G Unlimited Packages]")
				for i, b := range fourG {
					fmt.Printf("[%d] %s | Price: %s | Validity: %s\n", i+1, b.Title, b.Price, b.Validity)
				}
			}

		case "4":
			offers, err := client.GetSpecialOffers(ctx)
			if err != nil {
				fmt.Printf("GetSpecialOffers error: %v\n", err)
			} else {
				fmt.Println("\n[Special Offers]")
				for i, o := range offers {
					fmt.Printf("[%d] %s | Price: %s | Volume: %s | ID: %d\n", i+1, o.Title, o.Price, o.Volume, o.ID)
				}
			}

		case "5":
			fmt.Printf("\nMaster Wallet: %s\n", wallet)
			fmt.Println("1. Verify Incoming Transfer to Wallet")
			fmt.Println("2. Verify Outgoing Transfer")
			fmt.Print("Select [1-2]: ")
			subChoice, _ := reader.ReadString('\n')
			subChoice = strings.TrimSpace(subChoice)

			if subChoice == "1" {
				fmt.Print("Enter sender phone (e.g. 07701234567): ")
				sPhone, _ := reader.ReadString('\n')
				sPhone = strings.TrimSpace(sPhone)
				fmt.Print("Enter minimum amount (or 0): ")
				sAmtStr, _ := reader.ReadString('\n')
				sAmt, _ := strconv.ParseFloat(strings.TrimSpace(sAmtStr), 64)

				found, rec, err := client.VerifyIncomingTransfer(ctx, sPhone, sAmt)
				if err != nil {
					fmt.Printf("Verify error: %v\n", err)
				} else if found && rec != nil {
					fmt.Printf("Verified: YES | Sender: %s | Amount: %s | Date: %s\n", rec.MSISDN, rec.Amount, rec.CreatedAt)
				} else {
					fmt.Println("Transfer not found in recorded ledger.")
				}
			} else if subChoice == "2" {
				fmt.Print("Enter receiver phone (e.g. 07701234567): ")
				rPhone, _ := reader.ReadString('\n')
				rPhone = strings.TrimSpace(rPhone)
				fmt.Print("Enter minimum amount: ")
				rAmtStr, _ := reader.ReadString('\n')
				rAmt, _ := strconv.ParseFloat(strings.TrimSpace(rAmtStr), 64)

				found, rec, err := client.VerifyTransferTo(ctx, rPhone, rAmt)
				if err != nil {
					fmt.Printf("Verify error: %v\n", err)
				} else if found && rec != nil {
					fmt.Printf("Verified: YES | Receiver: %s | Amount: %s\n", rec.ReceiverMSISDN, rec.Amount)
				} else {
					fmt.Println("Transfer not found.")
				}
			}

		case "6":
			fmt.Print("Enter phone number to recharge: ")
			rcPhone, _ := reader.ReadString('\n')
			rcPhone = strings.TrimSpace(rcPhone)
			fmt.Print("Enter 14-digit voucher code: ")
			rcVoucher, _ := reader.ReadString('\n')
			rcVoucher = strings.TrimSpace(rcVoucher)

			if len(rcVoucher) >= 14 {
				resp, err := client.RechargeVoucher(ctx, rcPhone, rcVoucher, asiacell.RechargeTypeNormal)
				if err != nil {
					fmt.Printf("Recharge error: %v\n", err)
				} else {
					fmt.Printf("Result: Success=%v | Message: %s\n", resp.Success, resp.Message)
				}
			} else {
				fmt.Println("Invalid voucher code length.")
			}

		case "7":
			mgmt, err := client.GetServicesManagement(ctx)
			if err != nil {
				fmt.Printf("GetServicesManagement error: %v\n", err)
			} else {
				fmt.Printf("\n[Internet Protection: %s - %s]\n", mgmt.InternetControl.Title, mgmt.InternetControl.USSDCode)
				for i, g := range mgmt.CancellationGuides {
					fmt.Printf("[%d] %s: %s (%s)\n", i+1, g.Title, g.Code, g.Method)
				}
			}

		case "8":
			cities, err := client.GetCities(ctx)
			if err != nil {
				fmt.Printf("GetCities error: %v\n", err)
			} else {
				fmt.Println("\n[Covered Cities]")
				for i, c := range cities {
					fmt.Printf("[%2d] %-15s (ID: %d)\n", i+1, c.Name, c.ID)
				}
				fmt.Print("\nEnter City ID to view certified branches (e.g. 5 for Baghdad): ")
				cInput, _ := reader.ReadString('\n')
				cID, _ := strconv.Atoi(strings.TrimSpace(cInput))
				if cID > 0 {
					shops, err := client.GetCityShops(ctx, cID)
					if err != nil {
						fmt.Printf("GetCityShops error: %v\n", err)
					} else {
						fmt.Printf("\nFound %d branches:\n", len(shops))
						maxDisplay := 5
						if len(shops) < maxDisplay {
							maxDisplay = len(shops)
						}
						for i := 0; i < maxDisplay; i++ {
							s := shops[i]
							fmt.Printf("[%d] %s | %s | Tel: %s | Hours: %s - %s\n",
								i+1, s.Name, s.Address, s.Phone, s.StartHour, s.FinishHour)
						}
					}
				}
			}

		case "9":
			sw, err := client.GetSpinWheelStatus(ctx)
			if err != nil {
				fmt.Printf("GetSpinWheelStatus error: %v\n", err)
			} else {
				fmt.Printf("\nSpin Wheel: %s (Can Play: %v, Day: %d)\n", sw.Data.Title, sw.Data.IsPlayable, sw.Data.DayIndex)
			}
			shukran, err := client.GetShukranInfo(ctx)
			if err == nil && len(shukran.Headers) > 0 {
				fmt.Printf("Shukran: %s (%s)\n", shukran.Headers[0].Title, shukran.Headers[0].Credit)
			}

		case "0":
			fmt.Println("Exiting CLI. Goodbye!")
			cancel()
			return

		default:
			fmt.Println("Invalid choice. Please select [0-9].")
		}

		cancel()
	}
}
