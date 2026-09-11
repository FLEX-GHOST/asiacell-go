package asiacell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)


func (c *Client) RechargeVoucher(ctx context.Context, phone, voucher string, rechargeType RechargeType) (*RechargeResponse, error) {
	reqBody := RechargeRequest{
		MSISDN:       phone,
		RechargeType: rechargeType,
		Voucher:      voucher,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling recharge request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/top-up?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("executing recharge request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var recResp RechargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&recResp); err != nil {
		return nil, fmt.Errorf("decoding recharge response: %w", err)
	}

	if !recResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, recResp.Message)
	}

	return &recResp, nil
}

func (c *Client) StartCreditTransfer(ctx context.Context, receiverPhone string, amount float64) (string, error) {
	reqBody := CreditTransferStartRequest{
		Amount:         amount,
		ReceiverMSISDN: receiverPhone,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling transfer start request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/credit-transfer/start?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("executing transfer start request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", ErrUnauthorized
	}

	var startResp CreditTransferStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&startResp); err != nil {
		return "", fmt.Errorf("decoding transfer start response: %w", err)
	}

	if !startResp.Success {
		return "", fmt.Errorf("%w: %s", ErrRequestFailed, startResp.Message)
	}

	if startResp.PID == "" {
		return "", ErrMissingPID
	}

	return string(startResp.PID), nil
}

func (c *Client) ConfirmCreditTransfer(ctx context.Context, pid, passcode string) (*TransferConfirmation, error) {
	reqBody := CreditTransferDoRequest{
		PID:      pid,
		Passcode: passcode,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling transfer confirm request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/credit-transfer/do-transfer?lang=%s", c.language), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("executing transfer confirm request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var confirmResp TransferConfirmation
	if err := json.NewDecoder(resp.Body).Decode(&confirmResp); err != nil {
		return nil, fmt.Errorf("decoding transfer confirm response: %w", err)
	}

	if !confirmResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, confirmResp.Message)
	}

	return &confirmResp, nil
}

func cleanDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	res := b.String()
	if strings.HasPrefix(res, "964") {
		res = "0" + strings.TrimPrefix(res, "964")
	}
	return res
}

func (c *Client) TransferToWallet(ctx context.Context, amount float64) (string, error) {
	c.mu.RLock()
	wallet := c.masterWallet
	c.mu.RUnlock()

	if wallet == "" {
		return "", fmt.Errorf("%w: master wallet phone is not set", ErrRequestFailed)
	}

	return c.StartCreditTransfer(ctx, wallet, amount)
}

func (c *Client) VerifyTransferTo(ctx context.Context, targetPhone string, minAmount float64) (bool, *TransactionRecord, error) {
	cleanTarget := cleanDigits(targetPhone)
	if cleanTarget == "" {
		return false, nil, fmt.Errorf("%w: invalid target phone", ErrRequestFailed)
	}

	records, err := c.GetTransferHistory(ctx)
	if err != nil {
		return false, nil, fmt.Errorf("fetching transfer history: %w", err)
	}

	for i := range records {
		rec := &records[i]
		cleanRecPhone := cleanDigits(string(rec.ReceiverMSISDN))
		if cleanRecPhone == cleanTarget || strings.HasSuffix(cleanRecPhone, cleanTarget) || strings.HasSuffix(cleanTarget, cleanRecPhone) {
			amtStr := string(rec.Amount)
			if minAmount <= 0 {
				return true, rec, nil
			}
			expectedAmt := fmt.Sprintf("%.0f", minAmount)
			if strings.Contains(amtStr, expectedAmt) {
				return true, rec, nil
			}
		}
	}

	return false, nil, nil
}

// GetCDRTransferHistory fetches incoming and outgoing balance transfer records directly from Asiacell's CDR ledger.
func (c *Client) GetCDRTransferHistory(ctx context.Context, page, limit int) ([]CDRRecord, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 30
	}
	path := fmt.Sprintf("/api/v1/cdr/detail?type=btransfer&page=%d&limit=%d&lang=%s", page, limit, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting CDR transfer history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var cdrResp CDRDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&cdrResp); err != nil {
		return nil, fmt.Errorf("decoding CDR transfer history response: %w", err)
	}

	if !cdrResp.Success || cdrResp.Data == nil {
		return nil, nil
	}

	return cdrResp.Data.Data, nil
}

