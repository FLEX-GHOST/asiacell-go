package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetFanZoneHome(ctx context.Context, competitionID string) (*FanZoneHomeData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/home?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone home: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneHomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone home: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetFanZoneKickAndWin(ctx context.Context, competitionID string) (*FanZoneKickAndWinData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/kick-and-win/home?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone kick and win: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneKickAndWinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone kick and win: %w", err)
	}
	return res.Data, nil
}

func (c *Client) FinishFanZoneKickAndWin(ctx context.Context, competitionID string, score int) error {
	payload, err := json.Marshal(map[string]interface{}{
		"competitionId": competitionID,
		"score":         score,
	})
	if err != nil {
		return fmt.Errorf("encoding fanzone finish payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/fanzone/kick-and-win/play-finish?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("requesting fanzone play finish: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetFanZoneKickAndWinReward(ctx context.Context, competitionID, ticketID string) (*FanZoneRewardData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/kick-and-win/reward?competitionId=%s&ticketId=%s&lang=%s",
		url.QueryEscape(competitionID), url.QueryEscape(ticketID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone reward: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneRewardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone reward: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetFanZoneLeaderBoard(ctx context.Context, competitionID string) (*FanZoneLeaderBoardData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/leader-board?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone leaderboard: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneLeaderBoardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone leaderboard: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetFanZonePredictions(ctx context.Context, competitionID string) (*FanZonePredictData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/predict-and-win?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone predict: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZonePredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone predict: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetFanZoneGrandPrizes(ctx context.Context, competitionID string) (*FanZoneGrandPrizesData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/grand-prizes?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone grand prizes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneGrandPrizesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone grand prizes: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetFanZoneRewardsHistory(ctx context.Context, competitionID string) (*FanZoneRewardsHistoryData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/rewards-history?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone rewards history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneRewardsHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone rewards history: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GenerateFanZoneNickname(ctx context.Context, nickname string) (*NicknameData, error) {
	payload, err := json.Marshal(map[string]string{
		"nickname": nickname,
	})
	if err != nil {
		return nil, fmt.Errorf("encoding fanzone nickname payload: %w", err)
	}

	path := fmt.Sprintf("/api/v1/fanzone/onboarding/gen-nickname?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone nickname generation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res NicknameResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone nickname response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetFanZoneAnswerAndWin(ctx context.Context, competitionID string) (*FanZoneKickAndWinData, error) {
	path := fmt.Sprintf("/api/v1/fanzone/answer-and-win?competitionId=%s&lang=%s", url.QueryEscape(competitionID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting fanzone answer and win: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res FanZoneKickAndWinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding fanzone answer and win: %w", err)
	}
	return res.Data, nil
}

func (c *Client) PickFanZoneFavoriteTeam(ctx context.Context, competitionID, teamID string) error {
	path := fmt.Sprintf("/api/v1/fanzone/favourite-team/pick?competitionId=%s&teamId=%s&lang=%s",
		url.QueryEscape(competitionID), url.QueryEscape(teamID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return fmt.Errorf("picking favorite team: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) PickFanZoneChampionTeam(ctx context.Context, competitionID, teamID string) error {
	path := fmt.Sprintf("/api/v1/fanzone/champion-team/pick?competitionId=%s&teamId=%s&lang=%s",
		url.QueryEscape(competitionID), url.QueryEscape(teamID), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return fmt.Errorf("picking champion team: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}
