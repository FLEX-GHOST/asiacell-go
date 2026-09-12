package asiacell

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestPersistentImmutableDeviceID(t *testing.T) {
	var receivedDeviceID1, receivedXDeviceID1 string
	var receivedDeviceID2, receivedXDeviceID2 string
	var reqCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&reqCount, 1)
		if count == 1 {
			receivedDeviceID1 = r.Header.Get("DeviceId")
			receivedXDeviceID1 = r.Header.Get("x-device-id")
		} else {
			receivedDeviceID2 = r.Header.Get("DeviceId")
			receivedXDeviceID2 = r.Header.Get("x-device-id")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data":    map[string]any{},
		})
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient(
		WithTokens("mock-token", "mock-refresh"),
		WithDeviceID("custom-uuid-1234-5678"),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	if client.DeviceID() != "custom-uuid-1234-5678" {
		t.Fatalf("expected device id custom-uuid-1234-5678, got %s", client.DeviceID())
	}

	// First request
	_, _ = client.doRequest(ctx, http.MethodGet, "/api/v1/test1", nil)

	// Second request
	_, _ = client.doRequest(ctx, http.MethodGet, "/api/v1/test2", nil)

	if receivedDeviceID1 != "custom-uuid-1234-5678" || receivedXDeviceID1 != "custom-uuid-1234-5678" {
		t.Errorf("req 1 headers mismatch: DeviceId=%q, x-device-id=%q", receivedDeviceID1, receivedXDeviceID1)
	}
	if receivedDeviceID2 != "custom-uuid-1234-5678" || receivedXDeviceID2 != "custom-uuid-1234-5678" {
		t.Errorf("req 2 headers mismatch: DeviceId=%q, x-device-id=%q", receivedDeviceID2, receivedXDeviceID2)
	}
}

func TestThreeLayerAuth_AutoRefresh(t *testing.T) {
	var profileCalls int32
	var validateCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/profile":
			atomic.AddInt32(&profileCalls, 1)
			auth := r.Header.Get("Authorization")
			if auth != "Bearer token-renewed-v2" {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"message": "Token expired",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data": map[string]any{
					"bodies": []any{},
				},
			})

		case "/api/v1/validate":
			atomic.AddInt32(&validateCalls, 1)
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["refreshToken"] != "Bearer refresh-v1" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(SMSValidationResponse{
				Success:      true,
				AccessToken:  "token-renewed-v2",
				RefreshToken: "refresh-renewed-v2",
				Secret:       "enc-secret-from-validate",
				Username:     "07701234567",
			})
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient(
		WithTokens("token-stale-v1", "refresh-v1"),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL
	client.SetPhone("07701234567")

	_, err = client.GetProfile(ctx)
	if err != nil {
		t.Fatalf("GetProfile should have succeeded after auto-refresh: %v", err)
	}

	if atomic.LoadInt32(&validateCalls) != 1 {
		t.Errorf("expected 1 validate call, got %d", validateCalls)
	}
	if atomic.LoadInt32(&profileCalls) != 2 {
		t.Errorf("expected 2 profile calls (1 401 + 1 retry), got %d", profileCalls)
	}
	if client.BiometricSecret() != "enc-secret-from-validate" {
		t.Errorf("expected secret enc-secret-from-validate, got %s", client.BiometricSecret())
	}
}