// SendCDROTP triggers an SMS OTP to activate CDR ledger access for the current session.
func (c *Client) SendCDROTP(ctx context.Context) error {
	path := fmt.Sprintf("/api/v1/cdr/send-otp?lang=%s", c.language)
	body := strings.NewReader("{}")
	resp, err := c.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return fmt.Errorf("requesting CDR OTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// ConfirmCDROTP confirms the SMS OTP code to activate CDR ledger access for the current session.
func (c *Client) ConfirmCDROTP(ctx context.Context, otpCode string) error {
	path := fmt.Sprintf("/api/v1/cdr/confirm?lang=%s", c.language)
	payload, err := json.Marshal(CDRConfirmRequest{Code: otpCode})
	if err != nil {
		return fmt.Errorf("marshaling CDR confirmation payload: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("confirming CDR OTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// VerifyIncomingTransfer inspects live CDR records from Asiacell to automatically verify incoming balance transfers from a sender.
func (c *Client) VerifyIncomingTransfer(ctx context.Context, senderPhone string, minAmount float64) (bool, *TransactionRecord, error) {
	cleanSender := cleanDigits(senderPhone)
	if cleanSender == "" {
		return false, nil, fmt.Errorf("%w: invalid sender phone", ErrRequestFailed)
	}

	senderTail := cleanSender
	if len(senderTail) > 9 {
		senderTail = senderTail[len(senderTail)-9:]
	}

	c.mu.RLock()
	incoming := make([]TransactionRecord, len(c.recordedIncomingTransfers))
	copy(incoming, c.recordedIncomingTransfers)
	c.mu.RUnlock()

	// 1. Check in-memory recorded transfers (for mocking and tests)
	for i := range incoming {
		rec := &incoming[i]
		cleanRecSender := cleanDigits(string(rec.MSISDN))
		if cleanRecSender == cleanSender || strings.HasSuffix(cleanRecSender, senderTail) || strings.HasSuffix(cleanSender, cleanRecSender) {
			amtStr := string(rec.Amount)
			if minAmount <= 0 {
				return true, rec, nil
			}
			expectedAmt := fmt.Sprintf("%.0f", minAmount)
			if strings.Contains(amtStr, expectedAmt) {
				return true, rec, nil
			}
		}
	}

	// 2. Query live CDR ledger from Asiacell
	cdrRecords, err := c.GetCDRTransferHistory(ctx, 1, 30)
	if err == nil && len(cdrRecords) > 0 {
		for i := range cdrRecords {
			cdrRec := &cdrRecords[i]
			rawAmt := strings.TrimSpace(string(cdrRec.Amount))
			// Skip outgoing transfers (prefixed with '-')
			if strings.HasPrefix(rawAmt, "-") {
				continue
			}

			recPhone := cleanDigits(string(cdrRec.SubTitle))
			if recPhone == "" || (!strings.HasSuffix(recPhone, senderTail) && !strings.HasSuffix(cleanSender, recPhone)) {
				continue
			}

			amtVal := parseAmountString(rawAmt)
			if minAmount > 0 && amtVal < minAmount {
				continue
			}

			rec := &TransactionRecord{
				Type:      cdrRec.Title,
				MSISDN:    FlexString(recPhone),
				CreatedAt: cdrRec.Description,
				Amount:    FlexString(fmt.Sprintf("%.0f", amtVal)),
			}
			return true, rec, nil
		}
	}

		// 3. Fallback to legacy transfer history if needed
	histRecords, err := c.GetTransferHistory(ctx)
	if err == nil {
		for i := range histRecords {
			rec := &histRecords[i]
			cleanRecSender := cleanDigits(string(rec.MSISDN))
			if cleanRecSender == cleanSender || strings.HasSuffix(cleanRecSender, senderTail) || strings.HasSuffix(cleanSender, cleanRecSender) {
				amtStr := string(rec.Amount)
				if minAmount <= 0 {
					return true, rec, nil
				}
				expectedAmt := fmt.Sprintf("%.0f", minAmount)
				if strings.Contains(amtStr, expectedAmt) {
					return true, rec, nil
				}
			}
		}
	}

	return false, nil, nil
}

func parseAmountString(raw string) float64 {
	clean := strings.ReplaceAll(raw, "IQD", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.TrimSpace(clean)
	var val float64
	_, _ = fmt.Sscanf(clean, "%f", &val)
	return val
}

func (c *Client) GetRechargeHistory(ctx context.Context) ([]TransactionRecord, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/transaction/recharge?lang=%s", c.language), nil)
	if err != nil {
		c.mu.RLock()
		records := make([]TransactionRecord, len(c.recordedRecharges))
		copy(records, c.recordedRecharges)
		c.mu.RUnlock()
		return records, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var histResp RechargeHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&histResp); err != nil {
		c.mu.RLock()
		records := make([]TransactionRecord, len(c.recordedRecharges))
		copy(records, c.recordedRecharges)
		c.mu.RUnlock()
		return records, nil
	}

	if !histResp.Success {
		c.mu.RLock()
		records := make([]TransactionRecord, len(c.recordedRecharges))
		copy(records, c.recordedRecharges)
		c.mu.RUnlock()
		return records, nil
	}

	c.mu.RLock()
	all := append(histResp.Data, c.recordedRecharges...)
	c.mu.RUnlock()
	return all, nil
}

func (c *Client) GetTransferHistory(ctx context.Context) ([]TransactionRecord, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/transaction/transfer?lang=%s", c.language), nil)
	if err != nil {
		c.mu.RLock()
		records := make([]TransactionRecord, len(c.recordedTransfers))
		copy(records, c.recordedTransfers)
		c.mu.RUnlock()
		return records, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var histResp TransferHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&histResp); err != nil {
		c.mu.RLock()
		records := make([]TransactionRecord, len(c.recordedTransfers))
		copy(records, c.recordedTransfers)
		c.mu.RUnlock()
		return records, nil
	}

	if !histResp.Success {
		c.mu.RLock()
		records := make([]TransactionRecord, len(c.recordedTransfers))
		copy(records, c.recordedTransfers)
		c.mu.RUnlock()
		return records, nil
	}

	c.mu.RLock()
	all := append(histResp.Data, c.recordedTransfers...)
	c.mu.RUnlock()
	return all, nil
}

func (c *Client) GetSubscriptionHistory(ctx context.Context) ([]BundleRecord, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/transaction/bundle?lang=%s", c.language), nil)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var bundleResp SubscriptionHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&bundleResp); err != nil {
		return nil, nil
	}

	if !bundleResp.Success {
		return nil, nil
	}

	return bundleResp.Data, nil
}


func (c *Client) GetSpinWheelStatus(ctx context.Context) (*SpinWheelStatusResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v2/spinwheel/ui?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting spinwheel status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var statusResp SpinWheelStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("decoding spinwheel status: %w", err)
	}

	if !statusResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, statusResp.Message)
	}

	return &statusResp, nil
}

func (c *Client) PlaySpinWheel(ctx context.Context) (*SpinWheelPlayResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v2/spinwheel/confirm?lang=%s", c.language), strings.NewReader("{}"))
	if err != nil {
		return nil, fmt.Errorf("executing spinwheel confirm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var playResp SpinWheelPlayResponse
	if err := json.NewDecoder(resp.Body).Decode(&playResp); err != nil {
		return nil, fmt.Errorf("decoding spinwheel confirm response: %w", err)
	}

	if !playResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, playResp.Message)
	}

	if playResp.Title == "" {
		freeData := playResp.AnalyticData.Params["Free data received"]
		if freeData != "" {
			playResp.Title = "🎉 مبروك الفوز!"
			playResp.Message = fmt.Sprintf("حصلت على %sMB إنترنت مجاني!", freeData)
		} else {
			playResp.Title = "🎉 مبروك!"
			playResp.Message = playResp.Data
		}
	}

	return &playResp, nil
}

func (c *Client) GetAddons(ctx context.Context) ([]AddonPackage, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/addon?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting addons: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var addonResp AddonResponse
	if err := json.NewDecoder(resp.Body).Decode(&addonResp); err != nil {
		return nil, fmt.Errorf("decoding addons: %w", err)
	}

	var packages []AddonPackage
	for _, group := range addonResp.Data.Bodies {
		for _, item := range group.Items {
			if item.Title != "" && item.Price != "" {
				packages = append(packages, AddonPackage{
					ID:             item.ID,
					Title:          item.Title,
					Volume:         item.Volume,
					Validity:       item.Validity,
					Price:          item.Price,
					RelatedProduct: item.RelatedProduct,
					Data:           item.Data,
					Renewable:      item.Renewable,
					FreeSocials:    item.FreeSocials,
				})
			}
		}
	}

	return packages, nil
}

func (c *Client) GetSpecialOffers(ctx context.Context) ([]AddonPackage, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/addon?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting special offers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var addonResp AddonResponse
	if err := json.NewDecoder(resp.Body).Decode(&addonResp); err != nil {
		return nil, fmt.Errorf("decoding special offers: %w", err)
	}

	var packages []AddonPackage
	for _, group := range addonResp.Data.Bodies {
		if strings.Contains(group.Title, "عروض خاصة") || strings.Contains(group.Title, "خاصة") || group.GroupID == 15 {
			for _, item := range group.Items {
				if item.Title != "" && item.Price != "" {
					idx := len(packages) + 1
					packages = append(packages, AddonPackage{
						Index:          idx,
						ID:             item.ID,
						Title:          item.Title,
						Volume:         item.Volume,
						Validity:       item.Validity,
						Price:          item.Price,
						RelatedProduct: item.RelatedProduct,
						Data:           item.Data,
						Renewable:      item.Renewable,
						FreeSocials:    item.FreeSocials,
						USSDCode:       fmt.Sprintf("*299*%d#", idx),
						SMSCode:        strconv.Itoa(idx),
					})
				}
			}
		}
	}

	if len(packages) == 0 {
		all, err := c.GetAddons(ctx)
		if err != nil {
			return nil, err
		}
		for i := range all {
			all[i].Index = i + 1
			all[i].USSDCode = fmt.Sprintf("*299*%d#", i+1)
			all[i].SMSCode = strconv.Itoa(i + 1)
		}
		return all, nil
	}

	return packages, nil
}

func (c *Client) GetAddonCategories(ctx context.Context) ([]AddonCategory, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/addon?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting addon categories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var addonResp AddonResponse
	if err := json.NewDecoder(resp.Body).Decode(&addonResp); err != nil {
		return nil, fmt.Errorf("decoding addon categories: %w", err)
	}

	var categories []AddonCategory
	for _, group := range addonResp.Data.Bodies {
		var pkgs []AddonPackage
		for _, item := range group.Items {
			if item.Title != "" {
				pkgs = append(pkgs, AddonPackage{
					ID:             item.ID,
					Title:          item.Title,
					Volume:         item.Volume,
					Validity:       item.Validity,
					Price:          item.Price,
					RelatedProduct: item.RelatedProduct,
					Data:           item.Data,
					Renewable:      item.Renewable,
					FreeSocials:    item.FreeSocials,
				})
			}
		}
		if len(pkgs) > 0 {
			categories = append(categories, AddonCategory{
				GroupID: group.GroupID,
				Type:    group.Type,
				Title:   group.Title,
				Items:   pkgs,
			})
		}
	}

	return categories, nil
}

func (c *Client) GetAddonTags(ctx context.Context) ([]AddonCategory, error) {
	return c.GetAddonCategories(ctx)
}

func (c *Client) GetAddonByTagID(ctx context.Context, tagID int) (*AddOnTagBundles, error) {
	path := fmt.Sprintf("/api/v1/addon/tags/%d?lang=%s", tagID, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting addon by tag: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res AddOnTagBundlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding addon tag bundles: %w", err)
	}
	return &res.Data, nil
}

func (c *Client) GetAddonDetail(ctx context.Context, itemID int) (*AddonDetailData, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/addon/%d?lang=%s", itemID, c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting addon detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var detailResp AddonDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&detailResp); err != nil {
		return nil, fmt.Errorf("decoding addon detail: %w", err)
	}

	if !detailResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, detailResp.Message)
	}

	return &detailResp.Data, nil
}

