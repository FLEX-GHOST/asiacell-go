package asiacell

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL   = "https://app.asiacell.com"
	defaultLanguage  = "ar"
	defaultUserAgent = "okhttp/5.0.0-alpha.2"
	defaultAPIKey    = "1ccbc4c913bc4ce785a0a2de444aa0d6"
	defaultAppVer    = "4.2.5"
	defaultDeviceTyp = "[Android][INFINIX][Infinix X6871 15][VANILLA_ICE_CREAM][HMS][4.2.5:90000256]"
)

type Option func(*Client)

type Client struct {
	httpClient     *http.Client
	baseURL        string
	deviceID       string
	language       string
	userAgent      string
	accessToken    string
	refreshToken   string
	handshakeToken string
	userID            string
	username          string
	masterWallet      string
	lastRefresh       time.Time
	onTokenUpdate     func(*SessionData)
	recordedTransfers         []TransactionRecord
	recordedIncomingTransfers []TransactionRecord
	recordedRecharges         []TransactionRecord
	mu                        sync.RWMutex
}


func WithLanguage(lang string) Option {
	return func(c *Client) {
		if lang != "" {
			c.language = lang
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.httpClient.Timeout = timeout
		}
	}
}

func WithDeviceID(deviceID string) Option {
	return func(c *Client) {
		if deviceID != "" {
			c.deviceID = deviceID
		}
	}
}

func WithTokens(accessToken, refreshToken string) Option {
	return func(c *Client) {
		c.accessToken = accessToken
		c.refreshToken = refreshToken
	}
}

func WithMasterWallet(wallet string) Option {
	return func(c *Client) {
		c.masterWallet = wallet
	}
}

func WithOnTokenUpdate(fn func(*SessionData)) Option {
	return func(c *Client) {
		c.onTokenUpdate = fn
	}
}

func generateDeviceID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generating device id: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func NewClient(opts ...Option) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("creating cookie jar: %w", err)
	}

	devID, err := generateDeviceID()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}

	client := &Client{
		httpClient: &http.Client{
			Jar:       jar,
			Timeout:   20 * time.Second,
			Transport: transport,
		},
		baseURL:   defaultBaseURL,
		deviceID:  devID,
		language:  defaultLanguage,
		userAgent: defaultUserAgent,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Client) SetOnTokenUpdate(fn func(*SessionData)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onTokenUpdate = fn
}

func (c *Client) SetTokens(access, refresh, handshake, userID, username string) {
	c.mu.Lock()
	c.accessToken = access
	c.refreshToken = refresh
	c.handshakeToken = handshake
	c.userID = userID
	c.username = username
	c.lastRefresh = time.Now()
	cb := c.onTokenUpdate
	var data *SessionData
	if cb != nil && access != "" {
		data = &SessionData{
			AccessToken:    c.accessToken,
			RefreshToken:   c.refreshToken,
			HandshakeToken: c.handshakeToken,
			UserID:         FlexString(c.userID),
			Username:       c.username,
			DeviceID:       c.deviceID,
			Language:       c.language,
			LastRefresh:    c.lastRefresh,
		}
	}
	c.mu.Unlock()

	if cb != nil && data != nil {
		cb(data)
	}
}

func (c *Client) GetTokens() (access, refresh, handshake, userID, username string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.accessToken, c.refreshToken, c.handshakeToken, c.userID, c.username
}

func (c *Client) DeviceID() string {
	return c.deviceID
}

func (c *Client) RotateDeviceID() (string, error) {
	devID, err := generateDeviceID()
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.deviceID = devID
	c.mu.Unlock()
	return devID, nil
}

func (c *Client) Language() string {
	return c.language
}

func (c *Client) SetMasterWallet(wallet string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.masterWallet = wallet
}

func (c *Client) MasterWallet() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.masterWallet
}

func (c *Client) applyHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("DeviceId", c.deviceID)
	req.Header.Set("X-ODP-API-KEY", defaultAPIKey)
	req.Header.Set("X-OS-Version", "15")
	req.Header.Set("X-ODP-APP-VERSION", defaultAppVer)
	req.Header.Set("X-Device-Type", defaultDeviceTyp)
	req.Header.Set("X-FROM-APP", "odp")
	req.Header.Set("X-ODP-CHANNEL", "mobile")
	req.Header.Set("X-SCREEN-TYPE", "MOBILE")

	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	hosts := []string{c.baseURL}
	altHost := "https://odpapp.asiacell.com"
	if c.baseURL == "https://odpapp.asiacell.com" {
		altHost = "https://app.asiacell.com"
	}
	hosts = append(hosts, altHost)

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("reading request body: %w", err)
		}
	}

	var lastErr error
	for _, host := range hosts {
		fullURL := host + path
		var rdr io.Reader
		if bodyBytes != nil {
			rdr = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, rdr)
		if err != nil {
			lastErr = err
			continue
		}

		c.applyHeaders(req)
		if bodyBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized && path != "/api/v1/validate" && path != "/api/v1/login" && path != "/api/v1/smsvalidation" {
			c.mu.RLock()
			hasRefresh := c.refreshToken != ""
			c.mu.RUnlock()

			if hasRefresh && bodyBytes == nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				if refErr := c.RefreshToken(ctx); refErr == nil {
					retryReq, retryErr := http.NewRequestWithContext(ctx, method, fullURL, nil)
					if retryErr == nil {
						c.applyHeaders(retryReq)
						resp, err = c.httpClient.Do(retryReq)
						if err != nil {
							lastErr = err
							continue
						}
					}
				}
			}
		}

		br := bufio.NewReader(resp.Body)
		peekBytes, _ := br.Peek(512)
		peekTrimmed := strings.TrimSpace(string(peekBytes))
		if strings.HasPrefix(peekTrimmed, "<") || strings.Contains(peekTrimmed, "Access Blocked") || strings.Contains(peekTrimmed, "<!DOCTYPE") {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("host %s returned HTML response instead of JSON", host)
			continue
		}

		resp.Body = &readCloser{
			reader: br,
			closer: resp.Body,
		}

		return resp, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ErrRequestFailed
}

func (c *Client) RecordTransfer(rec TransactionRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedTransfers = append([]TransactionRecord{rec}, c.recordedTransfers...)
}

func (c *Client) RecordRecharge(rec TransactionRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedRecharges = append([]TransactionRecord{rec}, c.recordedRecharges...)
}

func (c *Client) RecordedTransfers() []TransactionRecord {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]TransactionRecord, len(c.recordedTransfers))
	copy(out, c.recordedTransfers)
	return out
}

func (c *Client) SetRecordedTransfers(records []TransactionRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedTransfers = records
}

func (c *Client) RecordIncomingTransfer(rec TransactionRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedIncomingTransfers = append([]TransactionRecord{rec}, c.recordedIncomingTransfers...)
}

func (c *Client) RecordedIncomingTransfers() []TransactionRecord {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]TransactionRecord, len(c.recordedIncomingTransfers))
	copy(out, c.recordedIncomingTransfers)
	return out
}

func (c *Client) SetRecordedIncomingTransfers(records []TransactionRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.recordedIncomingTransfers = records
}

type readCloser struct {
	reader io.Reader
	closer io.Closer
}

func (rc *readCloser) Read(p []byte) (n int, err error) {
	return rc.reader.Read(p)
}

func (rc *readCloser) Close() error {
	_, _ = io.Copy(io.Discard, rc.reader)
	return rc.closer.Close()
}