func TestThreeLayerAuth_BiometricFallback(t *testing.T) {
	var validateCalls int32
	var bioLoginCalls int32
	var profileCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/profile":
			atomic.AddInt32(&profileCalls, 1)
			auth := r.Header.Get("Authorization")
			if auth != "Bearer token-renewed-by-biometrics" {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"message": "Token expired",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data": map[string]any{
					"bodies": []any{},
				},
			})

		case "/api/v1/validate":
			atomic.AddInt32(&validateCalls, 1)
			// Layer 2 Refresh Token fails (e.g. 401 expired)
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": false,
				"message": "Refresh token expired",
			})

		case "/api/v1/biometrics/do-login":
			atomic.AddInt32(&bioLoginCalls, 1)
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)

			if body["msisdn"] != "07701234567" || body["secret"] != "mock-bio-secret" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			_ = json.NewEncoder(w).Encode(map[string]any{
				"success":        true,
				"access_token":   "token-renewed-by-biometrics",
				"refresh_token":  "new-refresh-v3",
				"secret":         "mock-bio-secret-rotated",
				"username":       "07701234567",
			})
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient(
		WithTokens("stale-access", "stale-refresh"),
		WithPhone("07701234567"),
		WithBiometricSecret("mock-bio-secret"),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	_, err = client.GetProfile(ctx)
	if err != nil {
		t.Fatalf("GetProfile should succeed via Layer 3 Biometric silent login: %v", err)
	}

	if atomic.LoadInt32(&validateCalls) != 1 {
		t.Errorf("expected 1 validate call, got %d", validateCalls)
	}
	if atomic.LoadInt32(&bioLoginCalls) != 1 {
		t.Errorf("expected 1 biometric login call, got %d", bioLoginCalls)
	}
	if atomic.LoadInt32(&profileCalls) != 2 {
		t.Errorf("expected 2 profile calls, got %d", profileCalls)
	}
	if client.BiometricSecret() != "mock-bio-secret-rotated" {
		t.Errorf("expected updated biometric secret, got %s", client.BiometricSecret())
	}
}

func TestFirstLogin_BiometricsRegistration(t *testing.T) {
	var bioRegCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/smsvalidation":
			_ = json.NewEncoder(w).Encode(SMSValidationResponse{
				Success:      true,
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				UserID:       "usr-100",
				Username:     "07709998877",
			})

		case "/api/v1/biometrics/register":
			atomic.AddInt32(&bioRegCalls, 1)
			_ = json.NewEncoder(w).Encode(BiometricRegisterResponse{
				Success: true,
				Secret:  "new-biometric-secret-xyz",
			})
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	resp, err := client.VerifySMS(ctx, "test-pid", "123456")
	if err != nil {
		t.Fatalf("VerifySMS failed: %v", err)
	}
	if resp.AccessToken != "new-access-token" {
		t.Errorf("unexpected access token: %s", resp.AccessToken)
	}

	if atomic.LoadInt32(&bioRegCalls) != 1 {
		t.Errorf("expected 1 biometrics register call upon first login, got %d", bioRegCalls)
	}
	if client.BiometricSecret() != "new-biometric-secret-xyz" {
		t.Errorf("expected biometric secret new-biometric-secret-xyz, got %s", client.BiometricSecret())
	}

	exported, err := client.ExportSession()
	if err != nil {
		t.Fatalf("ExportSession failed: %v", err)
	}
	if exported.BiometricSecret != "new-biometric-secret-xyz" {
		t.Errorf("expected exported session to have biometric secret, got %s", exported.BiometricSecret)
	}
	if exported.Phone != "07709998877" {
		t.Errorf("expected exported phone 07709998877, got %s", exported.Phone)
	}
}

