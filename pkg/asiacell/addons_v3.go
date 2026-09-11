package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetHomeDashboardV3(ctx context.Context, lat, lng float64, roaming bool) (*HomeDashboardData, error) {
	path := fmt.Sprintf("/api/v3/home?lat=%f&lng=%f&roaming=%t&lang=%s", lat, lng, roaming, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting home dashboard v3: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res HomeDashboardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding home dashboard v3: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetCYOBundles(ctx context.Context, groupID string) (*CYOBundlesData, error) {
	path := fmt.Sprintf("/api/v3/addon/cyo?groupId=%s&lang=%s", url.QueryEscape(groupID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting cyo bundles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res CYOBundlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding cyo bundles: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetActiveOffers(ctx context.Context) (*ActiveOffersData, error) {
	path := fmt.Sprintf("/api/v1/offer/active-offers?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting active offers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ActiveOffersResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding active offers: %w", err)
	}
	return res.Data, nil
}

func (c *Client) ManageQuickActions(ctx context.Context, actionIDs []int) error {
	payload, err := json.Marshal(actionIDs)
	if err != nil {
		return fmt.Errorf("encoding quick actions payload: %w", err)
	}

	path := fmt.Sprintf("/api/v3/profile/manage-quick-actions?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting quick actions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) SubmitAppFeedback(ctx context.Context, category, comment string, rating int) error {
	payload, err := json.Marshal(map[string]interface{}{
		"category": category,
		"comment":  comment,
		"rating":   rating,
	})
	if err != nil {
		return fmt.Errorf("encoding feedback payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/app-feedback?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting app feedback: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetSurveys(ctx context.Context) (*SurveyData, error) {
	path := fmt.Sprintf("/api/v1/survey?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting surveys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res SurveyResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding surveys: %w", err)
	}
	return res.Data, nil
}

func (c *Client) SubmitSurvey(ctx context.Context, surveyID string, answers map[string]string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"surveyId": surveyID,
		"answers":  answers,
	})
	if err != nil {
		return fmt.Errorf("encoding survey submit payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/survey?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting survey: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetVoiceOfCustomer(ctx context.Context) (*VoCEntity, error) {
	path := fmt.Sprintf("/api/v1/voc?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting voc: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res VoCResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding voc: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetVideoTutorials(ctx context.Context) ([]VideoTutorialItem, error) {
	path := fmt.Sprintf("/api/v1/promotions/video-tutorials?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting video tutorials: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res VideoTutorialsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding video tutorials: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetOneYadHome(ctx context.Context) (*OneYadData, error) {
	path := fmt.Sprintf("/api/v1/one-yad?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting one yad home: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res OneYadResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding one yad home: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetOneYadTeams(ctx context.Context) ([]OneYadTeam, error) {
	path := fmt.Sprintf("/api/v1/one-yad/teams?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting one yad teams: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res OneYadTeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding one yad teams: %w", err)
	}
	return res.Data, nil
}

func (c *Client) SubmitOneYadRequest(ctx context.Context, teamID string, amount float64) error {
	payload, err := json.Marshal(map[string]interface{}{
		"teamId": teamID,
		"amount": amount,
	})
	if err != nil {
		return fmt.Errorf("encoding one yad payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/one-yad/request?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting one yad request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetEpicLines(ctx context.Context) (*EpicLinesData, error) {
	path := fmt.Sprintf("/api/v1/epic?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting epic lines: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EpicLinesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding epic lines: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetEpicLineUsage(ctx context.Context, msisdn string) (*EpicLineUsageData, error) {
	path := fmt.Sprintf("/api/v1/epic/remaining/%s?lang=%s", url.PathEscape(msisdn), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting epic line usage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EpicLineUsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding epic line usage: %w", err)
	}
	return res.Data, nil
}
