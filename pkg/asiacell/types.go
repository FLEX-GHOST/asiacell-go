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

// CDRRecord represents an in-app Call Detail Record (CDR) ledger entry for balance transfers.
type CDRRecord struct {
	Amount      FlexString `json:"amount"`      // e.g. "1000 IQD" (or "-1000 IQD" for outgoing)
	Unit        FlexString `json:"unit"`        // "TRANSFERS"
	Title       FlexString `json:"title"`       // "تحويل الرصيد"
	SubTitle    FlexString `json:"subTitle"`    // Sender/Receiver phone number (e.g. "7744298878")
	Description FlexString `json:"description"` // Timestamp e.g. "١١/٠٩/٢٠٢٦ ٠٦:٢٣:٤٣"
}

// CDRDetailData contains the paginated list of CDR records.
type CDRDetailData struct {
	Total int         `json:"total"`
	Data  []CDRRecord `json:"data"`
}

// CDRDetailResponse represents the response envelope from /api/v1/cdr/detail.
type CDRDetailResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Success bool           `json:"success"`
	Data    *CDRDetailData `json:"data"`
}

// CDRConfirmRequest represents the body payload for /api/v1/cdr/confirm.
type CDRConfirmRequest struct {
	Code string `json:"code"`
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

type AddOnBundleTag struct {
	Tag   string `json:"tag"`
	Title string `json:"title"`
}

type AddOnTagBundles struct {
	ScreenTitle string           `json:"screenTitle"`
	Items       []AddonRawItem   `json:"items"`
	Tags        []AddOnBundleTag `json:"tags"`
}

type AddOnTagBundlesResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    AddOnTagBundles `json:"data"`
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



// --- Vanity VIP Numbers ---

type VanityClass struct {
	ID    FlexString `json:"id"`
	Title FlexString `json:"title"`
	Price FlexString `json:"price"`
	Icon  string     `json:"icon"`
}

type VanityClassesResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []VanityClass `json:"data"`
}

type VanityNumberItem struct {
	MSISDN    FlexString `json:"msisdn"`
	ClassID   FlexString `json:"classId"`
	ClassName FlexString `json:"className"`
	Price     FlexString `json:"price"`
	Currency  FlexString `json:"currency"`
	Status    FlexString `json:"status"`
}

type VanitySearchData struct {
	Total int                `json:"total"`
	List  []VanityNumberItem `json:"list"`
}

type VanitySearchResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *VanitySearchData `json:"data"`
}

type VanityDetailResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *VanityNumberItem `json:"data"`
}

type ReserveVanityRequest struct {
	MSISDN  string `json:"msisdn"`
	ClassID string `json:"classId"`
}

type ReserveVanityResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	PID     FlexString `json:"pid"`
}

// --- Gifting Addons ---

type SendGiftRequest struct {
	AddonID        int    `json:"addOnId"`
	ReceiverMSISDN string `json:"receiverMsisdn"`
}

type SendGiftResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Resolution Center / Support Tickets ---

type TicketCategory struct {
	ID    FlexString `json:"id"`
	Title FlexString `json:"title"`
	Icon  string     `json:"icon"`
}

type TicketCategoriesResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []TicketCategory `json:"data"`
}

type TicketItem struct {
	TicketNumber string     `json:"ticketNumber"`
	Category     FlexString `json:"category"`
	Status       FlexString `json:"status"`
	Subject      FlexString `json:"subject"`
	Description  FlexString `json:"description"`
	CreatedAt    FlexString `json:"createdAt"`
	UpdatedAt    FlexString `json:"updatedAt"`
}

type TicketsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []TicketItem `json:"data"`
}

type SubmitTicketRequest struct {
	CategoryID  string `json:"category"`
	Description string `json:"description"`
}

type SubmitTicketResponse struct {
	Success      bool       `json:"success"`
	Message      string     `json:"message"`
	TicketNumber FlexString `json:"ticketNumber"`
}

// --- Compensation System ---

type CompensationItem struct {
	ID          FlexString `json:"id"`
	Title       FlexString `json:"title"`
	Description FlexString `json:"description"`
	Benefit     FlexString `json:"benefit"`
	Eligible    bool       `json:"eligible"`
	Claimed     bool       `json:"claimed"`
}

type CompensationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    []CompensationItem `json:"data"`
}

// --- Yooz MGM Referral System ---

type YoozMGMData struct {
	ReferralCode string `json:"referralCode"`
	ShareLink    string `json:"shareLink"`
	TotalInvites int    `json:"totalInvites"`
	TotalEarned  string `json:"totalEarned"`
}

type YoozMGMResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *YoozMGMData `json:"data"`
}

type ApplyPromoRequest struct {
	PromoCode string `json:"promoCode"`
}

type ApplyPromoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- E-Voucher Packages (Digital Gaming / App Cards) ---

type EVoucherPackageItem struct {
	ID       int        `json:"id"`
	Title    FlexString `json:"title"`
	Category FlexString `json:"category"`
	Price    FlexString `json:"price"`
	ImageURL string     `json:"imageUrl"`
}

type EVoucherPackagesResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []EVoucherPackageItem `json:"data"`
}

// --- Active Subscriptions & Services (الاشتراكات والخدمات الفعالة) ---

type ActionButton struct {
	Title           FlexString `json:"title"`
	Icon            FlexString `json:"icon"`
	Action          FlexString `json:"action"`
	BackgroundColor FlexString `json:"backgroundColor"`
	Inverted        bool       `json:"inverted"`
	Style           FlexString `json:"style"`
	ViewID          int        `json:"viewId"`
	Disabled        bool       `json:"disabled"`
}

type MySubscriptionItem struct {
	Title        FlexString    `json:"title"`
	Validity     FlexString    `json:"validity"`
	ActionButton *ActionButton `json:"actionButton"`
}

type MySubscriptionsResponse struct {
	Success    bool                 `json:"success"`
	Message    string               `json:"message"`
	Title      string               `json:"title"`
	NextAction string               `json:"nextAction"`
	Data       []MySubscriptionItem `json:"data"`
}

// --- USSD Interactive Cloud Menu (قائمة أوامر USSD التفاعلية عبر السحابة) ---

type USSDAnswerItem struct {
	ID         int        `json:"id"`
	Content    FlexString `json:"content"`
	USSDID     FlexString `json:"ussdId"`
	SubContent FlexString `json:"subContent"`
}

type USSDQuestionsResponse struct {
	Success    bool             `json:"success"`
	Message    string           `json:"message"`
	Title      string           `json:"title"`
	NextAction string           `json:"nextAction"`
	Data       []USSDAnswerItem `json:"data"`
}

type USSDResultData struct {
	Title    FlexString    `json:"title"`
	Msg      FlexString    `json:"msg"`
	Image    string        `json:"image"`
	Positive *ActionButton `json:"positive"`
	Negative *ActionButton `json:"negative"`
}

type UssdQAResultResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Title      string          `json:"title"`
	NextAction string          `json:"nextAction"`
	Data       *USSDResultData `json:"data"`
}

// --- Active Bundle Details & Action Buttons (/api/v1/profile/bundle/{key}) ---

type BundleDetailItem struct {
	Title         FlexString     `json:"title"`
	Label         FlexString     `json:"label"`
	Volume        float64        `json:"volume"`
	Unit          FlexString     `json:"unit"`
	Unlimited     bool           `json:"unlimited"`
	TotalVolume   float64        `json:"totalVolume"`
	Validity      FlexString     `json:"validity"`
	BarColor      FlexString     `json:"barColor"`
	DetailLabel   FlexString     `json:"detailLabel"`
	ActionButtons []ActionButton `json:"actionButtons"`
}

type AccountBundleDetailData struct {
	Bodies []BundleDetailItem `json:"bodies"`
}

type ProfileBundleResponse struct {
	Success    bool                     `json:"success"`
	Message    string                   `json:"message"`
	Title      string                   `json:"title"`
	NextAction string                   `json:"nextAction"`
	Data       *AccountBundleDetailData `json:"data"`
}

// --- Multi-Account / Linked Lines Management (/api/v1/map-account) ---

type AccountsResponse struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Title      string     `json:"title"`
	NextAction string     `json:"nextAction"`
	Data       []string   `json:"data"`
}

type AccountMapResponse struct {
	Success    bool       `json:"success"`
	Message    string     `json:"message"`
	Title      string     `json:"title"`
	NextAction string     `json:"nextAction"`
}

// --- Data Caps & Bundle Limits ---

type SubmitLineLimitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Shake & Win (/api/v1/shake-and-win) ---

type ShakeAndWinData struct {
	Title             string         `json:"title"`
	SubTitle          string         `json:"subTitle"`
	Image             string         `json:"image"`
	Message           string         `json:"message"`
	Description       string         `json:"description"`
	ShakeAndWinAction string         `json:"shakeAndWinAction"`
	ActionButtons     []ActionButton `json:"actionButtons"`
}

type ShakeAndWinResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *ShakeAndWinData `json:"data"`
}

// --- Bill & Postpaid (/api/v1/top-up/bill-amount & pay-bill) ---

type BillInfoResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	DueDate string  `json:"dueDate"`
	Data    float64 `json:"data"`
}

// --- Search Engine (/api/v1/search & suggestions) ---

type SearchSuggestionsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type SearchItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ActionURL   string `json:"actionUrl"`
}

type SearchResultResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []SearchItem `json:"data"`
}

// --- International Services & Tariffs (/api/v1/international-services) ---

type InternationalCountryTariff struct {
	CountryCode string  `json:"countryCode"`
	CountryName string  `json:"countryName"`
	Flag        string  `json:"flag"`
	RatePerMin  float64 `json:"ratePerMin"`
	Currency    string  `json:"currency"`
}

type InternationalTariffResponse struct {
	Success bool                         `json:"success"`
	Message string                      `json:"message"`
	Data    []InternationalCountryTariff `json:"data"`
}

// --- Loyalty & Rewards (/api/v1/reward & /api/v1/eo) ---

type LoyaltyRewardItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Points      int    `json:"points"`
	ImageURL    string `json:"imageUrl"`
}

type LoyaltyRewardsResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []LoyaltyRewardItem `json:"data"`
}

type RedeemRewardResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	PID     string `json:"pid"`
}

// --- Data Line & Router Management (/api/v1/data-line & multi-line) ---

type DataLineInfo struct {
	MSISDN       string `json:"msisdn"`
	ICCID        string `json:"iccid"`
	Status       string `json:"status"`
	PackageName  string `json:"packageName"`
	RemainingVol string `json:"remainingVolume"`
}

type DataLineResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *DataLineInfo `json:"data"`
}

type MultiLineConnection struct {
	MSISDN   string `json:"msisdn"`
	Type     string `json:"type"`
	NickName string `json:"nickname"`
	Status   string `json:"status"`
}

type MultiLineResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []MultiLineConnection `json:"data"`
}

// --- Limits Inspection & Manage Lines (/api/v1/addon/datacap/limit & /api/v1/addon/share) ---

type LineLimitInfo struct {
	MSISDN      string  `json:"msisdn"`
	LimitMB     float64 `json:"limitMB"`
	ConsumedMB  float64 `json:"consumedMB"`
	RemainingMB float64 `json:"remainingMB"`
}

type LineLimitsResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []LineLimitInfo `json:"data"`
}

type SharedLineItem struct {
	MSISDN   string `json:"msisdn"`
	NickName string `json:"nickname"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type ManageLineResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []SharedLineItem `json:"data"`
}

// --- Ticket Detail & Form (/api/v1/resolution-center/{ticketNumber}) ---

type TicketDetailItem struct {
	TicketNumber string `json:"ticketNumber"`
	Category     string `json:"category"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	Description  string `json:"description"`
	Resolution   string `json:"resolution"`
}

type TicketDetailResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *TicketDetailItem `json:"data"`
}

type TicketFormField struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Options     []string `json:"options,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
}

type TicketFormResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []TicketFormField `json:"data"`
}

// --- FanZone & Gaming ---

type FanZoneHomeData struct {
	Title            string           `json:"title"`
	CompetitionID    string           `json:"competitionId"`
	Banner           string           `json:"banner"`
	ActionButton     PromotionButton  `json:"actionButton"`
	TermsLink        string           `json:"termsLink"`
	Status           string           `json:"status"`
}

type FanZoneHomeResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *FanZoneHomeData `json:"data"`
}

type FanZoneKickAndWinData struct {
	BackgroundImage string          `json:"backgroundImage"`
	Image           string          `json:"image"`
	AttemptsLeft    int             `json:"attemptsLeft"`
	MaxAttempts     int             `json:"maxAttempts"`
	ActionButton    PromotionButton `json:"actionButton"`
}

type FanZoneKickAndWinResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    *FanZoneKickAndWinData `json:"data"`
}

type FanZoneRewardData struct {
	RewardID    string `json:"rewardId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TicketID    string `json:"ticketId"`
	RewardType  string `json:"rewardType"`
	Value       string `json:"value"`
}

type FanZoneRewardResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *FanZoneRewardData `json:"data"`
}

type FanZoneLeaderItem struct {
	Rank     int    `json:"rank"`
	Nickname string `json:"nickname"`
	Score    int    `json:"score"`
	Avatar   string `json:"avatar"`
}