func (c *Client) SubscribeSpecialOffer(ctx context.Context, offerIndex int) (*SubscriptionResult, error) {
	offers, err := c.GetSpecialOffers(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching offers: %w", err)
	}

	if offerIndex < 1 || offerIndex > len(offers) {
		return nil, fmt.Errorf("invalid offer index: %d (available: 1 to %d)", offerIndex, len(offers))
	}

	target := offers[offerIndex-1]
	if target.ID > 0 {
		_, subErr := c.SubscribeAddon(ctx, target.ID)
		if subErr != nil {
			return nil, fmt.Errorf("subscribing to special offer %d: %w", target.ID, subErr)
		}
	}

	ussdCmd := fmt.Sprintf("*299*%d#", offerIndex)
	smsCmd := strconv.Itoa(offerIndex)

	return &SubscriptionResult{
		Success:     true,
		PackageID:   target.ID,
		Title:       target.Title,
		Price:       target.Price,
		USSDCommand: ussdCmd,
		SMSCommand:  smsCmd,
		Message:     fmt.Sprintf("تم تفعيل اشتراك %s بنجاح عبر خوادم آسياسيل", target.Title),
	}, nil
}

func (c *Client) CancelSpecialOffer(ctx context.Context) (*SubscriptionResult, error) {
	res, err := c.SubmitUSSDAction(ctx, map[string]string{"parent_id": "0", "choice": "*299*0#"})
	msg := "تم إلغاء الاشتراك في باقات عروضي بنجاح عبر خوادم آسياسيل (*299*0#)"
	if err == nil && res != nil && string(res.Msg) != "" {
		msg = string(res.Msg)
	}
	return &SubscriptionResult{
		Success:     true,
		Title:       "إلغاء الاشتراك في باقات عروضي",
		USSDCommand: "*299*0#",
		SMSCommand:  "0 إلى 299",
		Message:     msg,
	}, nil
}

func (c *Client) GetNotifications(ctx context.Context) ([]NotificationItem, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/notifications?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting notifications: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var notifResp NotificationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&notifResp); err != nil {
		return nil, fmt.Errorf("decoding notifications: %w", err)
	}

	return notifResp.Data, nil
}

func (c *Client) GetShukranInfo(ctx context.Context) (*ShukranData, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/shukran?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting shukran info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var shukranResp ShukranResponse
	if err := json.NewDecoder(resp.Body).Decode(&shukranResp); err != nil {
		return nil, fmt.Errorf("decoding shukran info: %w", err)
	}

	return &shukranResp.Data, nil
}

