package asiacell

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

type FlexString string

func (fs *FlexString) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*fs = ""
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*fs = FlexString(s)
		return nil
	}
	*fs = FlexString(strings.Trim(string(b), "\""))
	return nil
}

func (fs FlexString) String() string {
	return string(fs)
}

func (fs FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(fs))
}

type RechargeType int

const (
	RechargeTypeNormal   RechargeType = 1
	RechargeTypeInternet RechargeType = 2
)

type LoginRequest struct {
	Username    string `json:"username"`
	CaptchaCode string `json:"captchaCode"`
}

type LoginResponse struct {
	Success        bool         `json:"success"`
	NextURL        string       `json:"nextUrl"`
	Username       string       `json:"username"`
	Message        string       `json:"message"`
	RequireCaptcha bool         `json:"requireCaptcha"`
	Captcha        *CaptchaData `json:"captcha,omitempty"`
}

type SMSValidationRequest struct {
	PID      string `json:"PID"`
	Passcode string `json:"passcode"`
	Token    string `json:"token,omitempty"`
}

type SMSValidationResponse struct {
	Success        bool       `json:"success"`
	AccessToken    string     `json:"access_token"`
	RefreshToken   string     `json:"refresh_token"`
	HandshakeToken string     `json:"handshake_token"`
	UserID         FlexString `json:"userId"`
	Username       string     `json:"username"`
	UserType       FlexString `json:"userType"`
	Language       string     `json:"language"`
	Message        string     `json:"message"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type ActionDetail struct {
	Title  FlexString `json:"title"`
	Action string     `json:"action"`
}

type ProfileItem struct {
	Name            FlexString    `json:"name"`
	PhoneNumber     FlexString    `json:"phoneNumber"`
	Photo           string        `json:"photo"`
	Value           FlexString    `json:"value"`
	Validity        string        `json:"validity"`
	ActionButton    *ActionDetail `json:"actionButton"`
	Title           FlexString    `json:"title"`
	Action          *ActionDetail `json:"action"`
	RemainingVolume float64       `json:"remainingVolume"`
	TotalVolume     float64       `json:"totalVolume"`
	Unit            string        `json:"unit"`
	ExpireLabel     string        `json:"expireLabel"`
	ExpireDate      string        `json:"expireDate"`
	SubscribeDate   string        `json:"subscribeDate"`
	BundleKey       string        `json:"bundleKey"`
	FreeUnitRefName string        `json:"freeUnitRefName"`
	ActiveBundle    bool          `json:"activeBundle"`
}

type ProfileBody struct {
	Title FlexString    `json:"title"`
	Type  string        `json:"type"`
	Items []ProfileItem `json:"items"`
}

type ProfileData struct {
	Bodies []ProfileBody `json:"bodies"`
}

type ProfileResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    ProfileData `json:"data"`
}

type ActiveBundleInfo struct {
	Title           string  `json:"title"`
	RemainingVolume float64 `json:"remainingVolume"`
	TotalVolume     float64 `json:"totalVolume"`
	Unit            string  `json:"unit"`
	ExpireLabel     string  `json:"expireLabel"`
	ExpireDate      string  `json:"expireDate"`
	SubscribeDate   string  `json:"subscribeDate"`
	BundleKey       string  `json:"bundleKey"`
	FreeUnitRefName string  `json:"freeUnitRefName"`
	ActiveBundle    bool    `json:"activeBundle"`
}

type AccountOverview struct {
	Name           string             `json:"name"`
	PhoneNumber    string             `json:"phoneNumber"`
	Balance        string             `json:"balance"`
	Validity       string             `json:"validity"`
	RemainingData  string             `json:"remainingData"`
	RemainingCalls string             `json:"remainingCalls"`
	RemainingSMS   string             `json:"remainingSMS"`
	ActiveBundles  []ActiveBundleInfo `json:"activeBundles"`
}

type CaptchaData struct {
	CaptchaCode  string `json:"captchaCode"`
	CaptchaImage string `json:"captchaImage"`
	CaptchaURL   string `json:"captchaUrl"`
	OriginSource string `json:"originSource"`
	ResourceURL  string `json:"resourceUrl"`
}

func (c *CaptchaData) ImageURL() string {
	if c == nil {
		return ""
	}
	if c.CaptchaURL != "" {
		return c.CaptchaURL
	}
	if c.OriginSource != "" && c.ResourceURL != "" {
		return c.OriginSource + c.ResourceURL
	}
	if c.CaptchaImage != "" {
		return c.CaptchaImage
	}
	return ""
}

type CaptchaResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    CaptchaData `json:"data"`
}

type ProfileDetails struct {
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"thirdName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
	PhotoURL    string `json:"photo"`
}

type ProfileDetailsResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    ProfileDetails `json:"data"`
}

type RechargeRequest struct {
	MSISDN       string       `json:"msisdn"`
	RechargeType RechargeType `json:"rechargeType"`
	Voucher      string       `json:"voucher"`
}

type RechargeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CreditTransferStartRequest struct {
	Amount         float64 `json:"amount"`
	ReceiverMSISDN string  `json:"receiverMsisdn"`
}

type CreditTransferStartResponse struct {
	Success bool       `json:"success"`
	PID     FlexString `json:"PID"`
	Message string     `json:"message"`
}

type CreditTransferDoRequest struct {
	PID      string `json:"PID"`
	Passcode string `json:"passcode"`
}

type TransferConfirmation struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Title      FlexString `json:"title"`
	NextAction string     `json:"nextAction"`
}

type TransactionRecord struct {
	Type           FlexString `json:"type"`
	MSISDN         FlexString `json:"msisdn"`
	ReceiverMSISDN FlexString `json:"receiverMsisdn"`
	CreatedAt      FlexString `json:"createdAt"`
	Amount         FlexString `json:"amount"`
	Voucher        FlexString `json:"voucher"`
}

type RechargeHistoryResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []TransactionRecord `json:"data"`
}

type TransferHistoryResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []TransactionRecord `json:"data"`
}

type BundleRecord struct {
	BundleName   FlexString `json:"bundleName"`
	ActivationAt FlexString `json:"activationAt"`
	Amount       FlexString `json:"amount"`
}

type SubscriptionHistoryResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []BundleRecord `json:"data"`
}

type SpinWheelRemainingTime struct {
	Hours      int   `json:"hours"`
	Minutes    int   `json:"minutes"`
	Second     int   `json:"second"`
	Timestamp  int64 `json:"timestamp"`
	IsPlayable bool  `json:"isPlayable"`
	LuckyID    int   `json:"luckyId"`
	DayIndex   int   `json:"dayIndex"`
}

type SpinWheelSlide struct {
	GrandPrize bool   `json:"grandPrize"`
	Color      string `json:"color"`
	Background string `json:"background"`
	Title      string `json:"title"`
}

type SpinWheelData struct {
	IsPlayable    bool                   `json:"isPlayable"`
	Playable      bool                   `json:"playable"`
	LuckyID       int                    `json:"luckyId"`
	DayIndex      int                    `json:"dayIndex"`
	Title         string                 `json:"title"`
	SubTitle      string                 `json:"subTitle"`
	Desc          string                 `json:"desc"`
	RemainingTime SpinWheelRemainingTime `json:"remainingTime"`
	Slides        []SpinWheelSlide       `json:"slides"`
}

type SpinWheelStatusResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    SpinWheelData `json:"data"`
}

type SpinWheelPlayAnalytics struct {
	Event  string            `json:"event"`
	Params map[string]string `json:"params"`
}

type SpinWheelPlayResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	Data         string                 `json:"data"`
	AnalyticData SpinWheelPlayAnalytics `json:"analyticData"`
	Title        string                 `json:"title"`
	NextAction   string                 `json:"nextAction"`
}

type AddonPackage struct {
	Index          int      `json:"index"`
	ID             int      `json:"id"`
	Title          string   `json:"title"`
	Volume         string   `json:"volume"`
	Validity       string   `json:"validity"`
	Price          string   `json:"price"`
	RelatedProduct string   `json:"relatedProduct"`
	Data           string   `json:"data"`
	Renewable      bool     `json:"renewable"`
	FreeSocials    []string `json:"freeSocials"`
	USSDCode       string   `json:"ussdCode"`
	SMSCode        string   `json:"smsCode"`
}

type AddonRawItem struct {
	ID             int      `json:"id"`
	Title          string   `json:"title"`
	Volume         string   `json:"volume"`
	Validity       string   `json:"validity"`
	Price          string   `json:"price"`
	RelatedProduct string   `json:"relatedProduct"`
	Data           string   `json:"data"`
	Renewable      bool     `json:"renewable"`
	FreeSocials    []string `json:"freeSocials"`
	Type           string   `json:"type"`
}

type AddonRawGroup struct {
	GroupID int            `json:"groupId"`
	Type    string         `json:"type"`
	Title   string         `json:"title"`
	Items   []AddonRawItem `json:"items"`
}

type AddonResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Bodies []AddonRawGroup `json:"bodies"`
	} `json:"data"`
}

type AddonCategory struct {
	GroupID int            `json:"groupId"`
	Type    string         `json:"type"`
	Title   string         `json:"title"`
	Items   []AddonPackage `json:"items"`
}

type AddonDetailData struct {
	Headers []struct {
		Title           string `json:"title"`
		Validity        string `json:"validity"`
		Data            string `json:"data"`
		Minute          string `json:"minute"`
		SMS             string `json:"sms"`
		FeatureImage    string `json:"featureImage"`
		BackgroundImage string `json:"backgroundImage"`
		BackgroundColor string `json:"backgroundColor"`
	} `json:"headers"`
	Bodies []struct {
		GroupID int    `json:"groupId"`
		Type    string `json:"type"`
		Items   []struct {
			Price       string `json:"price"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Action      struct {
				Title  string `json:"title"`
				Action string `json:"action"`
			} `json:"action"`
		} `json:"items"`
	} `json:"bodies"`
}

type AddonDetailResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    AddonDetailData `json:"data"`
}

type SubscriptionResult struct {
	Success     bool   `json:"success"`
	PackageID   int    `json:"packageId"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	USSDCommand string `json:"ussdCommand"`
	SMSCommand  string `json:"smsCommand"`
	Message     string `json:"message"`
}

type ServiceActionItem struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	URL    string `json:"url,omitempty"`
}

type ServiceCancellationInfo struct {
	ID          string              `json:"id"`
	Title       string              `json:"title"`
	Code        string              `json:"code"`
	Method      string              `json:"method"`
	Description string              `json:"description"`
	Actions     []ServiceActionItem `json:"actions"`
}

type InternetControlInfo struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	USSDCode     string   `json:"ussdCode"`
	CustomerCare string   `json:"customerCare"`
	WhatsAppCare string   `json:"whatsAppCare"`
	Instructions []string `json:"instructions"`
}

type ServicesManagementOverview struct {
	ActiveBundles      []ActiveBundleInfo        `json:"activeBundles"`
	CancellationGuides []ServiceCancellationInfo `json:"cancellationGuides"`
	InternetControl    InternetControlInfo       `json:"internetControl"`
}

type ServiceActionResult struct {
	Success     bool   `json:"success"`
	ServiceID   string `json:"serviceId"`
	Title       string `json:"title"`
	ActionType  string `json:"actionType"`
	Value       string `json:"value"`
	URL         string `json:"url,omitempty"`
	Instruction string `json:"instruction"`
}

type NotificationItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}

type NotificationsResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    []NotificationItem `json:"data"`
}

type ShukranAction struct {
	Title string `json:"title"`
	Action string `json:"action"`
}

type ShukranItem struct {
	Label string `json:"label"`
	Type string `json:"type"`
	Action ShukranAction `json:"action"`
}

type ShukranHeader struct {
	ScreenTitle string `json:"screenTitle"`
	Title string `json:"title"`
	Data string `json:"data"`
	Credit string `json:"credit"`
	FeatureImage string `json:"featureImage"`
}

type ShukranBody struct {
	Type string `json:"type"`
	Items []ShukranItem `json:"items"`
}

type ShukranData struct {
	Headers []ShukranHeader `json:"headers"`
	Bodies []ShukranBody `json:"bodies"`
}

type ShukranResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data ShukranData `json:"data"`
}

type ShukranActionResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type RoamingTag struct {
	Tag string `json:"tag"`
	Title string `json:"title"`
	Ordering int `json:"ordering"`
}

type RoamingItem struct {
	Title string `json:"title"`
	Tag string `json:"tag"`
}

type RoamingData struct {
	ScreenTitle string `json:"screenTitle"`
	Tags []RoamingTag `json:"tags"`
	Items []RoamingItem `json:"items"`
}

type RoamingResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data RoamingData `json:"data"`
}

type PromotionButton struct {
	Title string `json:"title"`
	Action string `json:"action"`
}

type PromotionItem struct {
	ID int `json:"id"`
	TitleDesc string `json:"titleDesc"`
	BackgroundImage string `json:"backgroundImage"`
	ActionButton PromotionButton `json:"actionButton"`
}

type PromotionsData struct {
	Bodies []struct {
		Items []PromotionItem `json:"items"`
	} `json:"bodies"`
}

type PromotionsResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data PromotionsData `json:"data"`
}

type SpinWheelCheckResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	NextAction string `json:"nextAction"`
}

type ProfileViewData struct {
	FirstName FlexString `json:"firstName"`
	LastName FlexString `json:"lastName"`
	ThirdName FlexString `json:"thirdName"`
	Email FlexString `json:"email"`
	Phone FlexString `json:"phone"`
	Photo string `json:"photo"`
}

type ProfileViewResponse struct {
	Success bool `json:"success"`
	Data ProfileViewData `json:"data"`
}

type UserInterest struct {
	Key string `json:"key"`
	Title string `json:"title"`
	Checked bool `json:"checked"`
}

type UserCover struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Selected bool `json:"selected"`
}

type UserAvatar struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Selected bool `json:"selected"`
}

type ProfileV2Data struct {
	Interests []UserInterest `json:"interests"`
	Covers []UserCover `json:"covers"`
	Avatars []UserAvatar `json:"avatars"`
}

type ProfileV2Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data ProfileV2Data `json:"data"`
}

type SessionData struct {
	AccessToken    string     `json:"access_token"`
	RefreshToken   string     `json:"refresh_token"`
	HandshakeToken string     `json:"handshake_token"`
	UserID         FlexString `json:"userId"`
	Username       string     `json:"username"`
	DeviceID       string     `json:"deviceId"`
	Language       string     `json:"language"`
	LastRefresh    time.Time  `json:"last_refresh,omitempty"`
}

type HomeHeader struct {
	BackgroundColor   string             `json:"backgroundColor"`
	BackgroundImage   string             `json:"backgroundImage"`
	InvertedTextColor bool               `json:"invertedTextColor"`
	ActionButton      PromotionButton `json:"actionButton"`
}

type HomeItem struct {
	ID           int                `json:"id"`
	GroupID      int                `json:"groupId"`
	Title        string             `json:"title"`
	Icon         string             `json:"icon"`
	Unlimited    bool               `json:"unlimited"`
	Volume       string             `json:"volume"`
	Validity     string             `json:"validity"`
	Price        string             `json:"price"`
	Data         string             `json:"data"`
	Action       string             `json:"action"`
	Tag          string             `json:"tag"`
	Renewable    bool               `json:"renewable"`
	FreeSocials  []string           `json:"freeSocials"`
	ActionButton PromotionButton `json:"actionButton"`
}

type HomeBody struct {
	GroupID int        `json:"groupId"`
	Type    string     `json:"type"`
	Action  string     `json:"action"`
	Title   string     `json:"title"`
	Items   []HomeItem `json:"items"`
}

type HomeData struct {
	Headers []HomeHeader `json:"headers"`
	Bodies  []HomeBody   `json:"bodies"`
}

type HomeResponse struct {
	Data HomeData `json:"data"`
}

type CityPartner struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	TotalCategory int    `json:"totalCategory"`
	HomeDelivery  bool   `json:"homeDelivery"`
}

type CitiesResponse struct {
	Success bool          `json:"success"`
	Data    []CityPartner `json:"data"`
}

type DigitalServiceItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type DigitalServicesInfo struct {
	Title    string               `json:"title"`
	Services []DigitalServiceItem `json:"services"`
}

type ShopWorkingHour struct {
	WorkingDay  int    `json:"workingDay"`
	StartHour   int    `json:"startHour"`
	FinishHour  int    `json:"finishHour"`
	DayTitle    string `json:"dayTitle"`
	WorkingHour string `json:"workingHour"`
}

type ShopInfo struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Address      string            `json:"address"`
	Phone        string            `json:"phone"`
	Lat          float64           `json:"lat"`
	Lng          float64           `json:"lng"`
	Status       string            `json:"status"`
	CityName     string            `json:"cityName"`
	CityID       int               `json:"cityId"`
	StartHour    string            `json:"startHour"`
	FinishHour   string            `json:"finishHour"`
	WorkingDays  string            `json:"workingDays"`
	WorkingHours []ShopWorkingHour `json:"workingHours"`
}

type ShopsResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    []ShopInfo `json:"data"`
}

type AddonSummaryBalance struct {
	ID           int             `json:"id"`
	Method       string          `json:"method"`
	Title        string          `json:"title"`
	Desc         string          `json:"desc"`
	ActionButton PromotionButton `json:"actionButton"`
}

type AddonSummaryOnlinePayment struct {
	BalanceLabel string              `json:"balanceLabel"`
	Balance      AddonSummaryBalance `json:"balance"`
}

type AddonSummaryOption struct {
	Key           string                    `json:"key"`
	Title         string                    `json:"title"`
	Price         string                    `json:"price"`
	Validity      string                    `json:"validity"`
	Volumed       string                    `json:"volumed"`
	Total         float64                   `json:"total"`
	ActionButton  PromotionButton           `json:"actionButton"`
	OnlinePayment AddonSummaryOnlinePayment `json:"onlinePayment"`
}

type AddonSummaryData struct {
	ID      int                  `json:"id"`
	Options []AddonSummaryOption `json:"options"`
}

type AddonSummaryResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    AddonSummaryData `json:"data"`
}

type AddonSubscribeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}



