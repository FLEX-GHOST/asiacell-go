package asiacell

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)


func TestExtractPID(t *testing.T) {
	tests := []struct {
		url      string
		expected string
		hasErr   bool
	}{
		{
			url:      "https://www.asiacell.com/personal/my-account/login?PID=test-uuid-1234&foo=bar",
			expected: "test-uuid-1234",
			hasErr:   false,
		},
		{
			url:      "#/native/open/smsval/1?lang=ar&PID=8a444a4a-751e-4994-a686-468f87ab41c8",
			expected: "8a444a4a-751e-4994-a686-468f87ab41c8",
			hasErr:   false,
		},
		{
			url:      "/smsvalidation?PID=abcd-5678",
			expected: "abcd-5678",
			hasErr:   false,
		},
		{
			url:      "/smsvalidation?nopid=1",
			expected: "",
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		pid, err := extractPIDFromURL(tt.url)
		if tt.hasErr && err == nil {
			t.Fatalf("expected error for %s, got nil", tt.url)
		}
		if !tt.hasErr && err != nil {
			t.Fatalf("unexpected error for %s: %v", tt.url, err)
		}
		if pid != tt.expected {
			t.Fatalf("expected pid %q, got %q", tt.expected, pid)
		}
	}
}

func TestFlexString(t *testing.T) {
	type TestStruct struct {
		ID FlexString `json:"userId"`
	}

	tests := []struct {
		input    string
		expected string
	}{
		{input: `{"userId": 249019284}`, expected: "249019284"},
		{input: `{"userId": "987654321"}`, expected: "987654321"},
		{input: `{"userId": null}`, expected: ""},
		{input: `{}`, expected: ""},
	}

	for _, tt := range tests {
		var s TestStruct
		if err := json.Unmarshal([]byte(tt.input), &s); err != nil {
			t.Fatalf("failed to unmarshal %s: %v", tt.input, err)
		}
		if s.ID.String() != tt.expected {
			t.Fatalf("expected %q, got %q", tt.expected, s.ID.String())
		}
	}
}

func TestMockAuthAndProfileFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/login":
			var req LoginRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.Username == "fail" {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(LoginResponse{
					Success: false,
					Message: "Invalid phone number",
				})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(LoginResponse{
				Success:  true,
				NextURL:  "/verify?PID=mock-pid-123",
				Username: req.Username,
			})
		case "/api/v1/smsvalidation":
			var req SMSValidationRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			if req.Passcode != "123456" {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(SMSValidationResponse{
					Success: false,
					Message: "Invalid passcode",
				})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(SMSValidationResponse{
				Success:      true,
				AccessToken:  "mock-access-token",
				RefreshToken: "mock-refresh-token",
				UserID:       "mock-user-id",
				Username:     "07701234567",
			})
		case "/api/v1/profile":
			auth := r.Header.Get("Authorization")
			if auth != "Bearer mock-access-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"success": true,
				"data": {
					"bodies": [
						{"items": [{"name": "Ahmed Ali", "phoneNumber": "07701234567"}]},
						{},
						{"items": [{"value": 5000}]},
						{"items": [{"title": "5 GB"}]},
						{"items": [{"title": "100 Mins"}]},
						{"items": [{"title": "50 SMS"}]}
					]
				}
			}`))
		case "/api/v1/profile/view":
			auth := r.Header.Get("Authorization")
			if auth != "Bearer mock-access-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(ProfileDetailsResponse{
				Success: true,
				Data: ProfileDetails{
					FirstName:   "Ahmed",
					MiddleName:  "Ali",
					LastName:    "Hassan",
					PhoneNumber: "07701234567",
					Email:       "ahmed@example.com",
				},
			})
		case "/api/v1/top-up":
			var req RechargeRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(RechargeResponse{
				Success: true,
				Message: "Recharge successful",
			})
		case "/api/v1/credit-transfer/start":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(CreditTransferStartResponse{
				Success: true,
				PID:     "transfer-pid-999",
			})
		case "/api/v1/credit-transfer/do-transfer":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(TransferConfirmation{
				Success: true,
				Message: "Transfer complete",
			})
		case "/api/v1/transaction/recharge":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(RechargeHistoryResponse{
				Success: true,
				Data: []TransactionRecord{
					{Amount: "5000", CreatedAt: "2026-09-01", MSISDN: "07701234567"},
				},
			})
		case "/api/v1/spinwheel/", "/api/v2/spinwheel/ui":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(SpinWheelStatusResponse{
				Success: true,
				Data: SpinWheelData{
					RemainingTime: SpinWheelRemainingTime{
						IsPlayable: true,
						Hours:      0,
						Minutes:    0,
					},
				},
			})
		case "/api/v1/spinwheel/confirm", "/api/v2/spinwheel/confirm":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(SpinWheelPlayResponse{
				Success: true,
				Title:   "You won 1GB!",
			})
		case "/api/v1/addon":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(AddonResponse{
				Success: true,
				Data: struct {
					Bodies []AddonRawGroup `json:"bodies"`
				}{
					Bodies: []AddonRawGroup{
						{
							GroupID: 15,
							Title:   "عروض خاصة",
							Items: []AddonRawItem{
								{
									ID:       2406,
									Title:    "باقة RED 15,000",
									Price:    "15,000 دينار",
									Validity: "4 أسابيع",
									Volume:   "20 كيكابايت",
								},
							},
						},
					},
				},
			})
		case "/api/v1/addon/2406":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(AddonDetailResponse{
				Success: true,
				Data: AddonDetailData{
					Headers: []struct {
						Title           string `json:"title"`
						Validity        string `json:"validity"`
						Data            string `json:"data"`
						Minute          string `json:"minute"`
						SMS             string `json:"sms"`
						FeatureImage    string `json:"featureImage"`
						BackgroundImage string `json:"backgroundImage"`
						BackgroundColor string `json:"backgroundColor"`
					}{
						{
							Title:    "باقة RED 15,000",
							Validity: "4 أسابيع",
							Data:     "20 كيكابايت",
						},
					},
				},
			})
		case "/api/v2/home":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(HomeResponse{
				Data: HomeData{
					Headers: []HomeHeader{
						{
							BackgroundColor: "#FFFFFF",
							BackgroundImage: "https://app.asiacell.com/img/EOMainODPAR.jpg",
						},
					},
					Bodies: []HomeBody{
						{
							GroupID: 30,
							Title:   "باقات انترنت 4G بلاحدود",
							Items: []HomeItem{
								{
									ID:        2017,
									GroupID:   30,
									Title:     "باقة إنترنت 4G بلا حدود الیومیة",
									Unlimited: true,
									Price:     "3,000 دینار",
									Validity:  "1 يوم",
								},
							},
						},
					},
				},
			})
		case "/api/v1/partners/cities":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(CitiesResponse{
				Success: true,
				Data: []CityPartner{
					{
						ID:   4,
						Name: "بغداد",
					},
				},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = ts.URL

	_, err = client.Login(ctx, "fail")
	if err == nil {
		t.Fatalf("expected error for failed login")
	}

	pid, err := client.Login(ctx, "07701234567")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if pid != "mock-pid-123" {
		t.Fatalf("expected pid mock-pid-123, got: %s", pid)
	}

	_, err = client.VerifySMS(ctx, pid, "000000")
	if err != ErrInvalidPasscode {
		t.Fatalf("expected ErrInvalidPasscode, got: %v", err)
	}

	smsResp, err := client.VerifySMS(ctx, pid, "123456")
	if err != nil {
		t.Fatalf("VerifySMS failed: %v", err)
	}
	if smsResp.AccessToken != "mock-access-token" {
		t.Fatalf("expected access token mock-access-token, got: %s", smsResp.AccessToken)
	}

	overview, err := client.GetProfile(ctx)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if overview.Name != "Ahmed Ali" {
		t.Fatalf("expected name Ahmed Ali, got: %s", overview.Name)
	}
	if overview.Balance != "5000" {
		t.Fatalf("expected balance 5000, got: %s", overview.Balance)
	}
	if overview.RemainingData != "5 GB" {
		t.Fatalf("expected data 5 GB, got: %s", overview.RemainingData)
	}

	details, err := client.GetProfileDetails(ctx)
	if err != nil {
		t.Fatalf("GetProfileDetails failed: %v", err)
	}
	if details.FirstName != "Ahmed" {
		t.Fatalf("expected first name Ahmed, got: %s", details.FirstName)
	}

	recResp, err := client.RechargeVoucher(ctx, "07701234567", "12345678901234", RechargeTypeNormal)
	if err != nil {
		t.Fatalf("RechargeVoucher failed: %v", err)
	}
	if !recResp.Success {
		t.Fatalf("expected success in recharge")
	}

	transferPID, err := client.StartCreditTransfer(ctx, "07709876543", 2000)
	if err != nil {
		t.Fatalf("StartCreditTransfer failed: %v", err)
	}
	if transferPID != "transfer-pid-999" {
		t.Fatalf("expected transfer PID transfer-pid-999, got: %s", transferPID)
	}

	conf, err := client.ConfirmCreditTransfer(ctx, transferPID, "1234")
	if err != nil {
		t.Fatalf("ConfirmCreditTransfer failed: %v", err)
	}
	if !conf.Success {
		t.Fatalf("expected success in transfer confirmation")
	}

	rechargeHist, err := client.GetRechargeHistory(ctx)
	if err != nil {
		t.Fatalf("GetRechargeHistory failed: %v", err)
	}
	if len(rechargeHist) != 1 || rechargeHist[0].Amount != "5000" {
		t.Fatalf("unexpected recharge history: %+v", rechargeHist)
	}

	spinStatus, err := client.GetSpinWheelStatus(ctx)
	if err != nil {
		t.Fatalf("GetSpinWheelStatus failed: %v", err)
	}
	if !spinStatus.Data.RemainingTime.IsPlayable {
		t.Fatalf("expected spinwheel to be playable")
	}

	spinPlay, err := client.PlaySpinWheel(ctx)
	if err != nil {
		t.Fatalf("PlaySpinWheel failed: %v", err)
	}
	if spinPlay.Title != "You won 1GB!" {
		t.Fatalf("unexpected spin title: %s", spinPlay.Title)
	}

	specialOffers, err := client.GetSpecialOffers(ctx)
	if err != nil {
		t.Fatalf("GetSpecialOffers failed: %v", err)
	}
	if len(specialOffers) != 1 || specialOffers[0].ID != 2406 || specialOffers[0].USSDCode != "*299*1#" {
		t.Fatalf("unexpected special offers: %+v", specialOffers)
	}

	categories, err := client.GetAddonCategories(ctx)
	if err != nil {
		t.Fatalf("GetAddonCategories failed: %v", err)
	}
	if len(categories) != 1 || categories[0].Title != "عروض خاصة" {
		t.Fatalf("unexpected addon categories: %+v", categories)
	}

	detail, err := client.GetAddonDetail(ctx, 2406)
	if err != nil {
		t.Fatalf("GetAddonDetail failed: %v", err)
	}
	if len(detail.Headers) != 1 || detail.Headers[0].Title != "باقة RED 15,000" {
		t.Fatalf("unexpected addon detail: %+v", detail)
	}

	subResult, err := client.SubscribeSpecialOffer(ctx, 1)
	if err != nil {
		t.Fatalf("SubscribeSpecialOffer failed: %v", err)
	}
	if !subResult.Success || subResult.USSDCommand != "*299*1#" {
		t.Fatalf("unexpected subscription result: %+v", subResult)
	}

	cancelResult, err := client.CancelSpecialOffer(ctx)
	if err != nil {
		t.Fatalf("CancelSpecialOffer failed: %v", err)
	}
	if !cancelResult.Success || cancelResult.USSDCommand != "*299*0#" {
		t.Fatalf("unexpected cancel result: %+v", cancelResult)
	}

	servicesMgmt, err := client.GetServicesManagement(ctx)
	if err != nil {
		t.Fatalf("GetServicesManagement failed: %v", err)
	}
	if len(servicesMgmt.CancellationGuides) == 0 || servicesMgmt.InternetControl.USSDCode != "*223#" {
		t.Fatalf("unexpected services management result: %+v", servicesMgmt)
	}

	ctrl, err := client.GetServiceControl(ctx, "stop_payg_data")
	if err != nil {
		t.Fatalf("GetServiceControl failed: %v", err)
	}
	if ctrl.ID != "stop_payg_data" || len(ctrl.Actions) == 0 {
		t.Fatalf("unexpected service control result: %+v", ctrl)
	}

	actRes, err := client.ExecuteServiceAction(ctx, "stop_payg_data", 0)
	if err != nil {
		t.Fatalf("ExecuteServiceAction failed: %v", err)
	}
	if !actRes.Success || actRes.Value != "*223#" {
		t.Fatalf("unexpected service action result: %+v", actRes)
	}

	homeData, err := client.GetHomeLayout(ctx)
	if err != nil {
		t.Fatalf("GetHomeLayout failed: %v", err)
	}
	if len(homeData.Headers) != 1 || len(homeData.Bodies) != 1 {
		t.Fatalf("unexpected homeData: %+v", homeData)
	}

	fourGBundles, err := client.GetUnlimited4GBundles(ctx)
	if err != nil {
		t.Fatalf("GetUnlimited4GBundles failed: %v", err)
	}
	if len(fourGBundles) != 1 || fourGBundles[0].Price != "3,000 دینار" {
		t.Fatalf("unexpected 4G bundles: %+v", fourGBundles)
	}

	cities, err := client.GetCities(ctx)
	if err != nil {
		t.Fatalf("GetCities failed: %v", err)
	}
	if len(cities) != 1 || cities[0].Name != "بغداد" {
		t.Fatalf("unexpected cities: %+v", cities)
	}

	digSvc, err := client.GetDigitalServices(ctx)
	if err != nil {
		t.Fatalf("GetDigitalServices failed: %v", err)
	}
	if len(digSvc.Services) == 0 {
		t.Fatalf("unexpected digital services: %+v", digSvc)
	}



	sessionData, err := client.ExportSession()
	if err != nil {
		t.Fatalf("ExportSession failed: %v", err)
	}
	if sessionData.AccessToken != "mock-access-token" {
		t.Fatalf("expected exported token mock-access-token, got: %s", sessionData.AccessToken)
	}

	newClient, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create new client: %v", err)
	}
	if err := newClient.ImportSession(sessionData); err != nil {
		t.Fatalf("ImportSession failed: %v", err)
	}
	acc, _, _, _, _ := newClient.GetTokens()
	if acc != "mock-access-token" {
		t.Fatalf("expected imported token mock-access-token, got: %s", acc)
	}

	tmpDir := t.TempDir()
	sessPath := filepath.Join(tmpDir, "test_session.json")

	if err := client.SaveSessionToFile(sessPath); err != nil {
		t.Fatalf("SaveSessionToFile failed: %v", err)
	}

	if _, err := os.Stat(sessPath); err != nil {
		t.Fatalf("session file was not created: %v", err)
	}

	loadedClient, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := loadedClient.LoadSessionFromFile(sessPath); err != nil {
		t.Fatalf("LoadSessionFromFile failed: %v", err)
	}

	loadedAcc, _, _, _, _ := loadedClient.GetTokens()
	if loadedAcc != "mock-access-token" {
		t.Fatalf("expected loaded token mock-access-token, got: %s", loadedAcc)
	}
}

func TestMockCityShopsAndAddonSubscription(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/v1/shops"):
			if r.URL.Query().Get("cities") == "5" {
				_ = json.NewEncoder(w).Encode(ShopsResponse{
					Success: true,
					Message: "success",
					Data: []ShopInfo{
						{
							ID:       190,
							Name:     "بغداد",
							Address:  "شارع فلسطين",
							Phone:    "7728886200",
							CityName: "بغداد",
							CityID:   5,
						},
					},
				})
				return
			}
			_ = json.NewEncoder(w).Encode(ShopsResponse{
				Success: true,
				Data:    []ShopInfo{},
			})
		case strings.HasPrefix(r.URL.Path, "/api/v2/addon/summary/2017"):
			_ = json.NewEncoder(w).Encode(AddonSummaryResponse{
				Success: true,
				Message: "Success",
				Data: AddonSummaryData{
					ID: 2017,
					Options: []AddonSummaryOption{
						{
							Key:      "ONE_MONTH",
							Title:    "باقة إنترنت 4G بلا حدود الیومیة",
							Price:    "3,000 دینار",
							Validity: "1 يوم",
							Total:    3000,
						},
					},
				},
			})
		case strings.HasPrefix(r.URL.Path, "/api/v1/addon"):
			if r.Method == http.MethodPost {
				if r.URL.Query().Get("addOnId") == "2017" {
					_ = json.NewEncoder(w).Encode(AddonSubscribeResponse{
						Success: true,
						Message: "تم الاشتراك بنجاح",
					})
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(AddonSubscribeResponse{
					Success: false,
					Message: "فشل الاشتراك",
				})
				return
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	shops, err := client.GetCityShops(ctx, 5)
	if err != nil {
		t.Fatalf("GetCityShops failed: %v", err)
	}
	if len(shops) != 1 || shops[0].ID != 190 {
		t.Fatalf("unexpected shops result: %+v", shops)
	}

	summary, err := client.GetAddonSummary(ctx, 2017)
	if err != nil {
		t.Fatalf("GetAddonSummary failed: %v", err)
	}
	if summary.ID != 2017 || len(summary.Options) != 1 {
		t.Fatalf("unexpected summary result: %+v", summary)
	}

	sub, err := client.SubscribeAddon(ctx, 2017)
	if err != nil {
		t.Fatalf("SubscribeAddon failed: %v", err)
	}
	if !sub.Success {
		t.Fatalf("expected success subscription, got: %+v", sub)
	}
}

func TestVerifyIncomingTransfer(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	found, _, err := client.VerifyIncomingTransfer(ctx, "07701234567", 5000)
	if err != nil {
		t.Fatalf("VerifyIncomingTransfer unexpected error: %v", err)
	}
	if found {
		t.Fatalf("expected not found initially")
	}

	client.RecordIncomingTransfer(TransactionRecord{
		MSISDN:    "07701234567",
		Amount:    "5000",
		CreatedAt: "2026-09-08 10:00",
	})

	recorded := client.RecordedIncomingTransfers()
	if len(recorded) != 1 {
		t.Fatalf("expected 1 recorded incoming transfer, got %d", len(recorded))
	}

	found, rec, err := client.VerifyIncomingTransfer(ctx, "7701234567", 5000)
	if err != nil {
		t.Fatalf("VerifyIncomingTransfer failed: %v", err)
	}
	if !found || rec == nil || string(rec.Amount) != "5000" {
		t.Fatalf("expected verified incoming transfer, got found=%v, rec=%+v", found, rec)
	}

	foundDiffAmt, _, _ := client.VerifyIncomingTransfer(ctx, "07701234567", 10000)
	if foundDiffAmt {
		t.Fatalf("expected not found for different amount")
	}
}

func TestDeviceIDRotation(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	initialDevID := client.DeviceID()
	if initialDevID == "" {
		t.Fatalf("expected non-empty initial device ID")
	}

	newDevID, err := client.RotateDeviceID()
	if err != nil {
		t.Fatalf("RotateDeviceID failed: %v", err)
	}
	if newDevID == "" || newDevID == initialDevID {
		t.Fatalf("expected new device ID different from initial, got %s vs %s", newDevID, initialDevID)
	}
	if client.DeviceID() != newDevID {
		t.Fatalf("client.DeviceID() did not update: got %s, expected %s", client.DeviceID(), newDevID)
	}
}

func TestCaptchaDataImageURL(t *testing.T) {
	var nilData *CaptchaData
	if nilData.ImageURL() != "" {
		t.Fatalf("expected empty string for nil CaptchaData")
	}

	c1 := &CaptchaData{CaptchaURL: "https://example.com/c1.png"}
	if c1.ImageURL() != "https://example.com/c1.png" {
		t.Fatalf("unexpected image URL: %s", c1.ImageURL())
	}

	c2 := &CaptchaData{OriginSource: "https://example.com", ResourceURL: "/c2.png"}
	if c2.ImageURL() != "https://example.com/c2.png" {
		t.Fatalf("unexpected image URL: %s", c2.ImageURL())
	}

	c3 := &CaptchaData{CaptchaImage: "data:image/png;base64,ABC"}
	if c3.ImageURL() != "data:image/png;base64,ABC" {
		t.Fatalf("unexpected image URL: %s", c3.ImageURL())
	}
}

func TestCaptchaLoginRetry(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/login" {
			attempts++
			if attempts == 1 {
				_ = json.NewEncoder(w).Encode(LoginResponse{
					Success:        false,
					RequireCaptcha: true,
					Message:        "Captcha required",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(LoginResponse{
				Success: true,
				NextURL: "/validate?PID=pid_success_after_rotate",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = ts.URL

	ctx := context.Background()
	pid, err := client.Login(ctx, "07701234567")
	if err != nil {
		t.Fatalf("Login failed after rotate: %v", err)
	}
	if pid != "pid_success_after_rotate" {
		t.Fatalf("unexpected PID: %s", pid)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 login attempts, got %d", attempts)
	}
}




func TestVanityAndGifting(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/vanity/classes":
			_ = json.NewEncoder(w).Encode(VanityClassesResponse{
				Success: true,
				Data: []VanityClass{
					{ID: "1", Title: "Gold VIP", Price: "50000 IQD"},
					{ID: "2", Title: "Silver VIP", Price: "25000 IQD"},
				},
			})
		case "/api/v2/vanity":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(VanitySearchResponse{
					Success: true,
					Data: &VanitySearchData{
						Total: 1,
						List: []VanityNumberItem{
							{MSISDN: "07700001111", ClassName: "Gold VIP", Price: "50000"},
						},
					},
				})
			} else if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(ReserveVanityResponse{
					Success: true,
					PID:     "pid_reserve_1234",
				})
			}
		case "/api/v2/vanity/07700001111/detail":
			_ = json.NewEncoder(w).Encode(VanityDetailResponse{
				Success: true,
				Data: &VanityNumberItem{
					MSISDN: "07700001111",
					Price:  "50000",
				},
			})
		case "/api/v1/addon/send-as-gift":
			_ = json.NewEncoder(w).Encode(SendGiftResponse{
				Success: true,
				Message: "gift sent successfully",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	classes, err := client.GetVanityClasses(ctx)
	if err != nil || len(classes.Data) != 2 {
		t.Fatalf("GetVanityClasses failed: %v, len: %d", err, len(classes.Data))
	}

	searchRes, err := client.SearchVanityNumbers(ctx, "0770000", "1", 1, 10)
	if err != nil || searchRes.Data.Total != 1 {
		t.Fatalf("SearchVanityNumbers failed: %v", err)
	}

	detail, err := client.GetVanityDetail(ctx, "07700001111")
	if err != nil || string(detail.Data.Price) != "50000" {
		t.Fatalf("GetVanityDetail failed: %v", err)
	}

	reserve, err := client.ReserveVanityNumber(ctx, "07700001111", "1")
	if err != nil || string(reserve.PID) != "pid_reserve_1234" {
		t.Fatalf("ReserveVanityNumber failed: %v", err)
	}

	if err := client.SendGiftAddon(ctx, 42, "07711223344"); err != nil {
		t.Fatalf("SendGiftAddon failed: %v", err)
	}
}

func TestTicketsAndCompensation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/resolution-center/categories":
			_ = json.NewEncoder(w).Encode(TicketCategoriesResponse{
				Success: true,
				Data: []TicketCategory{
					{ID: "cat_1", Title: "Network Problem"},
				},
			})
		case "/api/v1/resolution-center":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(TicketsResponse{
					Success: true,
					Data: []TicketItem{
						{TicketNumber: "T-9988", Status: "OPEN", Subject: "Slow 4G"},
					},
				})
			} else if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(SubmitTicketResponse{
					Success:      true,
					TicketNumber: "T-9988",
				})
			}
		case "/api/v1/compensation":
			_ = json.NewEncoder(w).Encode(CompensationResponse{
				Success: true,
				Data: []CompensationItem{
					{ID: "comp_1", Title: "Maintenance Outage Refund", Benefit: "5GB Free Data", Eligible: true},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	cats, err := client.GetTicketCategories(ctx)
	if err != nil || len(cats) != 1 {
		t.Fatalf("GetTicketCategories failed: %v", err)
	}

	tickets, err := client.GetTickets(ctx)
	if err != nil || len(tickets) != 1 || tickets[0].TicketNumber != "T-9988" {
		t.Fatalf("GetTickets failed: %v", err)
	}

	created, err := client.CreateTicket(ctx, "cat_1", "Slow 4G in Baghdad")
	if err != nil || string(created.TicketNumber) != "T-9988" {
		t.Fatalf("CreateTicket failed: %v", err)
	}

	comp, err := client.CheckCompensation(ctx)
	if err != nil || len(comp.Data) != 1 || !comp.Data[0].Eligible {
		t.Fatalf("CheckCompensation failed: %v", err)
	}
}

func TestYoozAndEVouchers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/yooz-mgm":
			_ = json.NewEncoder(w).Encode(YoozMGMResponse{
				Success: true,
				Data: &YoozMGMData{
					ReferralCode: "YOOZ123",
					TotalInvites: 5,
				},
			})
		case "/api/v1/yooz-mgm/apply-code":
			_ = json.NewEncoder(w).Encode(ApplyPromoResponse{
				Success: true,
			})
		case "/api/v2/e-voucher/packages":
			_ = json.NewEncoder(w).Encode(EVoucherPackagesResponse{
				Success: true,
				Data: []EVoucherPackageItem{
					{ID: 1, Title: "PUBG 60 UC", Price: "1500 IQD"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	mgm, err := client.GetYoozMGM(ctx)
	if err != nil || mgm.Data.ReferralCode != "YOOZ123" {
		t.Fatalf("GetYoozMGM failed: %v", err)
	}

	if err := client.ApplyYoozMGMCode(ctx, "PROMO99"); err != nil {
		t.Fatalf("ApplyYoozMGMCode failed: %v", err)
	}

	vouchers, err := client.GetEVoucherPackages(ctx)
	if err != nil || len(vouchers.Data) != 1 || string(vouchers.Data[0].Title) != "PUBG 60 UC" {
		t.Fatalf("GetEVoucherPackages failed: %v", err)
	}
}