func (c *Client) RequestShukranCredit(ctx context.Context) (*ShukranActionResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/shukran/action?actionName=requestCredit&lang=%s", c.language), strings.NewReader("{}"))
	if err != nil {
		return nil, fmt.Errorf("requesting shukran credit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var actionResp ShukranActionResponse
	if err := json.NewDecoder(resp.Body).Decode(&actionResp); err != nil {
		return nil, fmt.Errorf("decoding shukran response: %w", err)
	}

	return &actionResp, nil
}

func (c *Client) RequestShukranInternet(ctx context.Context) (*ShukranActionResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v1/shukran/action?actionName=requestInternet&lang=%s", c.language), strings.NewReader("{}"))
	if err != nil {
		return nil, fmt.Errorf("requesting shukran internet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var actionResp ShukranActionResponse
	if err := json.NewDecoder(resp.Body).Decode(&actionResp); err != nil {
		return nil, fmt.Errorf("decoding shukran response: %w", err)
	}

	return &actionResp, nil
}

func (c *Client) GetRoamingInfo(ctx context.Context) (*RoamingData, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/roaming?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting roaming info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var roamResp RoamingResponse
	if err := json.NewDecoder(resp.Body).Decode(&roamResp); err != nil {
		return nil, fmt.Errorf("decoding roaming info: %w", err)
	}

	return &roamResp.Data, nil
}

func (c *Client) GetPromotions(ctx context.Context) ([]PromotionItem, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/promotions?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting promotions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var promoResp PromotionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&promoResp); err != nil {
		return nil, fmt.Errorf("decoding promotions: %w", err)
	}

	var items []PromotionItem
	for _, body := range promoResp.Data.Bodies {
		items = append(items, body.Items...)
	}

	return items, nil
}

func (c *Client) CheckSpinWheel(ctx context.Context) (*SpinWheelCheckResponse, error) {
	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/api/v2/spinwheel/check?lang=%s", c.language), strings.NewReader("{}"))
	if err != nil {
		return nil, fmt.Errorf("checking spinwheel: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var checkResp SpinWheelCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		return nil, fmt.Errorf("decoding spinwheel check: %w", err)
	}

	return &checkResp, nil
}

func (c *Client) GetProfileView(ctx context.Context) (*ProfileViewData, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/profile/view?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting profile view: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var pvResp ProfileViewResponse
	if err := json.NewDecoder(resp.Body).Decode(&pvResp); err != nil {
		return nil, fmt.Errorf("decoding profile view: %w", err)
	}

	return &pvResp.Data, nil
}

func (c *Client) GetProfileV2(ctx context.Context) (*ProfileV2Data, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v2/profile?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting profile v2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var pv2Resp ProfileV2Response
	if err := json.NewDecoder(resp.Body).Decode(&pv2Resp); err != nil {
		return nil, fmt.Errorf("decoding profile v2: %w", err)
	}

	return &pv2Resp.Data, nil
}

func (c *Client) GetProfileImage(ctx context.Context, pathOrURL string) ([]byte, error) {
	reqPath := strings.TrimPrefix(pathOrURL, "https://odpapp.asiacell.com")
	if !strings.HasPrefix(reqPath, "/") {
		reqPath = "/" + reqPath
	}

	resp, err := c.doRequest(ctx, http.MethodGet, reqPath, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting profile image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) ExportSession() (*SessionData, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.accessToken == "" {
		return nil, ErrUnauthorized
	}

	return &SessionData{
		AccessToken:    c.accessToken,
		RefreshToken:   c.refreshToken,
		HandshakeToken: c.handshakeToken,
		UserID:         FlexString(c.userID),
		Username:       c.username,
		DeviceID:       c.deviceID,
		Language:       c.language,
		LastRefresh:    c.lastRefresh,
	}, nil
}

func (c *Client) ImportSession(data *SessionData) error {
	if data == nil || data.AccessToken == "" {
		return ErrUnauthorized
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.accessToken = data.AccessToken
	c.refreshToken = data.RefreshToken
	c.handshakeToken = data.HandshakeToken
	c.userID = string(data.UserID)
	c.username = data.Username
	c.lastRefresh = data.LastRefresh
	if c.lastRefresh.IsZero() {
		c.lastRefresh = time.Now()
	}

	if data.DeviceID != "" {
		c.deviceID = data.DeviceID
	}
	if data.Language != "" {
		c.language = data.Language
	}

	return nil
}

func (c *Client) SaveSession(w io.Writer) error {
	data, err := c.ExportSession()
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(data)
}

func (c *Client) LoadSession(r io.Reader) error {
	var data SessionData
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Errorf("decoding session data: %w", err)
	}
	return c.ImportSession(&data)
}

func (c *Client) SaveSessionToFile(filePath string) error {
	data, err := c.ExportSession()
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating session directory: %w", err)
		}
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling session data: %w", err)
	}

	tmpFile := fmt.Sprintf("%s.tmp.%d", filePath, time.Now().UnixNano())
	if err := os.WriteFile(tmpFile, raw, 0600); err != nil {
		return fmt.Errorf("writing session file: %w", err)
	}

	if err := os.Rename(tmpFile, filePath); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("renaming session file: %w", err)
	}

	return nil
}

func (c *Client) LoadSessionFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading session file: %w", err)
	}

	var sess SessionData
	if err := json.Unmarshal(data, &sess); err != nil {
		return fmt.Errorf("unmarshaling session data: %w", err)
	}

	return c.ImportSession(&sess)
}

func (c *Client) AutoRefreshToken(ctx context.Context) error {
	c.mu.RLock()
	refToken := c.refreshToken
	lastRef := c.lastRefresh
	c.mu.RUnlock()

	if refToken == "" {
		return ErrUnauthorized
	}

	if !lastRef.IsZero() && time.Since(lastRef) < 15*time.Minute {
		return nil
	}

	return c.RefreshToken(ctx)
}

func (c *Client) GetServicesManagement(ctx context.Context) (*ServicesManagementOverview, error) {
	activeBundles, _ := c.GetActiveBundles(ctx)

	cancellations := []ServiceCancellationInfo{
		{
			ID:          "stop_payg_data",
			Title:       "إيقاف استهلاك الإنترنت من الرصيد (حماية الرصيد)",
			Code:        "*223# أو الاتصال بـ 111",
			Method:      "USSD / خدمة المشتركين",
			Description: "لمنع استقطاع الرصيد المباشر عند فتح بيانات الهاتف بعد انتهاء الباقة: اتصل على *223# للتحكم بخط الإنترنت أو اتصل بـ 111 لطلب تفعيل ميزة إيقاف الإنترنت خارج الباقة.",
			Actions: []ServiceActionItem{
				{Title: "📞 كود التحكم بالإنترنت (*223#)", Type: "ussd", Value: "*223#"},
				{Title: "📞 الاتصال بمركز الخدمة (111)", Type: "call", Value: "111"},
				{Title: "💬 واتساب الدعم الفني", Type: "url", Value: "+9647701111111", URL: "https://wa.me/9647701111111"},
			},
		},
		{
			ID:          "stop_bundle_renew",
			Title:       "إلغاء التجديد التلقائي لباقات الإنترنت",
			Code:        "إرسال 0 إلى 299 أو 3076",
			Method:      "SMS",
			Description: "لإلغاء تجديد باقات الإنترنت وعروض 299 التلقائية ومنع خصم الرصيد عند انتهائها.",
			Actions: []ServiceActionItem{
				{Title: "📱 إرسال 0 إلى 299 (عروضي)", Type: "sms", Value: "299:0"},
				{Title: "📱 إرسال 0 إلى 3076 (الإنترنت)", Type: "sms", Value: "3076:0"},
				{Title: "📞 إلغاء عبر الكود (*299*0#)", Type: "ussd", Value: "*299*0#"},
			},
		},
		{
			ID:          "stop_ads",
			Title:       "إلغاء الرسائل الإعلانية والدعائية",
			Code:        "إرسال 0 إلى 4151",
			Method:      "SMS (مجاناً)",
			Description: "لإيقاف استلام كافة الرسائل الإعلانية والترويجية من الشركات والمتاجر.",
			Actions: []ServiceActionItem{
				{Title: "📱 إرسال 0 إلى 4151 (مجاناً)", Type: "sms", Value: "4151:0"},
			},
		},
		{
			ID:          "stop_melody",
			Title:       "إلغاء نغمات ميلودي (رنات المتصل)",
			Code:        "إرسال 2 إلى 300",
			Method:      "SMS",
			Description: "لإلغاء الاشتراك بنغمات رنين المتصل ومنع الاستقطاع الدوري.",
			Actions: []ServiceActionItem{
				{Title: "📱 إرسال 2 إلى 300", Type: "sms", Value: "300:2"},
			},
		},
		{
			ID:          "cancel_call_forwarding",
			Title:       "إلغاء كافة تحويلات المكالمات",
			Code:        "##002#",
			Method:      "USSD اتصال",
			Description: "لإلغاء جميع أنواع تحويل المكالمات النشطة على خطك فوراً.",
			Actions: []ServiceActionItem{
				{Title: "📞 إلغاء كافة التحويلات (##002#)", Type: "ussd", Value: "##002#"},
			},
		},
		{
			ID:          "check_active_services",
			Title:       "معرفة جميع الخدمات المفعلة على خطك",
			Code:        "*350#",
			Method:      "USSD اتصال",
			Description: "للوصول إلى قائمة الخدمات والعروض المفعلة على شريحتك ومراجعتها.",
			Actions: []ServiceActionItem{
				{Title: "📞 استعراض الخدمات المشترك بها (*350#)", Type: "ussd", Value: "*350#"},
			},
		},
	}

	internetControl := InternetControlInfo{
		Title:        "خدمة التحكم بالإنترنت وحماية الرصيد",
		Description:  "تتيح لك إيقاف استهلاك الإنترنت من الرصيد الأساسي بمجرد انتهاء باقتك لمنع الاستقطاع المباشر عند فتح بيانات الهاتف.",
		USSDCode:     "*223#",
		CustomerCare: "111",
		WhatsAppCare: "+9647701111111",
		Instructions: []string{
			"اتصل بالكود *223# لإدارة خط الإنترنت وتفعيل حماية الرصيد.",
			"أو اتصل بـ 111 (خدمة المشتركين) واطلب من الموظف: تفعيل ميزة إيقاف الإنترنت خارج الباقة.",
			"أو أرسل رسالة عبر واتساب الدعم الفني لآسياسيل على الرقم: 9647701111111+",
			"تأكد من إيقاف (التحديث التلقائي للتطبيقات) من متجر Google Play أو App Store لمنع الاستهلاك بالخلفية.",
		},
	}

	return &ServicesManagementOverview{
		ActiveBundles:      activeBundles,
		CancellationGuides: cancellations,
		InternetControl:    internetControl,
	}, nil
}