type FanZoneLeaderBoardData struct {
	CompetitionID string              `json:"competitionId"`
	Leaders       []FanZoneLeaderItem `json:"leaders"`
	UserRank      int                 `json:"userRank"`
	UserScore     int                 `json:"userScore"`
}

type FanZoneLeaderBoardResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    *FanZoneLeaderBoardData `json:"data"`
}

type FanZoneMatchPredictItem struct {
	MatchID   string `json:"matchId"`
	TeamA     string `json:"teamA"`
	TeamB     string `json:"teamB"`
	Date      string `json:"date"`
	ScoreA    int    `json:"scoreA"`
	ScoreB    int    `json:"scoreB"`
	Predicted bool   `json:"predicted"`
}

type FanZonePredictData struct {
	CompetitionID string                    `json:"competitionId"`
	Matches       []FanZoneMatchPredictItem `json:"matches"`
}

type FanZonePredictResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *FanZonePredictData `json:"data"`
}

type FanZonePrizeItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type FanZoneGrandPrizesData struct {
	Prizes []FanZonePrizeItem `json:"prizes"`
}

type FanZoneGrandPrizesResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    *FanZoneGrandPrizesData `json:"data"`
}

type FanZoneRewardsHistoryData struct {
	History []FanZoneRewardData `json:"history"`
}

type FanZoneRewardsHistoryResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    *FanZoneRewardsHistoryData `json:"data"`
}

type NicknameData struct {
	Nickname string `json:"nickname"`
	Count    int    `json:"count"`
	Max      int    `json:"max"`
}

type NicknameResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *NicknameData `json:"data"`
}

// --- Asiacell Partners & Discounts Directory ---

type EOCategory struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	TotalPartner int    `json:"totalPartner"`
}

type EOCategoryResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []EOCategory `json:"data"`
}

type EOCity struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	TotalCategory int    `json:"totalCategory"`
}

type EOCityResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []EOCity `json:"data"`
}

type EOPartner struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Logo          string   `json:"logo"`
	Address       string   `json:"address"`
	CityName      string   `json:"cityName"`
	CategoryName  string   `json:"categoryName"`
	Phone         string   `json:"phone"`
	Email         string   `json:"email"`
	Website       string   `json:"website"`
	DiscountValue float64  `json:"discountValue"`
	Discounts     []string `json:"discounts"`
	Lat           float64  `json:"lat"`
	Lng           float64  `json:"lng"`
	IsVisible     bool     `json:"isVisible"`
}

type EOPartnerResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []EOPartner `json:"data"`
}

type PartnerRegisterRequest struct {
	Name       string   `json:"name"`
	Phone      string   `json:"phone"`
	Email      string   `json:"email"`
	CityID     int      `json:"cityId"`
	CategoryID int      `json:"categoryId"`
	Address    string   `json:"address"`
	Discounts  []string `json:"discounts"`
}

// --- Yooz (Avocado) Platform ---

type BannerItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	URL      string `json:"url"`
	DeepLink string `json:"deeplink"`
}

type YoozHomeData struct {
	Title        string          `json:"title"`
	Balance      string          `json:"balance"`
	Expiry       string          `json:"expiry"`
	InternetMB   float64         `json:"internetMb"`
	ActionButton PromotionButton `json:"actionButton"`
	Banners      []BannerItem    `json:"banners"`
}

type YoozHomeResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *YoozHomeData `json:"data"`
}

type YoozAddOnListResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    *YoozBundlesScreenData `json:"data"`
}

type YoozBundlesScreenData struct {
	GroupID string           `json:"groupId"`
	Title   string           `json:"title"`
	Bundles []YoozPlanEntity `json:"bundles"`
}

