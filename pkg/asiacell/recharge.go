package asiacell

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetRechargeNumbers(ctx context.Context, option string) (*RechargeNumberData, error) {
	path := fmt.Sprintf("/api/v1/recharge/screen1?option=%s&lang=%s", url.QueryEscape(option), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting recharge numbers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res RechargeNumberResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding recharge numbers: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetRechargeTypes(ctx context.Context, option, msisdn string) (*RechargeTypeData, error) {
	path := fmt.Sprintf("/api/v1/recharge/screen2?option=%s&msisdn=%s&lang=%s",
		url.QueryEscape(option), url.QueryEscape(msisdn), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting recharge types: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res RechargeTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding recharge types: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetRechargeMethods(ctx context.Context, option, msisdn, method string) (*RechargeMethodData, error) {
	path := fmt.Sprintf("/api/v1/recharge/screen3?option=%s&msisdn=%s&method=%s&lang=%s",
		url.QueryEscape(option), url.QueryEscape(msisdn), url.QueryEscape(method), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting recharge methods: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res RechargeMethodResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding recharge methods: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetOnlinePaymentDetails(ctx context.Context, option, msisdn, method, pgName string) (*OnlinePaymentData, error) {
	path := fmt.Sprintf("/api/v1/recharge/screen4?option=%s&msisdn=%s&method=%s&pgName=%s&lang=%s",
		url.QueryEscape(option), url.QueryEscape(msisdn), url.QueryEscape(method), url.QueryEscape(pgName), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting online payment details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res OnlinePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding online payment details: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetRechargeConfirmation(ctx context.Context, transactionID string) (*RechargeConfirmationData, error) {
	path := fmt.Sprintf("/api/v1/recharge/confirmation?transactionId=%s&lang=%s", url.QueryEscape(transactionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting recharge confirmation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res RechargeConfirmationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding recharge confirmation: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetPaymentSelection(ctx context.Context) (*PaymentSelectionData, error) {
	path := fmt.Sprintf("/api/v1/top-up/payment-selection?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting payment selection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res PaymentSelectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding payment selection: %w", err)
	}
	return res.Data, nil
}