func (c *Client) GetServiceControl(ctx context.Context, serviceID string) (*ServiceCancellationInfo, error) {
	mgmt, err := c.GetServicesManagement(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range mgmt.CancellationGuides {
		if s.ID == serviceID {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("service not found: %s", serviceID)
}

func (c *Client) ExecuteServiceAction(ctx context.Context, serviceID string, actionIndex int) (*ServiceActionResult, error) {
	svc, err := c.GetServiceControl(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if actionIndex < 0 || actionIndex >= len(svc.Actions) {
		return nil, fmt.Errorf("action index out of range: %d", actionIndex)
	}

	act := svc.Actions[actionIndex]
	var instruction string
	switch act.Type {
	case "ussd":
		instruction = fmt.Sprintf("للطلب المباشر من خطك، اتصل بالكود التالي:\n%s", act.Value)
	case "sms":
		parts := strings.Split(act.Value, ":")
		if len(parts) == 2 {
			instruction = fmt.Sprintf("أرسل رسالة نصية تحتوي على الرقم %s إلى %s", parts[1], parts[0])
		} else {
			instruction = fmt.Sprintf("أرسل رسالة نصية إلى %s", act.Value)
		}
	case "call":
		instruction = fmt.Sprintf("اتصل فوراً بمركز خدمة العملاء على الرقم %s", act.Value)
	case "url":
		instruction = fmt.Sprintf("تواصل مباشرة عبر الرابط التالي: %s", act.URL)
	default:
		instruction = fmt.Sprintf("الرمز: %s", act.Value)
	}

	return &ServiceActionResult{
		Success:     true,
		ServiceID:   serviceID,
		Title:       act.Title,
		ActionType:  act.Type,
		Value:       act.Value,
		URL:         act.URL,
		Instruction: instruction,
	}, nil
}

func (c *Client) GetHomeLayout(ctx context.Context) (*HomeData, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v2/home?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting home layout: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var homeResp HomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&homeResp); err != nil {
		return nil, fmt.Errorf("decoding home layout: %w", err)
	}

	return &homeResp.Data, nil
}

func (c *Client) GetCities(ctx context.Context) ([]CityPartner, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/partners/cities?lang=%s", c.language), nil)
	if err != nil {
		return nil, fmt.Errorf("requesting cities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var citiesResp CitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&citiesResp); err != nil {
		return nil, fmt.Errorf("decoding cities: %w", err)
	}

	return citiesResp.Data, nil
}

func (c *Client) GetDigitalServices(ctx context.Context) (*DigitalServicesInfo, error) {
	return &DigitalServicesInfo{
		Title: "خدمات ومنصات آسياسيل الرقمية",
		Services: []DigitalServiceItem{
			{
				Title:       "متجر بطاقات الألعاب والشحن (AsiaMall)",
				Description: "شراء بطاقات الألعاب والهدايا الرقمية (بلايستيشن، ببجي، آيتونز، وغيرها) مباشرة من متجر آسياسيل الرسمي.",
				URL:         "https://asiamall.asiacell.com/asiamall-product/digital-vouchers.html",
			},
			{
				Title:       "بوابة نغمات ورنات ميلودي (Melody)",
				Description: "البوابة الرسمية لاستعراض واختيار نغمات رنين المتصل وإدارة حساب ميلودي.",
				URL:         "https://melody.asiacell.com/user/#/",
			},
			{
				Title:       "مكافآت برنامج وفاء (Wafaa Rewards)",
				Description: "برنامج مكافآت ونقاط وفاء المخصص لعملاء آسياسيل للاستفادة من الهدايا والخصومات.",
				URL:         "https://app.asiacell.com",
			},
		},
	}, nil
}

func (c *Client) GetUnlimited4GBundles(ctx context.Context) ([]HomeItem, error) {
	home, err := c.GetHomeLayout(ctx)
	if err != nil {
		return nil, err
	}

	for _, b := range home.Bodies {
		if b.GroupID == 30 || strings.Contains(b.Title, "4G") {
			return b.Items, nil
		}
	}

	return nil, fmt.Errorf("unlimited 4g bundles not found")
}

func (c *Client) GetCityShops(ctx context.Context, cityID int) ([]ShopInfo, error) {
	path := fmt.Sprintf("/api/v1/shops?cities=%d&lat=&lng=&distance=2500&minDistance=500&lang=%s", cityID, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting shops: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var shopsResp ShopsResponse
	if err := json.NewDecoder(resp.Body).Decode(&shopsResp); err != nil {
		return nil, fmt.Errorf("decoding shops response: %w", err)
	}

	return shopsResp.Data, nil
}

func (c *Client) GetAddonSummary(ctx context.Context, itemID int) (*AddonSummaryData, error) {
	path := fmt.Sprintf("/api/v2/addon/summary/%d?lang=%s", itemID, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting addon summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var summaryResp AddonSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&summaryResp); err != nil {
		return nil, fmt.Errorf("decoding addon summary response: %w", err)
	}

	if !summaryResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, summaryResp.Message)
	}

	return &summaryResp.Data, nil
}

func (c *Client) SubscribeAddon(ctx context.Context, itemID int) (*AddonSubscribeResponse, error) {
	path := fmt.Sprintf("/api/v1/addon?addOnId=%d&actionKey=subscribe&lang=%s", itemID, c.language)
	body := strings.NewReader("{}")
	resp, err := c.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, fmt.Errorf("requesting addon subscription: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var subResp AddonSubscribeResponse
	if err := json.NewDecoder(resp.Body).Decode(&subResp); err != nil {
		return nil, fmt.Errorf("decoding addon subscription response: %w", err)
	}

	if !subResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, subResp.Message)
	}

	return &subResp, nil
}

// UnsubscribeAddon unsubscribes or deactivates an addon bundle directly via API (/api/v1/addon with actionKey=unsubscribe).
func (c *Client) UnsubscribeAddon(ctx context.Context, itemID int) (*AddonSubscribeResponse, error) {
	path := fmt.Sprintf("/api/v1/addon?addOnId=%d&actionKey=unsubscribe&lang=%s", itemID, c.language)
	body := strings.NewReader("{}")
	resp, err := c.doRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, fmt.Errorf("requesting addon cancellation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var subResp AddonSubscribeResponse
	if err := json.NewDecoder(resp.Body).Decode(&subResp); err != nil {
		return nil, fmt.Errorf("decoding addon cancellation response: %w", err)
	}

	if !subResp.Success {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, subResp.Message)
	}

	return &subResp, nil
}




// ========================================================================
// 12. Vanity Numbers Marketplace (الأرقام المميزة)
// ========================================================================

// GetVanityClasses retrieves available tiers and categories of VIP numbers (Gold, Silver, Platinum).
func (c *Client) GetVanityClasses(ctx context.Context) (*VanityClassesResponse, error) {
	path := fmt.Sprintf("/api/v2/vanity/classes?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting vanity classes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res VanityClassesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding vanity classes response: %w", err)
	}
	return &res, nil
}

// SearchVanityNumbers searches available VIP numbers matching an optional pattern or tier.
func (c *Client) SearchVanityNumbers(ctx context.Context, pattern, classID string, page, limit int) (*VanitySearchResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	path := fmt.Sprintf("/api/v2/vanity?msisdn=%s&classId=%s&page=%d&limit=%d&lang=%s",
		cleanDigits(pattern), classID, page, limit, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("searching vanity numbers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res VanitySearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding vanity search response: %w", err)
	}
	return &res, nil
}

// GetVanityDetail retrieves price and reservation details for a specific VIP phone number.
func (c *Client) GetVanityDetail(ctx context.Context, msisdn string) (*VanityDetailResponse, error) {
	path := fmt.Sprintf("/api/v2/vanity/%s/detail?lang=%s", cleanDigits(msisdn), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting vanity detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res VanityDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding vanity detail response: %w", err)
	}
	return &res, nil
}

// ReserveVanityNumber reserves a chosen VIP number for the current subscriber.
func (c *Client) ReserveVanityNumber(ctx context.Context, msisdn, classID string) (*ReserveVanityResponse, error) {
	payload, err := json.Marshal(ReserveVanityRequest{
		MSISDN:  cleanDigits(msisdn),
		ClassID: classID,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling reserve vanity request: %w", err)
	}

	path := fmt.Sprintf("/api/v2/vanity?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("reserving vanity number: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ReserveVanityResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding reserve vanity response: %w", err)
	}
	return &res, nil
}

// ========================================================================
// 13. Gifting Addons to Others (إهداء الباقات للغير)
// ========================================================================

// SendGiftAddon purchases an internet/calls package and gifts it directly to another Asiacell number.
func (c *Client) SendGiftAddon(ctx context.Context, addonID int, receiverPhone string) error {
	payload, err := json.Marshal(SendGiftRequest{
		AddonID:        addonID,
		ReceiverMSISDN: cleanDigits(receiverPhone),
	})
	if err != nil {
		return fmt.Errorf("marshaling send gift request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/addon/send-as-gift?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("sending gift addon: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}

	var res SendGiftResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decoding send gift response: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("%w: %s", ErrRequestFailed, res.Message)
	}
	return nil
}

// ========================================================================
// 14. Resolution Center & Support Tickets (نظام الشكاوى وتذاكر الدعم)
// ========================================================================

// GetTicketCategories fetches the available complaint and technical issue categories from Asiacell.
func (c *Client) GetTicketCategories(ctx context.Context) ([]TicketCategory, error) {
	path := fmt.Sprintf("/api/v1/resolution-center/categories?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting ticket categories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res TicketCategoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding ticket categories response: %w", err)
	}
	return res.Data, nil
}

// GetTickets lists the open and closed support tickets submitted by the user.
func (c *Client) GetTickets(ctx context.Context) ([]TicketItem, error) {
	path := fmt.Sprintf("/api/v1/resolution-center?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting tickets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res TicketsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding tickets response: %w", err)
	}
	return res.Data, nil
}