func TestCDR_SendOTP_ExtractPID_Confirm(t *testing.T) {
	var sendCalls int32
	var confirmCalls int32
	var receivedPID, receivedPasscode string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/cdr/send-otp":
			atomic.AddInt32(&sendCalls, 1)
			_ = json.NewEncoder(w).Encode(CDROTPResponse{
				Success: true,
				NextURL: "/verify?PID=cdr-mock-uuid-555&type=otp",
			})

		case "/api/v1/cdr/confirm":
			atomic.AddInt32(&confirmCalls, 1)
			var dto GenericSMSConfirmationDTO
			_ = json.NewDecoder(r.Body).Decode(&dto)
			receivedPID = dto.PID
			receivedPasscode = dto.Passcode

			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"message": "CDR verified successfully",
			})

		case "/api/v1/smsvalidation":
			t.Errorf("CRITICAL VIOLATION: /api/v1/smsvalidation was called during CDR confirmation!")
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := NewClient(WithTokens("mock-token", "mock-refresh"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	// 1. Send CDR OTP and verify PID extracted
	pid, err := client.SendCDROTP(ctx)
	if err != nil {
		t.Fatalf("SendCDROTP failed: %v", err)
	}
	if pid != "cdr-mock-uuid-555" {
		t.Errorf("expected PID cdr-mock-uuid-555, got %s", pid)
	}

	// 2. Confirm CDR OTP using cached PID
	if err := client.ConfirmCDROTP(ctx, "998877"); err != nil {
		t.Fatalf("ConfirmCDROTP failed: %v", err)
	}

	if receivedPID != "cdr-mock-uuid-555" {
		t.Errorf("expected received PID cdr-mock-uuid-555, got %s", receivedPID)
	}
	if receivedPasscode != "998877" {
		t.Errorf("expected received passcode 998877, got %s", receivedPasscode)
	}

	// 3. Confirm CDR OTP using explicit PID
	if err := client.ConfirmCDROTP(ctx, "custom-pid-111", "332211"); err != nil {
		t.Fatalf("ConfirmCDROTP with explicit PID failed: %v", err)
	}
	if receivedPID != "custom-pid-111" || receivedPasscode != "332211" {
		t.Errorf("explicit PID confirm mismatch: pid=%s, passcode=%s", receivedPID, receivedPasscode)
	}
}

func TestKeepAlive_PulseAndCDRExpired(t *testing.T) {
	var profileCalls int32
	var cdrCalls int32
	var cdrExpiredTriggered int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/v1/profile":
			atomic.AddInt32(&profileCalls, 1)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data": map[string]any{
					"bodies": []any{},
				},
			})

		case "/api/v1/cdr/detail":
			calls := atomic.AddInt32(&cdrCalls, 1)
			if calls >= 2 {
				// Second pulse: simulate CDR session expiration
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"message": "CDR session expired",
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data": map[string]any{
					"total": 0,
					"data":  []any{},
				},
			})
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := NewClient(
		WithTokens("mock-token", "mock-refresh"),
		WithOnCDRExpired(func() {
			atomic.AddInt32(&cdrExpiredTriggered, 1)
		}),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	pulseChan := client.StartKeepAlive(ctx, 30*time.Millisecond)

	// Receive pulse 1 (both profile and CDR OK)
	select {
	case pulse, ok := <-pulseChan:
		if !ok {
			t.Fatal("pulseChan closed unexpectedly")
		}
		if !pulse.ProfileOK || !pulse.CDROK {
			t.Errorf("pulse 1 expected ProfileOK and CDROK true, got profile=%v, cdr=%v", pulse.ProfileOK, pulse.CDROK)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for pulse 1")
	}

	// Receive pulse 2 (CDR expired)
	select {
	case pulse, ok := <-pulseChan:
		if !ok {
			t.Fatal("pulseChan closed unexpectedly")
		}
		if !pulse.CDRExpired {
			t.Errorf("pulse 2 expected CDRExpired true, got %v", pulse.CDRExpired)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for pulse 2")
	}

	if atomic.LoadInt32(&cdrExpiredTriggered) == 0 {
		t.Errorf("expected OnCDRExpired callback to be triggered")
	}

	// Cancel context to test clean termination
	cancel()

	// Verify pulseChan closes cleanly without leak
	select {
	case _, ok := <-pulseChan:
		if ok {
			// Might drain one last message
			select {
			case _, ok2 := <-pulseChan:
				if ok2 {
					t.Errorf("expected channel to close after cancellation")
				}
			case <-time.After(500 * time.Millisecond):
			}
		}
	case <-time.After(500 * time.Millisecond):
		t.Errorf("expected channel close after cancel")
	}
}

func TestSessionStorage_FileAndMemory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	initialSession := &SessionData{
		AccessToken:     "access-123",
		RefreshToken:    "refresh-456",
		HandshakeToken:  "handshake-789",
		UserID:          "usr-001",
		Username:        "07701112233",
		Phone:           "07701112233",
		DeviceID:        "dev-uuid-999",
		BiometricSecret: "secret-bio-777",
		Language:        "ar",
		LastRefresh:     time.Now(),
	}

	t.Run("MemorySessionStorage", func(t *testing.T) {
		memStorage := NewMemorySessionStorage()
		if err := memStorage.SaveSession(ctx, initialSession); err != nil {
			t.Fatalf("SaveSession failed: %v", err)
		}

		loaded, err := memStorage.LoadSession(ctx)
		if err != nil {
			t.Fatalf("LoadSession failed: %v", err)
		}

		if loaded.AccessToken != initialSession.AccessToken || loaded.BiometricSecret != initialSession.BiometricSecret {
			t.Errorf("loaded session mismatch: %+v", loaded)
		}
	})

	t.Run("FileSessionStorage", func(t *testing.T) {
		tmpDir := t.TempDir()
		sessionFile := filepath.Join(tmpDir, "asiacell_session.json")

		fileStorage := NewFileSessionStorage(sessionFile)
		if err := fileStorage.SaveSession(ctx, initialSession); err != nil {
			t.Fatalf("SaveSession failed: %v", err)
		}

		// Verify file exists on disk
		if _, err := os.Stat(sessionFile); err != nil {
			t.Fatalf("session file was not created: %v", err)
		}

		loaded, err := fileStorage.LoadSession(ctx)
		if err != nil {
			t.Fatalf("LoadSession failed: %v", err)
		}

		if loaded.DeviceID != "dev-uuid-999" || loaded.BiometricSecret != "secret-bio-777" {
			t.Errorf("loaded file session mismatch: %+v", loaded)
		}
	})

	t.Run("ClientStorageIntegration", func(t *testing.T) {
		memStorage := NewMemorySessionStorage()
		client, err := NewClient(
			WithSessionStorage(memStorage),
			WithDeviceID("dev-uuid-999"),
		)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}

		client.SetBiometricSecret("secret-bio-777")
		client.SetTokens("acc-token-99", "ref-token-99", "hand-99", "uid-99", "07701112233")

		// Storage should have been updated automatically by SetTokens
		saved, err := memStorage.LoadSession(ctx)
		if err != nil {
			t.Fatalf("storage should have saved session automatically: %v", err)
		}
		if saved.AccessToken != "acc-token-99" || saved.BiometricSecret != "secret-bio-777" {
			t.Errorf("unexpected saved session: %+v", saved)
		}

		// Create a new client and load from storage
		client2, err := NewClient()
		if err != nil {
			t.Fatalf("NewClient 2 failed: %v", err)
		}
		if err := client2.LoadFromStorage(ctx, memStorage); err != nil {
			t.Fatalf("LoadFromStorage failed: %v", err)
		}
		if client2.BiometricSecret() != "secret-bio-777" || client2.DeviceID() != "dev-uuid-999" {
			t.Errorf("client 2 state mismatch: secret=%s, devID=%s", client2.BiometricSecret(), client2.DeviceID())
		}
	})
}

