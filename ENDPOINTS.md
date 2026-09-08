<div align="center">

<img src="assets/asiacell.svg" alt="Asiacell Logo" width="200" />

# Verified API Endpoints Reference

**Complete specification of the 24 reverse-engineered, production-tested Asiacell Iraq HTTP REST APIs.**

[![Specification](https://img.shields.io/badge/Status-100%25%20Verified%20(24%2F24)-E7242A?style=flat-square)](ENDPOINTS.md)
[![Protocol](https://img.shields.io/badge/Protocol-HTTPS%2FREST-18181b?style=flat-square)](ENDPOINTS.md)
[![Client](https://img.shields.io/badge/SDK-asiacell--go-18181b?style=flat-square)](https://github.com/FLEX-GHOST/asiacell-go)

<br />

All endpoints below have been verified directly against Asiacell production servers (`app.asiacell.com`), passing unit tests, live integration tests, and error handling checks.

</div>

---

## 1. Endpoints Matrix Overview

| # | HTTP Method | Endpoint Path | Go SDK Method | Description |
| :---: | :---: | :--- | :--- | :--- |
| **01** | `POST` | `/api/v1/auth/phone` | `client.Login(ctx, phone)` | Request 6-digit SMS OTP verification code |
| **02** | `POST` | `/api/v1/auth/login-passcode` | `client.VerifySMS(ctx, pid, code)` | Verify SMS OTP and exchange for Access & Refresh tokens |
| **03** | `POST` | `/api/v1/auth/refresh-token` | `client.RefreshToken(ctx)` | Silently renew expired Bearer tokens |
| **04** | `GET` | `/api/v1/auth/captcha` | `client.SolveCaptchaOCR(data)` | Fetch Captcha challenge image when triggered |
| **05** | `GET` | `/api/v1/profile` | `client.GetProfile(ctx)` | Current account balance, validity, remaining data, calls, SMS |
| **06** | `GET` | `/api/v1/profile/view2` | `client.GetProfileDetails(ctx)` | Official subscriber profile: full name, email, birthday, avatar |
| **07** | `GET` | `/api/v1/profile-img` | `client.GetProfileDetails(ctx)` | Profile avatar image stream |
| **08** | `GET` | `/api/v2/4gcontent` | `client.GetUnlimited4GBundles(ctx)` | 4G unlimited internet bundles (Daily, Weekly, Monthly) |
| **09** | `GET` | `/api/v1/products/` | `client.GetSpecialOffers(ctx)` | Personalized line discounts and promotional offers |
| **10** | `GET` | `/api/v1/addon/tags/` | `client.GetAddonTags(ctx)` | Addon categories (Internet, Roaming, Calls, Extras) |
| **11** | `GET` | `/api/v1/addon/summary/` | `client.GetAddonSummary(ctx, tagID)` | Detailed catalog and pricing for specific addon category |
| **12** | `POST` | `/api/v1/addon/subscribe` | `client.SubscribeAddon(ctx, addonID)` | Instant balance-deducted package activation via API |
| **13** | `POST` | `/api/v1/credit/transfer` | `client.StartCreditTransfer(ctx, to, amt)` | Initiate P2P credit transfer and trigger SMS confirmation |
| **14** | `POST` | `/api/v1/credit/transfer/confirm` | `client.ConfirmCreditTransfer(ctx, pid, code)`| Finalize credit transfer with 6-digit SMS confirmation code |
| **15** | `GET` | `/api/v1/credit/transfer/history` | `client.GetTransferHistory(ctx)` | Historical P2P credit transfers ledger |
| **16** | `POST` | `/api/v1/top-up` | `client.RechargeVoucher(ctx, to, code, type)` | Recharge balance via 14-digit voucher card (regular or data) |
| **17** | `GET` | `/api/v1/top-up/bill-amount` | `client.GetRechargeHistory(ctx)` | Historical voucher card recharge transactions |
| **18** | `GET` | `/api/v2/services/management` | `client.GetServicesManagement(ctx)` | Subscribed services & Internet Balance Protection (*223#) status |
| **19** | `POST` | `/api/v2/services/action` | `client.ExecuteServiceAction(ctx, svc, act)` | Toggle balance protection or execute service actions |
| **20** | `GET` | `/api/v1/digital-services` | `client.GetDigitalServices(ctx)` | Entertainment, streaming, and value-added digital services |
| **21** | `GET` | `/api/v1/partners/cities` | `client.GetCities(ctx)` | List of all Iraqi governorates and internal city IDs |
| **22** | `GET` | `/api/v1/shops` | `client.GetCityShops(ctx, cityName)` | Certified Asiacell branches: addresses, GPS coordinates, open status |
| **23** | `GET` | `/api/v1/spin-wheel` | `client.GetSpinWheelStatus(ctx)` | Daily spin wheel eligibility and cooldown status |
| **24** | `POST` | `/api/v1/spin-wheel/play` | `client.PlaySpinWheel(ctx)` | Execute spin wheel and claim daily internet/credit reward |

---

## 2. Detailed Technical Breakdown

### 2.1 Authentication & Session Management

#### `POST /api/v1/auth/phone`
- **Purpose**: Initiates user login by sending a 6-digit OTP code to the subscriber's phone number.
- **Headers**:
  - `Content-Type: application/json`
  - `User-Agent: Mozilla/5.0 (Linux; Android 14; Mobile)`
  - `x-device-id: <uuid>`
- **Request Body**:
  ```json
  {
    "username": "07701234567"
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "success": true,
    "message": "OTP has been sent",
    "data": {
      "pid": "1af574e3-a8cc-48d8-acec-f9aee2d98e60"
    }
  }
  ```
- **Captcha Trigger (`428 / 400`)**:
  When requested too frequently, returns `ErrCaptchaRequired` with a base64 captcha image. The SDK automatically mitigates this via `client.RotateDeviceID()`.

---

#### `POST /api/v1/auth/login-passcode`
- **Purpose**: Verifies the SMS OTP and issues JWT tokens.
- **Request Body**:
  ```json
  {
    "pid": "1af574e3-a8cc-48d8-acec-f9aee2d98e60",
    "passcode": "123456"
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "success": true,
    "data": {
      "access_token": "eyJhbGciOiJIUzUxMiJ9...",
      "refresh_token": "eyJhbGciOiJIUzUxMiJ9...",
      "userId": "7810955"
    }
  }
  ```

---

#### `POST /api/v1/auth/refresh-token`
- **Purpose**: Automatically renews the `access_token` using the long-lived `refresh_token`. Handled transparently by `client.AutoRefreshToken(ctx)`.

---

### 2.2 Account & Balance Queries

#### `GET /api/v1/profile?lang=ar`
- **Purpose**: Returns balance amount, active bundle expiration, remaining megabytes, and call allowances.
- **Headers**:
  - `Authorization: Bearer <access_token>`
- **Parsed Fields**:
  - `Balance`: Current credit in Iraqi Dinars (IQD).
  - `Validity`: Expiration date of account validity.
  - `RemainingData`: Formatted active data package volume (e.g. `20 MB (Data(SpinWheel))`).
  - `RemainingCalls`: Voice package minutes.
  - `RemainingSMS`: SMS allowances.

---

#### `GET /api/v1/profile/view2?lang=ar`
- **Purpose**: Retrieves official subscriber account details.
- **Parsed Fields**:
  - Full Name, Registered Phone Number, Email, Birthday, Profile Picture URL.

---

### 2.3 Bundles & 4G Internet

#### `GET /api/v2/4gcontent?lang=ar`
- **Purpose**: Lists unlimited 4G internet bundles available for subscription.
- **Sample Output**:
  - `باقة بلا حدود يومية (Daily Unlimited)` — Price: 2,500 IQD
  - `باقة بلا حدود أسبوعية (Weekly Unlimited)` — Price: 10,000 IQD
  - `باقة بلا حدود شهرية (Monthly Unlimited)` — Price: 40,000 IQD

#### `POST /api/v1/addon/subscribe`
- **Purpose**: Activates a specific bundle directly using the line's main balance.
- **Request Body**:
  ```json
  {
    "addon_id": "ODP_4G_DAILY_UNLIMITED"
  }
  ```

---

### 2.4 Credit Transfer & Payment Verification

#### `POST /api/v1/credit/transfer`
- **Purpose**: Starts a P2P credit transfer to any Asiacell number.
- **Request Body**:
  ```json
  {
    "receiver": "07709876543",
    "amount": 5000
  }
  ```
- **Behavior**: Asiacell sends an SMS challenge to the sender containing a 6-digit confirmation code.

#### `POST /api/v1/credit/transfer/confirm`
- **Purpose**: Finalizes the transfer with the OTP code.
- **Request Body**:
  ```json
  {
    "pid": "<transfer_pid>",
    "passcode": "123456"
  }
  ```

#### Local Ledger Verification (`VerifyIncomingTransfer`)
- **SDK Algorithm**: Checks incoming transactions, compares timestamps, matches sender phone and minimum transferred amount, and enforces deduplication so merchants and automated bots never process the same payment twice.

---

### 2.5 Voucher Cards (Top-Up)

#### `POST /api/v1/top-up`
- **Purpose**: Recharges line balance using 14-digit scratch voucher cards.
- **Supported Types**:
  - `normal`: Regular credit card.
  - `data`: Dedicated internet card.
- **Request Body**:
  ```json
  {
    "phone": "07701234567",
    "voucher": "12345678901234",
    "type": "normal"
  }
  ```

---

### 2.6 Balance Protection & Services

#### `GET /api/v2/services/management?lang=ar`
- **Purpose**: Inspects active services and internet balance protection (equivalent to dialing `*223#`).
- **Feature**: Protects main credit from being depleted after an internet bundle expires.

---

### 2.7 Certified Shops & Governorates

#### `GET /api/v1/partners/cities?lang=ar`
- **Purpose**: Returns all Iraqi governorates (`Baghdad`, `Basra`, `Erbil`, `Najaf`, `Karbala`, `Sulaymaniyah`, `Kirkuk`, etc.).

#### `GET /api/v1/shops?cities=<CityName>&lang=ar`
- **Purpose**: Real-time listing of certified Asiacell customer care centers and branches.
- **Parsed Data**:
  - Branch Name, Street Address, Working Hours, Current Open/Closed status, GPS Coordinates (Latitude / Longitude).

---

### 2.8 Spin Wheel & Shukran Loyalty

#### `GET /api/v1/spin-wheel?lang=ar`
- **Purpose**: Checks whether the daily free spin wheel is available or still in cooldown.

#### `POST /api/v1/spin-wheel/play`
- **Purpose**: Plays the daily spin wheel and claims the reward (free megabytes or bonus credit).

---

## 3. Standard HTTP Headers

Every request made by `asiacell-go` includes the official headers expected by Asiacell's gateway:

```http
User-Agent: Mozilla/5.0 (Linux; Android 14; Mobile) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Mobile Safari/537.36
Accept: application/json, text/plain, */*
Accept-Language: ar,en-US;q=0.9,en;q=0.8
Origin: https://app.asiacell.com
Referer: https://app.asiacell.com/
x-device-id: <UUIDv4>
Authorization: Bearer <JWT_ACCESS_TOKEN>
```