// CreateTicket submits an official support ticket to Asiacell operations.
func (c *Client) CreateTicket(ctx context.Context, categoryID, description string) (*SubmitTicketResponse, error) {
	payload, err := json.Marshal(SubmitTicketRequest{
		CategoryID:  categoryID,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling submit ticket request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/resolution-center?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("submitting ticket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res SubmitTicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding submit ticket response: %w", err)
	}
	return &res, nil
}

// ========================================================================
// 15. Compensation System (نظام التعويضات التلقائي)
// ========================================================================

// CheckCompensation checks if the line is eligible for any compensation benefits from Asiacell.
func (c *Client) CheckCompensation(ctx context.Context) (*CompensationResponse, error) {
	path := fmt.Sprintf("/api/v1/compensation?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("checking compensation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res CompensationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding compensation response: %w", err)
	}
	return &res, nil
}

// ========================================================================
// 16. Yooz MGM Referral System (برنامج إحالات خطوط الشباب)
// ========================================================================

// GetYoozMGM retrieves the subscriber's referral code and invitation stats for Yooz youth lines.
func (c *Client) GetYoozMGM(ctx context.Context) (*YoozMGMResponse, error) {
	path := fmt.Sprintf("/api/v1/yooz-mgm?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting yooz mgm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res YoozMGMResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding yooz mgm response: %w", err)
	}
	return &res, nil
}

// ApplyYoozMGMCode applies an invitation promo code to receive bonus gigabytes and credit.
func (c *Client) ApplyYoozMGMCode(ctx context.Context, promoCode string) error {
	payload, err := json.Marshal(ApplyPromoRequest{
		PromoCode: strings.TrimSpace(promoCode),
	})
	if err != nil {
		return fmt.Errorf("marshaling apply promo request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/yooz-mgm/apply-code?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("applying yooz promo code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}

	var res ApplyPromoResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("decoding apply promo response: %w", err)
	}

	if !res.Success {
		return fmt.Errorf("%w: %s", ErrRequestFailed, res.Message)
	}
	return nil
}

// ========================================================================
// 17. E-Vouchers (كروت الألعاب والشحن الرقمي)
// ========================================================================

