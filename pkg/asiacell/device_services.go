package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (c *Client) UploadProfileImage(ctx context.Context, filename string, imageReader io.Reader) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, imageReader); err != nil {
		return fmt.Errorf("copying image data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing multipart writer: %w", err)
	}

	path := fmt.Sprintf("/api/v1/profile/upload?lang=%s", c.language)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &body)
	if err != nil {
		return fmt.Errorf("creating upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("DeviceId", c.deviceID)
	req.Header.Set("X-ODP-API-KEY", defaultAPIKey)
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) UploadPartnerLogo(ctx context.Context, filename string, imageReader io.Reader) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("creating logo form file: %w", err)
	}
	if _, err := io.Copy(part, imageReader); err != nil {
		return fmt.Errorf("copying logo data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing multipart writer: %w", err)
	}

	path := fmt.Sprintf("/api/v2/logo/upload?lang=%s", c.language)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, &body)
	if err != nil {
		return fmt.Errorf("creating logo upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("DeviceId", c.deviceID)
	req.Header.Set("X-ODP-API-KEY", defaultAPIKey)
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing logo upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) RegisterNotificationToken(ctx context.Context, token, osName string) error {
	payload, err := json.Marshal(map[string]string{
		"token": token,
		"os":    osName,
	})
	if err != nil {
		return fmt.Errorf("encoding notification token payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/notifications/register?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("requesting notification register: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) LogNotificationRead(ctx context.Context, notificationID string) error {
	payload, err := json.Marshal(map[string]string{
		"notificationId": notificationID,
	})
	if err != nil {
		return fmt.Errorf("encoding notification log payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/notifications/log?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("requesting notification log: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) RegisterBiometrics(ctx context.Context) error {
	path := fmt.Sprintf("/api/v1/biometrics/register?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("requesting biometrics register: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) BiometricLogin(ctx context.Context, biometricsKey string) (*LoginResponse, error) {
	payload, err := json.Marshal(map[string]string{
		"key": biometricsKey,
	})
	if err != nil {
		return nil, fmt.Errorf("encoding biometric login payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/biometrics/do-login?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("requesting biometric login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding biometric login response: %w", err)
	}
	return &res, nil
}

func (c *Client) GetWatchDashboard(ctx context.Context) (*WatchDashboardData, error) {
	path := fmt.Sprintf("/api/v1/watch/home?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting watch dashboard: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res WatchDashboardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding watch dashboard response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetProtectedPaymentStatus(ctx context.Context, transactionID string) (*ProtectedPaymentStatusData, error) {
	path := fmt.Sprintf("/protected/v1/payments/%s/status?lang=%s", url.PathEscape(transactionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting protected payment status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ProtectedPaymentStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding protected payment status: %w", err)
	}
	return res.Data, nil
}

func (c *Client) CancelProtectedPayment(ctx context.Context, transactionID string) error {
	path := fmt.Sprintf("/protected/v1/payments/%s/cancel?lang=%s", url.PathEscape(transactionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("canceling protected payment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) TopUpOmega(ctx context.Context, phone, voucher string) error {
	payload, err := json.Marshal(map[string]string{
		"msisdn":  phone,
		"voucher": voucher,
	})
	if err != nil {
		return fmt.Errorf("encoding omega topup payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/top-up/omega?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("requesting omega topup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) DeleteAccount(ctx context.Context) error {
	path := fmt.Sprintf("/api/v1/delete?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("requesting account deletion: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetAsiaverseHome(ctx context.Context) (*AsiaverseHomeData, error) {
	path := fmt.Sprintf("/api/v1/asiaverse?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting asiaverse home: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res AsiaverseHomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding asiaverse home response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetShazamScanStatus(ctx context.Context) (*ScanToWinData, error) {
	path := fmt.Sprintf("/api/v1/shazam?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting shazam status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ScanToWinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding shazam status response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) SubmitShazamScan(ctx context.Context, qrCode string) error {
	payload, err := json.Marshal(map[string]string{
		"code": qrCode,
	})
	if err != nil {
		return fmt.Errorf("encoding shazam scan payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/shazam?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting shazam scan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetUserInterests(ctx context.Context) ([]string, error) {
	path := fmt.Sprintf("/api/v1/profile/interests?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting user interests: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res UserInterestsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding user interests: %w", err)
	}
	return res.Data, nil
}

func (c *Client) SaveUserInterests(ctx context.Context, interests []string) error {
	payload, err := json.Marshal(interests)
	if err != nil {
		return fmt.Errorf("encoding user interests: %w", err)
	}

	path := fmt.Sprintf("/api/v1/avocado/profile/save-interests?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("saving user interests: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetCDRSummary(ctx context.Context) (*CDRSummaryData, error) {
	path := fmt.Sprintf("/api/v1/cdr/summary?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting cdr summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res CDRSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding cdr summary: %w", err)
	}
	return res.Data, nil
}

func (c *Client) ResendLinkedAccountSMS(ctx context.Context, phone string) error {
	payload, err := json.Marshal(map[string]string{
		"msisdn": phone,
	})
	if err != nil {
		return fmt.Errorf("encoding resend sms payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/map-account/resend?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("resending linked account sms: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}
