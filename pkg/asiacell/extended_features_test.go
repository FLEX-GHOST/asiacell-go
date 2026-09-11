package asiacell

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFanZoneAndGamingFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/fanzone/home":
			_ = json.NewEncoder(w).Encode(FanZoneHomeResponse{
				Success: true,
				Data: &FanZoneHomeData{
					Title:         "Iraq Stars League",
					CompetitionID: "comp-101",
					Status:        "active",
				},
			})
		case "/api/v1/fanzone/kick-and-win/home":
			_ = json.NewEncoder(w).Encode(FanZoneKickAndWinResponse{
				Success: true,
				Data: &FanZoneKickAndWinData{
					AttemptsLeft: 3,
					MaxAttempts:  5,
				},
			})
		case "/api/v1/fanzone/kick-and-win/play-finish":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]bool{"success": true})
		case "/api/v1/fanzone/kick-and-win/reward":
			_ = json.NewEncoder(w).Encode(FanZoneRewardResponse{
				Success: true,
				Data: &FanZoneRewardData{
					RewardID: "rew-99",
					Title:    "5GB Free Data",
					Value:    "5000",
				},
			})
		case "/api/v1/fanzone/leader-board":
			_ = json.NewEncoder(w).Encode(FanZoneLeaderBoardResponse{
				Success: true,
				Data: &FanZoneLeaderBoardData{
					CompetitionID: "comp-101",
					UserRank:      12,
					UserScore:     850,
					Leaders: []FanZoneLeaderItem{
						{Rank: 1, Nickname: "BaghdadSniper", Score: 1200},
					},
				},
			})
		case "/api/v1/fanzone/predict-and-win":
			_ = json.NewEncoder(w).Encode(FanZonePredictResponse{
				Success: true,
				Data: &FanZonePredictData{
					CompetitionID: "comp-101",
					Matches: []FanZoneMatchPredictItem{
						{MatchID: "m-1", TeamA: "Al-Shorta", TeamB: "Al-Zawraa", ScoreA: 2, ScoreB: 1},
					},
				},
			})
		case "/api/v1/fanzone/grand-prizes":
			_ = json.NewEncoder(w).Encode(FanZoneGrandPrizesResponse{
				Success: true,
				Data: &FanZoneGrandPrizesData{
					Prizes: []FanZonePrizeItem{
						{ID: "p-1", Title: "Toyota Land Cruiser", Description: "Season grand prize"},
					},
				},
			})
		case "/api/v1/fanzone/rewards-history":
			_ = json.NewEncoder(w).Encode(FanZoneRewardsHistoryResponse{
				Success: true,
				Data: &FanZoneRewardsHistoryData{
					History: []FanZoneRewardData{
						{RewardID: "hist-1", Title: "100 Free Minutes"},
					},
				},
			})
		case "/api/v1/fanzone/onboarding/gen-nickname":
			_ = json.NewEncoder(w).Encode(NicknameResponse{
				Success: true,
				Data: &NicknameData{
					Nickname: "Falcon_964",
					Count:    1,
				},
			})
		case "/api/v1/fanzone/answer-and-win":
			_ = json.NewEncoder(w).Encode(FanZoneKickAndWinResponse{
				Success: true,
				Data: &FanZoneKickAndWinData{
					AttemptsLeft: 2,
				},
			})
		case "/api/v1/fanzone/favourite-team/pick":
			w.WriteHeader(http.StatusOK)
		case "/api/v1/fanzone/champion-team/pick":
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	home, err := client.GetFanZoneHome(ctx, "comp-101")
	if err != nil || home.CompetitionID != "comp-101" {
		t.Fatalf("GetFanZoneHome failed: %v", err)
	}

	kick, err := client.GetFanZoneKickAndWin(ctx, "comp-101")
	if err != nil || kick.AttemptsLeft != 3 {
		t.Fatalf("GetFanZoneKickAndWin failed: %v", err)
	}

	if err := client.FinishFanZoneKickAndWin(ctx, "comp-101", 100); err != nil {
		t.Fatalf("FinishFanZoneKickAndWin failed: %v", err)
	}

	rew, err := client.GetFanZoneKickAndWinReward(ctx, "comp-101", "tick-1")
	if err != nil || rew.RewardID != "rew-99" {
		t.Fatalf("GetFanZoneKickAndWinReward failed: %v", err)
	}

	lb, err := client.GetFanZoneLeaderBoard(ctx, "comp-101")
	if err != nil || lb.UserRank != 12 || len(lb.Leaders) != 1 {
		t.Fatalf("GetFanZoneLeaderBoard failed: %v", err)
	}

	pred, err := client.GetFanZonePredictions(ctx, "comp-101")
	if err != nil || len(pred.Matches) != 1 {
		t.Fatalf("GetFanZonePredictions failed: %v", err)
	}

	prizes, err := client.GetFanZoneGrandPrizes(ctx, "comp-101")
	if err != nil || len(prizes.Prizes) != 1 {
		t.Fatalf("GetFanZoneGrandPrizes failed: %v", err)
	}

	hist, err := client.GetFanZoneRewardsHistory(ctx, "comp-101")
	if err != nil || len(hist.History) != 1 {
		t.Fatalf("GetFanZoneRewardsHistory failed: %v", err)
	}

	nick, err := client.GenerateFanZoneNickname(ctx, "Falcon")
	if err != nil || nick.Nickname != "Falcon_964" {
		t.Fatalf("GenerateFanZoneNickname failed: %v", err)
	}

	ans, err := client.GetFanZoneAnswerAndWin(ctx, "comp-101")
	if err != nil || ans.AttemptsLeft != 2 {
		t.Fatalf("GetFanZoneAnswerAndWin failed: %v", err)
	}

	if err := client.PickFanZoneFavoriteTeam(ctx, "comp-101", "team-1"); err != nil {
		t.Fatalf("PickFanZoneFavoriteTeam failed: %v", err)
	}

	if err := client.PickFanZoneChampionTeam(ctx, "comp-101", "team-1"); err != nil {
		t.Fatalf("PickFanZoneChampionTeam failed: %v", err)
	}
}

func TestPartnersAndDiscountsFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/partners/categories":
			_ = json.NewEncoder(w).Encode(EOCategoryResponse{
				Success: true,
				Data: []EOCategory{
					{ID: 1, Name: "Restaurants", TotalPartner: 45},
				},
			})
		case "/api/v2/partners/cities":
			_ = json.NewEncoder(w).Encode(EOCityResponse{
				Success: true,
				Data: []EOCity{
					{ID: 1, Name: "Baghdad", TotalCategory: 10},
				},
			})
		case "/api/v2/partners/cities/1/categories":
			_ = json.NewEncoder(w).Encode(EOCategoryResponse{
				Success: true,
				Data: []EOCategory{
					{ID: 1, Name: "Restaurants", TotalPartner: 30},
				},
			})
		case "/api/v2/categories/1/cities/1/partners":
			_ = json.NewEncoder(w).Encode(EOPartnerResponse{
				Success: true,
				Data: []EOPartner{
					{ID: 101, Name: "Zaytouna Cafe", DiscountValue: 15.0, CityName: "Baghdad"},
				},
			})
		case "/api/v2/partners/register":
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	cats, err := client.GetPartnerCategories(ctx)
	if err != nil || len(cats) != 1 || cats[0].Name != "Restaurants" {
		t.Fatalf("GetPartnerCategories failed: %v", err)
	}

	cities, err := client.GetPartnerCities(ctx)
	if err != nil || len(cities) != 1 || cities[0].Name != "Baghdad" {
		t.Fatalf("GetPartnerCities failed: %v", err)
	}

	cityCats, err := client.GetPartnerCityCategories(ctx, 1)
	if err != nil || len(cityCats) != 1 {
		t.Fatalf("GetPartnerCityCategories failed: %v", err)
	}

	partners, err := client.GetPartnersByCategoryAndCity(ctx, 1, 1)
	if err != nil || len(partners) != 1 || partners[0].DiscountValue != 15.0 {
		t.Fatalf("GetPartnersByCategoryAndCity failed: %v", err)
	}

	err = client.RegisterPartner(ctx, PartnerRegisterRequest{
		Name:       "New Store",
		Phone:      "07700000000",
		Email:      "store@example.com",
		CityID:     1,
		CategoryID: 1,
	})
	if err != nil {
		t.Fatalf("RegisterPartner failed: %v", err)
	}
}

func TestYoozAvocadoFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v5/avocado/home":
			_ = json.NewEncoder(w).Encode(YoozHomeResponse{
				Success: true,
				Data: &YoozHomeData{
					Title:      "Yooz Club",
					Balance:    "15,000 IQD",
					InternetMB: 10240,
				},
			})
		case "/api/v3/avocado/bundles/screen":
			_ = json.NewEncoder(w).Encode(YoozAddOnListResponse{
				Success: true,
				Data: &YoozBundlesScreenData{
					GroupID: "grp-1",
					Title:   "Youth Bundles",
					Bundles: []YoozPlanEntity{
						{ID: 1, Title: "Yooz Daily", Price: "1000 IQD"},
					},
				},
			})
		case "/api/v3/avocado/bundles/classic-plans":
			_ = json.NewEncoder(w).Encode(YoozClassicPlansResponse{
				Success: true,
				Data: &YoozClassicPlansData{
					Plans: []YoozPlanEntity{
						{ID: 2, Title: "Classic Monthly", Price: "15000 IQD"},
					},
				},
			})
		case "/api/v3/avocado/bundles/omega-plans":
			_ = json.NewEncoder(w).Encode(YoozOmegaPlansResponse{
				Success: true,
				Data: &YoozOmegaPlansData{
					Voucher: "v-123",
					Plans: []YoozPlanEntity{
						{ID: 3, Title: "Omega Max", Price: "25000 IQD"},
					},
				},
			})
		case "/api/v2/avocado/bundles":
			_ = json.NewEncoder(w).Encode(YoozBundlesResponse{
				Success: true,
				Data: []YoozPlanEntity{
					{ID: 4, Title: "Yooz Social", Price: "5000 IQD"},
				},
			})
		case "/api/v2/avocado/data-cap":
			_ = json.NewEncoder(w).Encode(YoozDataCapResponse{
				Success: true,
				Data: &YoozDataCapData{
					CurrentLimitMB: 5120,
					MaxLimitMB:     20480,
					IsActive:       true,
				},
			})
		case "/api/v1/avocado/data-cap":
			w.WriteHeader(http.StatusOK)
		case "/api/v1/avocado/reward":
			_ = json.NewEncoder(w).Encode(YoozRewardResponse{
				Success: true,
				Data: &YoozRewardData{
					RewardTitle: "Yooz Points",
					Points:      500,
					Status:      "available",
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	home, err := client.GetYoozHome(ctx)
	if err != nil || home.InternetMB != 10240 {
		t.Fatalf("GetYoozHome failed: %v", err)
	}

	screen, err := client.GetYoozBundlesScreen(ctx, "grp-1")
	if err != nil || len(screen.Bundles) != 1 {
		t.Fatalf("GetYoozBundlesScreen failed: %v", err)
	}

	classic, err := client.GetYoozClassicPlans(ctx)
	if err != nil || len(classic.Plans) != 1 {
		t.Fatalf("GetYoozClassicPlans failed: %v", err)
	}

	omega, err := client.GetYoozOmegaPlans(ctx, "v-123", "07700000000")
	if err != nil || len(omega.Plans) != 1 {
		t.Fatalf("GetYoozOmegaPlans failed: %v", err)
	}

	bundles, err := client.GetYoozBundles(ctx)
	if err != nil || len(bundles) != 1 {
		t.Fatalf("GetYoozBundles failed: %v", err)
	}

	cap, err := client.GetYoozDataCap(ctx)
	if err != nil || cap.CurrentLimitMB != 5120 {
		t.Fatalf("GetYoozDataCap failed: %v", err)
	}

	if err := client.SetYoozDataCap(ctx, 10240); err != nil {
		t.Fatalf("SetYoozDataCap failed: %v", err)
	}

	rew, err := client.GetYoozReward(ctx)
	if err != nil || rew.Points != 500 {
		t.Fatalf("GetYoozReward failed: %v", err)
	}
}

func TestRechargeAndOnlinePaymentFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/recharge/screen1":
			_ = json.NewEncoder(w).Encode(RechargeNumberResponse{
				Success: true,
				Data: &RechargeNumberData{
					Title: "Select Number",
					Items: []RechargeNumberItem{
						{MSISDN: "07701234567", Selected: true},
					},
				},
			})
		case "/api/v1/recharge/screen2":
			_ = json.NewEncoder(w).Encode(RechargeTypeResponse{
				Success: true,
				Data: &RechargeTypeData{
					Title: "Select Recharge Type",
					Types: []RechargeTypeItem{
						{TypeID: "scratch_card", Title: "Prepaid Card"},
					},
				},
			})
		case "/api/v1/recharge/screen3":
			_ = json.NewEncoder(w).Encode(RechargeMethodResponse{
				Success: true,
				Data: &RechargeMethodData{
					Title: "Payment Gateways",
					OnlinePayments: []OnlinePaymentProvider{
						{ID: "zaincash", Name: "ZainCash"},
						{ID: "fib", Name: "First Iraqi Bank"},
					},
				},
			})
		case "/api/v1/recharge/screen4":
			_ = json.NewEncoder(w).Encode(OnlinePaymentResponse{
				Success: true,
				Data: &OnlinePaymentData{
					Title: "Payment Amount",
					PaymentPackages: []OnlinePaymentPackage{
						{ID: "pkg-10k", Price: 10000, Label: "10,000 IQD"},
					},
				},
			})
		case "/api/v1/recharge/confirmation":
			_ = json.NewEncoder(w).Encode(RechargeConfirmationResponse{
				Success: true,
				Data: &RechargeConfirmationData{
					TransactionID: "tx-777",
					Amount:        10000,
					Status:        "SUCCESS",
				},
			})
		case "/api/v1/top-up/payment-selection":
			_ = json.NewEncoder(w).Encode(PaymentSelectionResponse{
				Success: true,
				Data: &PaymentSelectionData{
					RechargeType: 1,
					ForOthers:    true,
					Methods: []OnlinePaymentProvider{
						{ID: "fastpay", Name: "FastPay"},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	nums, err := client.GetRechargeNumbers(ctx, "recharge")
	if err != nil || len(nums.Items) != 1 {
		t.Fatalf("GetRechargeNumbers failed: %v", err)
	}

	types, err := client.GetRechargeTypes(ctx, "recharge", "07701234567")
	if err != nil || len(types.Types) != 1 {
		t.Fatalf("GetRechargeTypes failed: %v", err)
	}

	methods, err := client.GetRechargeMethods(ctx, "recharge", "07701234567", "online")
	if err != nil || len(methods.OnlinePayments) != 2 {
		t.Fatalf("GetRechargeMethods failed: %v", err)
	}

	payment, err := client.GetOnlinePaymentDetails(ctx, "recharge", "07701234567", "online", "zaincash")
	if err != nil || len(payment.PaymentPackages) != 1 {
		t.Fatalf("GetOnlinePaymentDetails failed: %v", err)
	}

	confirm, err := client.GetRechargeConfirmation(ctx, "tx-777")
	if err != nil || confirm.Status != "SUCCESS" {
		t.Fatalf("GetRechargeConfirmation failed: %v", err)
	}

	selection, err := client.GetPaymentSelection(ctx)
	if err != nil || len(selection.Methods) != 1 {
		t.Fatalf("GetPaymentSelection failed: %v", err)
	}
}

func TestAddonsV3AndSurveysFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v3/home":
			_ = json.NewEncoder(w).Encode(HomeDashboardResponse{
				Success: true,
				Data: &HomeDashboardData{
					Title:      "Dashboard v3",
					Balance:    "20,000 IQD",
					InternetMB: 15360,
					QuickActions: []QuickActionItem{
						{ID: "act-1", Title: "Transfer"},
					},
				},
			})
		case "/api/v3/addon/cyo":
			_ = json.NewEncoder(w).Encode(CYOBundlesResponse{
				Success: true,
				Data: &CYOBundlesData{
					Title:       "Create Your Own",
					MinInternet: 1,
					MaxInternet: 100,
				},
			})
		case "/api/v1/offer/active-offers":
			_ = json.NewEncoder(w).Encode(ActiveOffersResponse{
				Success: true,
				Data: &ActiveOffersData{
					Offers: []ActiveOfferItem{
						{OfferID: "off-1", Title: "Double Data Bonus", Price: "5000 IQD"},
					},
				},
			})
		case "/api/v3/profile/manage-quick-actions":
			w.WriteHeader(http.StatusOK)
		case "/api/v1/app-feedback":
			w.WriteHeader(http.StatusOK)
		case "/api/v1/survey":
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(SurveyResponse{
					Success: true,
					Data: &SurveyData{
						SurveyID: "surv-1",
						Title:    "Network Satisfaction",
					},
				})
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case "/api/v1/voc":
			_ = json.NewEncoder(w).Encode(VoCResponse{
				Success: true,
				Data: &VoCEntity{
					Username: "user123",
					NMFloID:  "flow-99",
				},
			})
		case "/api/v1/promotions/video-tutorials":
			_ = json.NewEncoder(w).Encode(VideoTutorialsResponse{
				Success: true,
				Data: []VideoTutorialItem{
					{ID: "vid-1", Title: "How to Recharge", VideoURL: "https://video.asiacell.com/1"},
				},
			})
		case "/api/v1/one-yad":
			_ = json.NewEncoder(w).Encode(OneYadResponse{
				Success: true,
				Data: &OneYadData{
					Title:        "One Yad Initiative",
					TotalDonated: "50,000,000 IQD",
				},
			})
		case "/api/v1/one-yad/teams":
			_ = json.NewEncoder(w).Encode(OneYadTeamsResponse{
				Success: true,
				Data: []OneYadTeam{
					{TeamID: "team-baghdad", Name: "Baghdad Volunteers", City: "Baghdad"},
				},
			})
		case "/api/v1/one-yad/request":
			w.WriteHeader(http.StatusOK)
		case "/api/v1/epic":
			_ = json.NewEncoder(w).Encode(EpicLinesResponse{
				Success: true,
				Data: &EpicLinesData{
					Lines: []EpicLineItem{
						{MSISDN: "07709998877", PlanName: "Corporate Gold"},
					},
				},
			})
		case "/api/v1/epic/remaining/07709998877":
			_ = json.NewEncoder(w).Encode(EpicLineUsageResponse{
				Success: true,
				Data: &EpicLineUsageData{
					MSISDN:      "07709998877",
					RemainingMB: 8192,
					TotalMB:     10240,
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = ts.URL
	ctx := context.Background()

	dash, err := client.GetHomeDashboardV3(ctx, 33.3, 44.4, false)
	if err != nil || dash.InternetMB != 15360 || len(dash.QuickActions) != 1 {
		t.Fatalf("GetHomeDashboardV3 failed: %v", err)
	}

	cyo, err := client.GetCYOBundles(ctx, "grp-cyo")
	if err != nil || cyo.MaxInternet != 100 {
		t.Fatalf("GetCYOBundles failed: %v", err)
	}

	offers, err := client.GetActiveOffers(ctx)
	if err != nil || len(offers.Offers) != 1 {
		t.Fatalf("GetActiveOffers failed: %v", err)
	}

	if err := client.ManageQuickActions(ctx, []int{1, 2, 3}); err != nil {
		t.Fatalf("ManageQuickActions failed: %v", err)
	}

	if err := client.SubmitAppFeedback(ctx, "network", "Excellent coverage", 5); err != nil {
		t.Fatalf("SubmitAppFeedback failed: %v", err)
	}

	surv, err := client.GetSurveys(ctx)
	if err != nil || surv.SurveyID != "surv-1" {
		t.Fatalf("GetSurveys failed: %v", err)
	}

	if err := client.SubmitSurvey(ctx, "surv-1", map[string]string{"q1": "5"}); err != nil {
		t.Fatalf("SubmitSurvey failed: %v", err)
	}

	voc, err := client.GetVoiceOfCustomer(ctx)
	if err != nil || voc.NMFloID != "flow-99" {
		t.Fatalf("GetVoiceOfCustomer failed: %v", err)
	}

	vids, err := client.GetVideoTutorials(ctx)
	if err != nil || len(vids) != 1 {
		t.Fatalf("GetVideoTutorials failed: %v", err)
	}

	yad, err := client.GetOneYadHome(ctx)
	if err != nil || yad.Title != "One Yad Initiative" {
		t.Fatalf("GetOneYadHome failed: %v", err)
	}

	teams, err := client.GetOneYadTeams(ctx)
	if err != nil || len(teams) != 1 {
		t.Fatalf("GetOneYadTeams failed: %v", err)
	}

	if err := client.SubmitOneYadRequest(ctx, "team-baghdad", 5000); err != nil {
		t.Fatalf("SubmitOneYadRequest failed: %v", err)
	}

	epic, err := client.GetEpicLines(ctx)
	if err != nil || len(epic.Lines) != 1 {
		t.Fatalf("GetEpicLines failed: %v", err)
	}

	usage, err := client.GetEpicLineUsage(ctx, "07709998877")
	if err != nil || usage.RemainingMB != 8192 {
		t.Fatalf("GetEpicLineUsage failed: %v", err)
	}
}
