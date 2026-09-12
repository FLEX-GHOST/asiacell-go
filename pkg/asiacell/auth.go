package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultFCMToken = "ersrsDcnSmeDyGflUeMnwN:APA91bHneFTyD2RLqLK-JU1kUXQM9q0xJAFkGetsnfEFyI0Nn-fbRWxee9oWTn8SKFjUMqBhIVcjsLMr47c3y0Q2HXZUae9fWwHWG2tUR4LoMnOGhg2yQMQ"

func extractPIDFromURL(rawURL string) (string, error) {
	idx := strings.Index(rawURL, "PID=")
	if idx != -1 {
		rest := rawURL[idx+4:]
		end := strings.IndexAny(rest, "&?# ")
		if end != -1 {
			return rest[:end], nil
		}
		return rest, nil
	}

	parsed, err := url.Parse(rawURL)
	if err == nil {
		if pid := parsed.Query().Get("PID"); pid != "" {
			return pid, nil
		}
	}

	return "", ErrMissingPID
}

func (c *Client) Login(ctx context.Context, phone string) (string, error) {
	reqBody := LoginRequest{
		Username:    phone,
		CaptchaCode: "",
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling login payload: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/login?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("executing login request: %w", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("decoding login response: %w", err)
	}

	if !loginResp.Success {
		if loginResp.RequireCaptcha || strings.Contains(strings.ToLower(loginResp.Message), "captcha") {
			retryResp, retryErr := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/login?lang=%s", c.language), bytes.NewReader(payload))
			if retryErr == nil {
				defer retryResp.Body.Close()
				var retryLoginResp LoginResponse
				if decodeErr := json.NewDecoder(retryResp.Body).Decode(&retryLoginResp); decodeErr == nil && retryLoginResp.Success {
					return extractPIDFromURL(retryLoginResp.NextURL)
				}
			}
			return "", ErrCaptchaRequired
		}
		return "", fmt.Errorf("%w: %s", ErrRequestFailed, loginResp.Message)
	}

	pid, err := extractPIDFromURL(loginResp.NextURL)
	if err != nil {
		return "", err
	}

	return pid, nil
}

func (c *Client) VerifySMS(ctx context.Context, pid, passcode string) (*SMSValidationResponse, error) {
	reqBody := SMSValidationRequest{
		PID:      pid,
		Passcode: passcode,
		Token:    defaultFCMToken,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling sms validation payload: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/smsvalidation?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("executing sms validation request: %w", err)
	}
	defer resp.Body.Close()

	var smsResp SMSValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return nil, fmt.Errorf("decoding sms validation response: %w", err)
	}

	if (!smsResp.Success && smsResp.AccessToken == "") || smsResp.AccessToken == "" {
		msg := strings.ToLower(smsResp.Message)
		if strings.Contains(msg, "passcode") || strings.Contains(msg, "code") || strings.Contains(smsResp.Message, "رمز") || strings.Contains(smsResp.Message, "تأكيد") {
			return nil, ErrInvalidPasscode
		}
		if smsResp.Message != "" {
			return nil, fmt.Errorf("%w: %s", ErrRequestFailed, smsResp.Message)
		}
		return nil, ErrRequestFailed
	}

	c.SetTokens(smsResp.AccessToken, smsResp.RefreshToken, smsResp.HandshakeToken, string(smsResp.UserID), smsResp.Username)
	if smsResp.Secret != "" {
		c.SetBiometricSecret(smsResp.Secret)
	}
	_ = c.RegisterBiometrics(ctx)

	return &smsResp, nil
}

func (c *Client) RefreshToken(ctx context.Context) error {
	c.mu.RLock()
	refToken := c.refreshToken
	c.mu.RUnlock()

	if refToken == "" {
		return ErrUnauthorized
	}

	reqBody := RefreshTokenRequest{
		RefreshToken: "Bearer " + refToken,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling refresh payload: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/validate?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("executing token refresh request: %w", err)
	}
	defer resp.Body.Close()

	var refreshResp SMSValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return fmt.Errorf("decoding refresh response: %w", err)
	}

	if !refreshResp.Success {
		return fmt.Errorf("%w: %s", ErrRequestFailed, refreshResp.Message)
	}

	if refreshResp.Secret != "" {
		c.SetBiometricSecret(refreshResp.Secret)
	}

	c.SetTokens(refreshResp.AccessToken, refreshResp.RefreshToken, refreshResp.HandshakeToken, string(refreshResp.UserID), refreshResp.Username)

	return nil
}

func (c *Client) RefreshSession(ctx context.Context) error {
	if err := c.RefreshToken(ctx); err == nil {
		return nil
	}

	c.mu.RLock()
	sec := c.biometricSecret
	phone := c.phone
	if phone == "" {
		phone = c.username
	}
	c.mu.RUnlock()

	if sec != "" && phone != "" {
		if _, err := c.LoginBiometric(ctx, sec); err == nil {
			return nil
		}
	}

	return ErrUnauthorized
}