func TestProactiveTokenRefresh(t *testing.T) {
	ctx := context.Background()

	// 1. Test JWT expiration parser with the real token from Asiacell session
	tok := "eyJhbGciOiJIUzUxMiJ9.eyJzZXNzaW9uSUQiOiJhMmY1MWU1NC1mMzM1LTQxYWEtOGE2Ny0wM2VhOGE4N2FhZmYiLCJleHAiOjE3ODkyOTI5NDV9.TdF2wMYW90b-FC7fYhDLAx7SsFDDFOJ71ORVHnKlv2gn9rmIljQaU7nd1ALqy-jMh2i0TNkrLJ9IqwPOCSfiVw"
	exp := parseJWTExpiration(tok)
	if exp.IsZero() {
		t.Fatalf("expected non-zero expiration from valid JWT token")
	}
	if exp.Unix() != 1789292945 {
		t.Fatalf("expected exp 1789292945, got %d", exp.Unix())
	}

	// 2. Test proactive refresh trigger when token is expiring soon
	var refreshCalled int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/validate":
			refreshCalled++
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(SMSValidationResponse{
				Success:      true,
				AccessToken:  "new-proactive-access-token",
				RefreshToken: "new-proactive-refresh-token",
			})
		case "/api/v1/profile":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data":    map[string]any{},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create an expired/soon-to-expire token: exp = now + 10 minutes (< 2 hours)
	expSoon := time.Now().Add(10 * time.Minute).Unix()
	claimsJSON := fmt.Sprintf(`{"sessionID":"test-sess","exp":%d}`, expSoon)
	payloadBase64 := base64.RawURLEncoding.EncodeToString([]byte(claimsJSON))
	soonExpiringToken := "eyJhbGciOiJIUzUxMiJ9." + payloadBase64 + ".sig"

	client, err := NewClient(
		WithTokens(soonExpiringToken, "existing-refresh-token"),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	// Before request, TokenExpiration should report the soonExpiringToken
	if client.TokenExpiration().Unix() != expSoon {
		t.Fatalf("expected client TokenExpiration to match expSoon")
	}

	// Calling GetProfile should trigger proactive refresh BEFORE sending the request
	_, err = client.GetProfile(ctx)
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}

	if refreshCalled != 1 {
		t.Fatalf("expected proactive refresh to be called 1 time, got %d", refreshCalled)
	}

	// Token should now be the new one
	if client.AccessToken() != "new-proactive-access-token" {
		t.Fatalf("expected updated access token, got %s", client.AccessToken())
	}
}