type YoozPlanEntity struct {
	ID              int             `json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Price           string          `json:"price"`
	Volume          string          `json:"volume"`
	Icon            string          `json:"icon"`
	BackgroundImage string          `json:"backgroundImage"`
	Selected        bool            `json:"selected"`
	ActionDetail    *PromotionButton `json:"actionDetail,omitempty"`
}

type YoozBundlesResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    []YoozPlanEntity `json:"data"`
}

type YoozClassicPlansData struct {
	Plans []YoozPlanEntity `json:"plans"`
}

type YoozClassicPlansResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    *YoozClassicPlansData `json:"data"`
}

type YoozOmegaPlansData struct {
	Voucher string           `json:"voucher"`
	MSISDN  string           `json:"msisdn"`
	Plans   []YoozPlanEntity `json:"plans"`
}

type YoozOmegaPlansResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *YoozOmegaPlansData `json:"data"`
}

type YoozDataCapData struct {
	CurrentLimitMB int  `json:"currentLimitMb"`
	MaxLimitMB     int  `json:"maxLimitMb"`
	IsActive       bool `json:"isActive"`
}

type YoozDataCapResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *YoozDataCapData `json:"data"`
}

type YoozRewardData struct {
	RewardTitle string `json:"rewardTitle"`
	Points      int    `json:"points"`
	Status      string `json:"status"`
}

type YoozRewardResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    *YoozRewardData `json:"data"`
}

type YoozMigrateLineRequest struct {
	DOB    string `json:"dob"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type YoozMigrateOutHomeData struct {
	VID           *int              `json:"vid,omitempty"`
	Title         string            `json:"title"`
	Image         string            `json:"image"`
	Desc          string            `json:"desc"`
	ActionButton  *PromotionButton  `json:"actionButton,omitempty"`
	ActionButtons []PromotionButton `json:"actionButtons,omitempty"`
}

type YoozMigrateOutHomeResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    *YoozMigrateOutHomeData `json:"data"`
}

type YoozMigrateOutLocationItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type YoozMigrateOutLocationData struct {
	Title        string                       `json:"title"`
	Desc         string                       `json:"desc"`
	Items        []YoozMigrateOutLocationItem `json:"items"`
	ActionButton *PromotionButton             `json:"actionButton,omitempty"`
}

type YoozMigrateOutLocationResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    *YoozMigrateOutLocationData `json:"data"`
}

type YoozMigrateOutResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	NextAction string `json:"nextAction"`
	Title      string `json:"title"`
}


// --- Multi-Step Online Payment & Recharge ---

type RechargeNumberItem struct {
	MSISDN   string `json:"msisdn"`
	Selected bool   `json:"selected"`
	Label    string `json:"label"`
}

type RechargeNumberData struct {
	Title        string               `json:"title"`
	Items        []RechargeNumberItem `json:"items"`
	ActionButton PromotionButton      `json:"actionButton"`
}

type RechargeNumberResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *RechargeNumberData `json:"data"`
}

type RechargeTypeItem struct {
	TypeID   string `json:"typeId"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type RechargeTypeData struct {
	Title string             `json:"title"`
	Types []RechargeTypeItem `json:"types"`
}

type RechargeTypeResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *RechargeTypeData `json:"data"`
}

type OnlinePaymentProvider struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Logo     string `json:"logo"`
	Method   string `json:"method"`
	Currency string `json:"currency"`
}

type RechargeMethodData struct {
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	OnlinePayments []OnlinePaymentProvider `json:"onlinePayments"`
	ActionButton   PromotionButton         `json:"actionButton"`
}

type RechargeMethodResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *RechargeMethodData `json:"data"`
}

type OnlinePaymentPackage struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Label string  `json:"label"`
}

type OnlinePaymentData struct {
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	PaymentPackages []OnlinePaymentPackage `json:"paymentPackages"`
	OtherCardsLabel string                 `json:"otherCardsLabel"`
	ActionButton    PromotionButton        `json:"actionButton"`
}

type OnlinePaymentResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *OnlinePaymentData `json:"data"`
}

type RechargeConfirmationData struct {
	TransactionID string  `json:"transactionId"`
	MSISDN        string  `json:"msisdn"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	ReceiptDate   string  `json:"receiptDate"`
}

type RechargeConfirmationResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    *RechargeConfirmationData `json:"data"`
}

type PaymentSelectionData struct {
	RechargeType int                     `json:"rechargeType"`
	ForOthers    bool                    `json:"forOthers"`
	Methods      []OnlinePaymentProvider `json:"methods"`
}

type PaymentSelectionResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    *PaymentSelectionData `json:"data"`
}

// --- Home Dashboard v3, CYO & Quick Actions ---

type HomeDashboardData struct {
	Title        string           `json:"title"`
	Balance      string           `json:"balance"`
	Expiry       string           `json:"expiry"`
	InternetMB   float64          `json:"internetMb"`
	Minutes      int              `json:"minutes"`
	SMS          int              `json:"sms"`
	QuickActions []QuickActionItem `json:"quickActions"`
	Banners      []BannerItem     `json:"banners"`
}

type QuickActionItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Icon  string `json:"icon"`
	URL   string `json:"url"`
}

type HomeDashboardResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *HomeDashboardData `json:"data"`
}