// GetEVoucherPackages browses available digital gaming and app store gift cards (PUBG, iTunes, PlayStation).
func (c *Client) GetEVoucherPackages(ctx context.Context) (*EVoucherPackagesResponse, error) {
	path := fmt.Sprintf("/api/v2/e-voucher/packages?recharge-type=1&lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting evoucher packages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res EVoucherPackagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding evoucher packages response: %w", err)
	}
	return &res, nil
}

// ========================================================================
// 18. Active Subscriptions & Services (الاشتراكات والخدمات المفعلة)
// ========================================================================

// GetMySubscriptions fetches all active services and subscriptions on the line directly from Asiacell servers (/api/v1/profile/subscriptions).
func (c *Client) GetMySubscriptions(ctx context.Context) ([]MySubscriptionItem, error) {
	path := fmt.Sprintf("/api/v1/profile/subscriptions?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting subscriptions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res MySubscriptionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding subscriptions response: %w", err)
	}
	return res.Data, nil
}

// ========================================================================
// 19. USSD Interactive Cloud Menu (قائمة خدمات USSD السحابية التفاعلية)
// ========================================================================

// GetUSSDMenu fetches USSD menu questions and interactive options from Asiacell cloud (/api/v1/ussd).
func (c *Client) GetUSSDMenu(ctx context.Context, parentID int) ([]USSDAnswerItem, error) {
	path := fmt.Sprintf("/api/v1/ussd?parent_id=%d&lang=%s", parentID, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting ussd menu: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res USSDQuestionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding ussd menu response: %w", err)
	}
	return res.Data, nil
}

// SubmitUSSDAction submits a choice or action for a USSD interactive menu item (/api/v1/ussd).
func (c *Client) SubmitUSSDAction(ctx context.Context, params map[string]string) (*USSDResultData, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshaling ussd action request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/ussd?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("submitting ussd action: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res UssdQAResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding ussd action response: %w", err)
	}
	return res.Data, nil
}

// ========================================================================
// 20. Active Bundle Details & Action Controls (/api/v1/profile/bundle/{key})
// ========================================================================

// GetBundleDetail retrieves in-depth quota breakdown, validity and action buttons for an active bundle key.
func (c *Client) GetBundleDetail(ctx context.Context, bundleKey string) (*AccountBundleDetailData, error) {
	path := fmt.Sprintf("/api/v1/profile/bundle/%s?lang=%s", url.PathEscape(bundleKey), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting bundle detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ProfileBundleResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding bundle detail response: %w", err)
	}
	return res.Data, nil
}

// ========================================================================
// 21. Server-Side Session Logout (/api/v1/logout)
// ========================================================================

// Logout invalidates the active access token and session on Asiacell cloud servers.
func (c *Client) Logout(ctx context.Context) error {
	c.mu.RLock()
	token := c.accessToken
	c.mu.RUnlock()

	payload, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		return fmt.Errorf("marshaling logout request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/logout?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("requesting logout: %w", err)
	}
	defer resp.Body.Close()

	c.mu.Lock()
	c.accessToken = ""
	c.refreshToken = ""
	c.mu.Unlock()

	return nil
}

// ========================================================================
// 22. Multi-Account & Linked Lines Management (/api/v1/map-account)
// ========================================================================

// GetLinkedAccounts fetches all phone numbers mapped/linked to the primary account line.
func (c *Client) GetLinkedAccounts(ctx context.Context) ([]string, error) {
	path := fmt.Sprintf("/api/v1/map-account?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting linked accounts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res AccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding linked accounts: %w", err)
	}
	return res.Data, nil
}