func TestCDROnboardingAndDetailCategories(t *testing.T) {
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/cdr":
			_ = json.NewEncoder(w).Encode(CDROnboardingResponse{
				Success: true,
				Data: &CDROnboardingData{
					Headers: []CDROnboardingHeader{
						{BackgroundImage: "https://odpapp.asiacell.com/img/my-pocket-onboarding-ar.jpg"},
					},
					Body: []CDROnboardingBody{
						{
							Desc:       "مع خدمة سجل الاستخدام...",
							Disclaimer: "سوف نتحقق من رقمك باستخدام رمز PIN",
							ActionButton: &ServiceActionItem{
								Title: "أطلب الـ PIN عبر SMS",
							},
						},
					},
				},
			})
		case "/api/v1/cdr/detail":
			cdrType := r.URL.Query().Get("type")
			_ = json.NewEncoder(w).Encode(CDRDetailResponse{
				Success: true,
				Data: &CDRDetailData{
					Total: 1,
					Data: []CDRRecord{
						{
							Title:  FlexString("سجل " + cdrType),
							Amount: FlexString("100"),
							Unit:   FlexString(cdrType),
						},
					},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := NewClient(WithTokens("valid-tok", "valid-ref"))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	client.baseURL = server.URL

	onboarding, err := client.GetCDROnboarding(ctx)
	if err != nil {
		t.Fatalf("GetCDROnboarding failed: %v", err)
	}
	if !onboarding.Success || len(onboarding.Data.Body) == 0 || onboarding.Data.Body[0].ActionButton.Title != "أطلب الـ PIN عبر SMS" {
		t.Fatalf("unexpected onboarding response: %+v", onboarding)
	}

	for _, category := range []string{"voice", "data", "sms", "btransfer", "cmp"} {
		records, err := client.GetCDRDetail(ctx, category, 1, 10)
		if err != nil {
			t.Fatalf("GetCDRDetail(%s) failed: %v", category, err)
		}
		if len(records) != 1 || string(records[0].Unit) != category {
			t.Fatalf("unexpected records for category %s: %+v", category, records)
		}
	}
}


