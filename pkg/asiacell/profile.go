package asiacell

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetProfile(ctx context.Context) (*AccountOverview, error) {
	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()

	if token == "" {
		return nil, ErrUnauthorized
	}

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/profile?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: "profile endpoint failed"}
	}

	var profResp ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profResp); err != nil {
		return nil, fmt.Errorf("decoding profile response: %w", err)
	}

	if !profResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, profResp.Message)
	}

	overview := &AccountOverview{}
	bodies := profResp.Data.Bodies

	for idx, body := range bodies {
		bType := body.Type
		if bType == "" {
			switch idx {
			case 0:
				bType = "profile"
			case 2:
				bType = "balance"
			case 3:
				bType = "remainingData"
			case 4:
				if overview.RemainingCalls == "" && len(body.Items) > 0 {
					overview.RemainingCalls = string(body.Items[0].Title)
				}
			case 5:
				if overview.RemainingSMS == "" && len(body.Items) > 0 {
					overview.RemainingSMS = string(body.Items[0].Title)
				}
			}
		}

		switch bType {
		case "profile":
			if len(body.Items) > 0 {
				rawName := string(body.Items[0].Name)
				overview.Name = strings.TrimSpace(strings.ReplaceAll(rawName, "null", ""))
				overview.PhoneNumber = string(body.Items[0].PhoneNumber)
			}
		case "balance":
			if len(body.Items) > 0 {
				overview.Balance = string(body.Items[0].Value)
				overview.Validity = body.Items[0].Validity
			}
		case "remainingData":
			for _, it := range body.Items {
				if it.RemainingVolume > 0 || it.ActiveBundle {
					bInfo := ActiveBundleInfo{
						Title:           string(it.Title),
						RemainingVolume: it.RemainingVolume,
						TotalVolume:     it.TotalVolume,
						Unit:            it.Unit,
						ExpireLabel:     it.ExpireLabel,
						ExpireDate:      it.ExpireDate,
						SubscribeDate:   it.SubscribeDate,
						BundleKey:       it.BundleKey,
						FreeUnitRefName: it.FreeUnitRefName,
						ActiveBundle:    it.ActiveBundle,
					}
					overview.ActiveBundles = append(overview.ActiveBundles, bInfo)
					if overview.RemainingData == "" {
						name := it.FreeUnitRefName
						if name == "" {
							name = string(it.Title)
						}
						overview.RemainingData = fmt.Sprintf("%.0f %s (%s)", it.RemainingVolume, it.Unit, name)
					}
				}
			}
			if overview.RemainingData == "" && len(body.Items) > 0 {
				overview.RemainingData = string(body.Items[0].Title)
			}
		}

		titleStr := string(body.Title)
		if strings.Contains(titleStr, "المکالمات") || strings.Contains(titleStr, "مكالمات") {
			if len(body.Items) > 0 {
				overview.RemainingCalls = string(body.Items[0].Title)
			}
		}
		if strings.Contains(titleStr, "الرسائل") || strings.Contains(titleStr, "رسائل") {
			if len(body.Items) > 0 {
				overview.RemainingSMS = string(body.Items[0].Title)
			}
		}
	}

	if overview.Name == "" && len(bodies) > 0 && len(bodies[0].Items) > 0 {
		rawName := string(bodies[0].Items[0].Name)
		overview.Name = strings.TrimSpace(strings.ReplaceAll(rawName, "null", ""))
		overview.PhoneNumber = string(bodies[0].Items[0].PhoneNumber)
	}

	return overview, nil
}

func (c *Client) GetActiveBundles(ctx context.Context) ([]ActiveBundleInfo, error) {
	overview, err := c.GetProfile(ctx)
	if err != nil {
		return nil, err
	}
	return overview.ActiveBundles, nil
}

func (c *Client) GetProfileDetails(ctx context.Context) (*ProfileDetails, error) {
	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()

	if token == "" {
		return nil, ErrUnauthorized
	}

	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/profile/view?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting profile details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Message: "profile view endpoint failed"}
	}

	var detailsResp ProfileDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&detailsResp); err != nil {
		return nil, fmt.Errorf("decoding profile details: %w", err)
	}

	if !detailsResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, detailsResp.Message)
	}

	return &detailsResp.Data, nil
}