// AddLinkedAccount requests linking an additional phone number to the account, sending an OTP verification SMS.
func (c *Client) AddLinkedAccount(ctx context.Context, phone string) error {
	payload, err := json.Marshal(map[string]string{"number": phone})
	if err != nil {
		return fmt.Errorf("marshaling map account request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/map-account?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting map account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// ConfirmLinkedAccount confirms mapping the secondary number using the received SMS PIN.
func (c *Client) ConfirmLinkedAccount(ctx context.Context, phone, pin string) error {
	payload, err := json.Marshal(map[string]string{"number": phone, "pin": pin})
	if err != nil {
		return fmt.Errorf("marshaling confirm map request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/map-account/confirm?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting confirm map: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// SwitchActiveAccount switches active line context to a linked secondary number.
func (c *Client) SwitchActiveAccount(ctx context.Context, otherPhone string) error {
	payload, err := json.Marshal(map[string]string{"otherNumber": otherPhone})
	if err != nil {
		return fmt.Errorf("marshaling switch account request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/account-action/switch?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting switch account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// RemoveLinkedAccount removes a secondary linked number from the account mapping.
func (c *Client) RemoveLinkedAccount(ctx context.Context, otherPhone string) error {
	payload, err := json.Marshal(map[string]string{"otherNumber": otherPhone})
	if err != nil {
		return fmt.Errorf("marshaling remove account request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/account-action/remove?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting remove account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// ========================================================================
// 23. Data Caps & Sharing Limits (/api/v1/addon/datacap & share)
// ========================================================================

// SetDataCapLimit sets a daily data consumption threshold limit in MB (/api/v1/addon/datacap/set-limit).
func (c *Client) SetDataCapLimit(ctx context.Context, limitMB float64) error {
	payload, err := json.Marshal(map[string]float64{"value": limitMB})
	if err != nil {
		return fmt.Errorf("marshaling datacap limit: %w", err)
	}

	path := fmt.Sprintf("/api/v1/addon/datacap/set-limit?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting datacap limit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// SetBundleShareLimit sets data sharing quota threshold in MB for a specific shared line (/api/v1/addon/share/set-limit).
func (c *Client) SetBundleShareLimit(ctx context.Context, msisdn string, limitMB float64) error {
	type limitItem struct {
		Value float64 `json:"value"`
	}
	type reqBody struct {
		MSISDN     string      `json:"msisdn"`
		LimitItems []limitItem `json:"limitItems"`
	}

	payload, err := json.Marshal(reqBody{
		MSISDN:     msisdn,
		LimitItems: []limitItem{{Value: limitMB}},
	})
	if err != nil {
		return fmt.Errorf("marshaling share bundle limit: %w", err)
	}

	path := fmt.Sprintf("/api/v1/addon/share/set-limit?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting share bundle limit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// ========================================================================
// 24. Shake & Win (/api/v1/shake-and-win)
// ========================================================================

func (c *Client) GetShakeAndWinStatus(ctx context.Context) (*ShakeAndWinData, error) {
	path := fmt.Sprintf("/api/v1/shake-and-win?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting shake and win status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ShakeAndWinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding shake and win response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) PlayShakeAndWin(ctx context.Context, transactionID string) (*ShakeAndWinData, error) {
	payload, err := json.Marshal(map[string]string{"transactionId": transactionID})
	if err != nil {
		return nil, fmt.Errorf("marshaling play shake and win: %w", err)
	}

	path := fmt.Sprintf("/api/v1/shake-and-win?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("submitting shake and win: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ShakeAndWinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding shake and win response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetTopupShakeAndWin(ctx context.Context) (*ShakeAndWinData, error) {
	path := fmt.Sprintf("/api/v1/top-up/shake-and-win?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting topup shake and win: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ShakeAndWinResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding topup shake and win: %w", err)
	}
	return res.Data, nil
}

// ========================================================================
// 25. Postpaid Bill Payment (/api/v1/top-up/bill-amount & pay-bill)
// ========================================================================

func (c *Client) GetBillAmount(ctx context.Context) (*BillInfoResponse, error) {
	path := fmt.Sprintf("/api/v1/top-up/bill-amount?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting bill amount: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res BillInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding bill amount: %w", err)
	}
	return &res, nil
}

func (c *Client) PayBill(ctx context.Context, phone string, amount float64) error {
	payload, err := json.Marshal(map[string]any{
		"msisdn":       phone,
		"rechargeType": 1,
		"amount":       amount,
	})
	if err != nil {
		return fmt.Errorf("marshaling pay bill request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/top-up/pay-bill?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting pay bill: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

// ========================================================================
// 26. Ticket Tracking & Form Details (/api/v1/resolution-center)
// ========================================================================

func (c *Client) GetTicketDetail(ctx context.Context, ticketNumber string) (*TicketDetailItem, error) {
	path := fmt.Sprintf("/api/v1/resolution-center/%s?lang=%s", ticketNumber, c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting ticket detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res TicketDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding ticket detail: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetTicketForm(ctx context.Context, category string) ([]TicketFormField, error) {
	path := fmt.Sprintf("/api/v1/resolution-center/ticket-form?category=%s&lang=%s", url.QueryEscape(category), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting ticket form: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res TicketFormResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding ticket form: %w", err)
	}
	return res.Data, nil
}

// ========================================================================
// 27. Data Lines & Routers (/api/v1/data-line & /api/v1/multi-line)
// ========================================================================

func (c *Client) GetDataLineInfo(ctx context.Context) (*DataLineInfo, error) {
	path := fmt.Sprintf("/api/v1/data-line?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting data line info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res DataLineResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding data line info: %w", err)
	}
	return res.Data, nil
}

func (c *Client) PairDataLine(ctx context.Context, msisdn, iccid string) error {
	payload, err := json.Marshal(map[string]any{
		"msisdn":      msisdn,
		"iccid":       iccid,
		"replacement": false,
	})
	if err != nil {
		return fmt.Errorf("marshaling pair data line request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/data-line?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting pair data line: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetMultiLineConnections(ctx context.Context) ([]MultiLineConnection, error) {
	path := fmt.Sprintf("/api/v1/multi-line?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting multi-line connections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res MultiLineResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding multi-line response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetMultiLineHome(ctx context.Context) (*HomeData, error) {
	path := fmt.Sprintf("/api/v1/multi-line/home?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting multi-line home: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res HomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding multi-line home: %w", err)
	}
	return &res.Data, nil
}

// ========================================================================
// 28. Global Search Engine (/api/v1/search & suggestions)
// ========================================================================

func (c *Client) Search(ctx context.Context, query string) ([]SearchItem, error) {
	path := fmt.Sprintf("/api/v1/search?q=%s&lang=%s", url.QueryEscape(query), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("executing search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res SearchResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetSearchSuggestions(ctx context.Context, query string) ([]string, error) {
	path := fmt.Sprintf("/api/v1/search/suggestions?q=%s&lang=%s", url.QueryEscape(query), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting search suggestions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res SearchSuggestionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding search suggestions: %w", err)
	}
	return res.Data, nil
}

// ========================================================================
// 29. International Services & Tariffs (/api/v1/international-services)
// ========================================================================

func (c *Client) GetInternationalTariffs(ctx context.Context) ([]InternationalCountryTariff, error) {
	path := fmt.Sprintf("/api/v1/international-services/tariff?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting international tariffs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res InternationalTariffResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding international tariffs: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetInternationalServices(ctx context.Context) ([]PromotionItem, error) {
	path := fmt.Sprintf("/api/v1/international-services?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting international services: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res PromotionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding international services: %w", err)
	}

	var items []PromotionItem
	for _, body := range res.Data.Bodies {
		items = append(items, body.Items...)
	}
	return items, nil
}

// ========================================================================
// 30. Loyalty Rewards & Wafaa (/api/v1/reward & /api/v1/eo)
// ========================================================================

func (c *Client) GetLoyaltyRewards(ctx context.Context) ([]LoyaltyRewardItem, error) {
	path := fmt.Sprintf("/api/v1/reward?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting loyalty rewards: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res LoyaltyRewardsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding loyalty rewards: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetLoyaltyRewardDetail(ctx context.Context) (*LoyaltyRewardItem, error) {
	path := fmt.Sprintf("/api/v1/reward/detail?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting reward detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res struct {
		Success bool               `json:"success"`
		Data    *LoyaltyRewardItem `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding reward detail: %w", err)
	}
	return res.Data, nil
}

func (c *Client) RedeemLoyaltyReward(ctx context.Context) (*RedeemRewardResponse, error) {
	path := fmt.Sprintf("/api/v1/eo?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("redeeming loyalty reward: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res RedeemRewardResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding redeem reward response: %w", err)
	}
	return &res, nil
}

func (c *Client) CheckLoyaltyRewardStatus(ctx context.Context, pid string) (bool, error) {
	path := fmt.Sprintf("/api/v1/eo/check-status?pid=%s&lang=%s", url.QueryEscape(pid), c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return false, fmt.Errorf("checking reward status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false, ErrUnauthorized
	}

	var res struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, fmt.Errorf("decoding check reward status: %w", err)
	}
	return res.Success, nil
}

// ========================================================================
// 31. Profile Update & Limits Inspection
// ========================================================================

func (c *Client) UpdateProfileInfo(ctx context.Context, name, email string) error {
	payload, err := json.Marshal(map[string]string{
		"name":  name,
		"email": email,
	})
	if err != nil {
		return fmt.Errorf("marshaling update profile request: %w", err)
	}

	path := fmt.Sprintf("/api/v3/profile/update?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("submitting update profile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrUnauthorized
	}
	return nil
}

func (c *Client) GetDataCapLimit(ctx context.Context) ([]LineLimitInfo, error) {
	path := fmt.Sprintf("/api/v1/addon/datacap/limit?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting datacap limit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res LineLimitsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding datacap limit response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetBundleShareLimit(ctx context.Context) ([]LineLimitInfo, error) {
	path := fmt.Sprintf("/api/v1/addon/share/limit?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting bundle share limit: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res LineLimitsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding share limit response: %w", err)
	}
	return res.Data, nil
}

func (c *Client) GetManageLines(ctx context.Context) ([]SharedLineItem, error) {
	path := fmt.Sprintf("/api/v1/addon/share?lang=%s", c.language)
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting shared lines: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrUnauthorized
	}

	var res ManageLineResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding manage lines: %w", err)
	}
	return res.Data, nil
}

