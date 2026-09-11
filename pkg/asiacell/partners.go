package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetPartnerCategories(ctx context.Context) ([]EOCategory, error) {
	path := fmt.Sprintf("/api/v1/partners/categories?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting partner categories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EOCategoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding partner categories: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetPartnerCities(ctx context.Context) ([]EOCity, error) {
	path := fmt.Sprintf("/api/v2/partners/cities?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting partner cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EOCityResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding partner cities: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetPartnerCityCategories(ctx context.Context, cityID int) ([]EOCategory, error) {
	path := fmt.Sprintf("/api/v2/partners/cities/%d/categories?lang=%s", cityID, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting city categories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EOCategoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding city categories: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetPartnersByCategoryAndCity(ctx context.Context, categoryID, cityID int) ([]EOPartner, error) {
	path := fmt.Sprintf("/api/v2/categories/%d/cities/%d/partners?lang=%s", categoryID, cityID, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting partners list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EOPartnerResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding partners list: %w", err)
	}
	return res.Data, nil
}

func (c *Client) RegisterPartner(ctx context.Context, req PartnerRegisterRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("encoding partner register payload: %w", err)
	}

	path := fmt.Sprintf("/api/v2/partners/register?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("requesting partner registration: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}
