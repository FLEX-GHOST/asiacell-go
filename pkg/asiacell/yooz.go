package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetYoozHome(ctx context.Context) (*YoozHomeData, error) {
	path := fmt.Sprintf("/api/v5/avocado/home?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz home: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozHomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz home: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetYoozBundlesScreen(ctx context.Context, groupID string) (*YoozBundlesScreenData, error) {
	path := fmt.Sprintf("/api/v3/avocado/bundles/screen?groupId=%s&lang=%s", url.QueryEscape(groupID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz bundles screen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozAddOnListResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz bundles screen: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetYoozClassicPlans(ctx context.Context) (*YoozClassicPlansData, error) {
	path := fmt.Sprintf("/api/v3/avocado/bundles/classic-plans?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz classic plans: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozClassicPlansResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz classic plans: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetYoozOmegaPlans(ctx context.Context, voucher, msisdn string) (*YoozOmegaPlansData, error) {
	path := fmt.Sprintf("/api/v3/avocado/bundles/omega-plans?voucher=%s&msisdn=%s&lang=%s",
		url.QueryEscape(voucher), url.QueryEscape(msisdn), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz omega plans: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozOmegaPlansResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz omega plans: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetYoozBundles(ctx context.Context) ([]YoozPlanEntity, error) {
	path := fmt.Sprintf("/api/v2/avocado/bundles?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz bundles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozBundlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz bundles: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetYoozDataCap(ctx context.Context) (*YoozDataCapData, error) {
	path := fmt.Sprintf("/api/v2/avocado/data-cap?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz datacap: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozDataCapResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz datacap: %w", err)
	}
	return res.Data, nil
}

func (c *Client) SetYoozDataCap(ctx context.Context, limitMB int) error {
	payload, err := json.Marshal(map[string]int{
		"limit": limitMB,
	})
	if err != nil {
		return fmt.Errorf("encoding datacap payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/avocado/data-cap?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting yooz datacap: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetYoozReward(ctx context.Context) (*YoozRewardData, error) {
	path := fmt.Sprintf("/api/v1/avocado/reward?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz reward: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozRewardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz reward: %w", err)
	}
	return res.Data, nil
}