func (c *Client) LoginBiometric(ctx context.Context, optionalSecret ...string) (*LoginResponse, error) {
	c.mu.RLock()
	sec := c.biometricSecret
	phone := c.phone
	if phone == "" {
		phone = c.username
	}
	c.mu.RUnlock()

	if len(optionalSecret) > 0 && optionalSecret[0] != "" {
		sec = optionalSecret[0]
	}

	if sec == "" {
		return nil, fmt.Errorf("biometric login requires secret: %w", ErrUnauthorized)
	}

	reqMap := map[string]string{
		"secret": sec,
		"key":    sec,
	}
	if phone != "" {
		reqMap["msisdn"] = phone
	}

	payload, err := json.Marshal(reqMap)
	if err != nil {
		return nil, fmt.Errorf("marshaling biometric login payload: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/biometrics/do-login?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("executing biometric login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden || resp.StatusCode == 493 {
		return nil, ErrUnauthorized
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading biometric login response: %w", err)
	}

	var bioResp struct {
		Success        bool       `json:"success"`
		AccessToken    string     `json:"access_token"`
		RefreshToken   string     `json:"refresh_token"`
		HandshakeToken string     `json:"handshake_token"`
		UserID         FlexString `json:"userId"`
		Username       string     `json:"username"`
		Secret         string     `json:"secret,omitempty"`
		Message        string     `json:"message,omitempty"`
	}

	if err := json.Unmarshal(bodyBytes, &bioResp); err != nil {
		return nil, fmt.Errorf("decoding biometric login response: %w", err)
	}

	if !bioResp.Success && bioResp.AccessToken == "" {
		if bioResp.Message != "" {
			return nil, fmt.Errorf("%w: %s", ErrRequestFailed, bioResp.Message)
		}
		return nil, ErrRequestFailed
	}

	if bioResp.Username == "" {
		bioResp.Username = phone
	}

	if bioResp.Secret != "" {
		c.SetBiometricSecret(bioResp.Secret)
	}

	if bioResp.AccessToken != "" {
		c.SetTokens(bioResp.AccessToken, bioResp.RefreshToken, bioResp.HandshakeToken, string(bioResp.UserID), bioResp.Username)
	}

	return &LoginResponse{
		Success:  bioResp.Success,
		Username: bioResp.Username,
		Message:  bioResp.Message,
	}, nil
}

func (c *Client) GetCaptcha(ctx context.Context) (*CaptchaData, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/captcha?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting captcha: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading captcha response: %w", err)
	}
	if len(strings.TrimSpace(string(bodyBytes))) == 0 {
		return nil, nil
	}

	var captResp CaptchaResponse
	if err := json.Unmarshal(bodyBytes, &captResp); err != nil {
		return nil, fmt.Errorf("decoding captcha response: %w", err)
	}

	return &captResp.Data, nil
}

func (c *Client) LoginWithCaptcha(ctx context.Context, phone, captchaCode string) (string, error) {
	reqBody := LoginRequest{
		Username:    phone,
		CaptchaCode: captchaCode,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling login payload: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/login?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("executing login request: %w", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("decoding login response: %w", err)
	}

	if !loginResp.Success {
		return "", fmt.Errorf("%w: %s", ErrRequestFailed, loginResp.Message)
	}

	pid, err := extractPIDFromURL(loginResp.NextURL)
	if err != nil {
		return "", err
	}

	return pid, nil
}

func (c *Client) SolveCaptchaOCR(ctx context.Context, imageURLOrBase64, apiKey string) (string, error) {
	if apiKey == "" {
		apiKey = "K88888888888957"
	}

	data := url.Values{}
	data.Set("language", "eng")
	data.Set("isOverlayRequired", "false")
	data.Set("OCREngine", "2")
	if strings.HasPrefix(imageURLOrBase64, "data:") || !strings.HasPrefix(imageURLOrBase64, "http") {
		data.Set("base64Image", imageURLOrBase64)
	} else {
		data.Set("url", imageURLOrBase64)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.ocr.space/parse/image", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating ocr request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("apikey", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling ocr service: %w", err)
	}
	defer resp.Body.Close()

	var ocrResp struct {
		ParsedResults []struct {
			ParsedText string `json:"ParsedText"`
		} `json:"ParsedResults"`
		IsErroredOnProcessing bool `json:"IsErroredOnProcessing"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return "", fmt.Errorf("decoding ocr response: %w", err)
	}

	if len(ocrResp.ParsedResults) == 0 {
		return "", fmt.Errorf("no text recognized in captcha")
	}

	cleanText := strings.TrimSpace(ocrResp.ParsedResults[0].ParsedText)
	cleanText = strings.ReplaceAll(cleanText, " ", "")
	cleanText = strings.ReplaceAll(cleanText, "\r", "")
	cleanText = strings.ReplaceAll(cleanText, "\n", "")
	cleanText = strings.ReplaceAll(cleanText, "\t", "")

	if cleanText == "" {
		return "", fmt.Errorf("recognized captcha text is empty")
	}

	return cleanText, nil
}