type CYOBundlesData struct {
	GroupID      string           `json:"groupId"`
	Title        string           `json:"title"`
	CustomSteps  []string         `json:"customSteps"`
	MinInternet  int              `json:"minInternet"`
	MaxInternet  int              `json:"maxInternet"`
	MinMinutes   int              `json:"minMinutes"`
	MaxMinutes   int              `json:"maxMinutes"`
	ActionButton PromotionButton  `json:"actionButton"`
}

type CYOBundlesResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    *CYOBundlesData `json:"data"`
}

type ActiveOfferItem struct {
	OfferID     string `json:"offerId"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	Validity    string `json:"validity"`
	Description string `json:"description"`
}

type ActiveOffersData struct {
	Offers []ActiveOfferItem `json:"offers"`
}

type ActiveOffersResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    *ActiveOffersData `json:"data"`
}

// --- Feedback, Surveys, Video Tutorials & Voice of Customer ---

type SurveyQuestion struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Type     string   `json:"type"`
	Options  []string `json:"options,omitempty"`
}

type SurveyData struct {
	SurveyID  string           `json:"surveyId"`
	Title     string           `json:"title"`
	Questions []SurveyQuestion `json:"questions"`
}

type SurveyResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *SurveyData `json:"data"`
}

type VideoTutorialItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	VideoURL  string `json:"videoUrl"`
	Thumbnail string `json:"thumbnail"`
	Duration  string `json:"duration"`
}

type VideoTutorialsResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []VideoTutorialItem `json:"data"`
}

type VoCEntity struct {
	Username  string `json:"username"`
	NMFloID   string `json:"nmfloId"`
	CreatedAt int64  `json:"createdAt"`
}

type VoCResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    *VoCEntity `json:"data"`
}

// --- One-Yad & Epic Corporate ---

type OneYadData struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	ActionButton PromotionButton `json:"actionButton"`
	TotalDonated string          `json:"totalDonated"`
}

type OneYadResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *OneYadData `json:"data"`
}

type OneYadTeam struct {
	TeamID string `json:"teamId"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	City   string `json:"city"`
}

type OneYadTeamsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    []OneYadTeam `json:"data"`
}

type EpicLineItem struct {
	MSISDN   string `json:"msisdn"`
	Account  string `json:"account"`
	Status   string `json:"status"`
	PlanName string `json:"planName"`
}

type EpicLinesData struct {
	Lines []EpicLineItem `json:"lines"`
}

type EpicLinesResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    *EpicLinesData `json:"data"`
}

type EpicLineUsageData struct {
	MSISDN      string  `json:"msisdn"`
	RemainingMB float64 `json:"remainingMb"`
	TotalMB     float64 `json:"totalMb"`
	Minutes     int     `json:"minutes"`
	SMS         int     `json:"sms"`
	Expiry      string  `json:"expiry"`
}

type EpicLineUsageResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *EpicLineUsageData `json:"data"`
}

// --- Watch, Protected Payments, Asiaverse & Scanning ---

type WatchDashboardData struct {
	MSISDN     string  `json:"msisdn"`
	Balance    string  `json:"balance"`
	InternetMB float64 `json:"internetMb"`
	Minutes    int     `json:"minutes"`
	SMS        int     `json:"sms"`
	Expiry     string  `json:"expiry"`
}

type WatchDashboardResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    *WatchDashboardData `json:"data"`
}

type ProtectedPaymentStatusData struct {
	TransactionID string  `json:"transactionId"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

type ProtectedPaymentStatusResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    *ProtectedPaymentStatusData `json:"data"`
}

type AsiaverseHomeData struct {
	Title       string `json:"title"`
	Banner      string `json:"banner"`
	Status      string `json:"status"`
	RedirectURL string `json:"redirectUrl"`
}

type AsiaverseHomeResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    *AsiaverseHomeData `json:"data"`
}

type ScanToWinData struct {
	CampaignID string `json:"campaignId"`
	Title      string `json:"title"`
	ScanLimit  int    `json:"scanLimit"`
	IsActive   bool   `json:"isActive"`
}

type ScanToWinResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    *ScanToWinData `json:"data"`
}

type CDRSummaryData struct {
	TotalTransfersIn  float64 `json:"totalTransfersIn"`
	TotalTransfersOut float64 `json:"totalTransfersOut"`
	TransfersCount    int     `json:"transfersCount"`
}

type CDRSummaryResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    *CDRSummaryData `json:"data"`
}

type UserInterestsResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}



