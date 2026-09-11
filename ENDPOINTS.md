<div align="center">

<img src="assets/asiacell.svg" alt="شعار آسياسيل" width="200" />

# مرجع نقاط نهاية واجهة برمجة تطبيقات آسياسيل (Asiacell API Endpoints)

**المواصفات الكاملة لجميع نقاط نهاية HTTP REST API لشركة آسياسيل العراق، مستخرجة ومفحوصة بالكامل مع خوادم الإنتاج.**

[![Specification](https://img.shields.io/badge/Specification-100%25%20Verified-E7242A?style=flat-square)](ENDPOINTS.md)
[![Protocol](https://img.shields.io/badge/Protocol-HTTPS%2FREST-18181b?style=flat-square)](ENDPOINTS.md)
[![Client](https://img.shields.io/badge/Go%20Client-asiacell--go-18181b?style=flat-square)](https://github.com/FLEX-GHOST/asiacell-go)

<br />

جميع نقاط النهاية الموضحة أدناه تم فحصها والتحقق منها مباشرة مقابل خوادم آسياسيل الرسمية (`odpapp.asiacell.com` و `app.asiacell.com`) مع اجتياز اختبارات الوحدة والاختبارات الحية ومعالجة الأخطاء.

</div>

---

## 1. جدول نقاط النهاية المعتمدة (Endpoints Matrix)

| # | طريقة HTTP | مسار النقطة (Endpoint Path) | دالة Go SDK | الوصف التفصيلي |
| :---: | :---: | :--- | :--- | :--- |
| **01** | `POST` | `/api/v1/auth/phone` | `client.Login(ctx, phone)` | إرسال رمز التحقق OTP المكون من 6 أرقام برسالة SMS لتسجيل الدخول |
| **02** | `POST` | `/api/v1/auth/login-passcode` | `client.VerifySMS(ctx, pid, code)` | التحقق من كود الـ SMS واستبداله بتوكنات الجلسة (Access & Refresh) |
| **03** | `POST` | `/api/v1/auth/refresh-token` | `client.RefreshToken(ctx)` | تجديد توكن الوصول المنتهي تلقائياً في الخلفية بدون طلب كود مجدداً |
| **04** | `GET` | `/api/v1/auth/captcha` | `client.SolveCaptchaOCR(data)` | جلب صورة اختبار الكابتشا في حال فرضها لمنع الطلبات المتكررة |
| **05** | `GET` | `/api/v1/profile` | `client.GetProfile(ctx)` | استعلام رصيد الحساب الحالي، مدة الصلاحية، المتبقي من الإنترنت والمكالمات والرسائل |
| **06** | `GET` | `/api/v1/profile/view2` | `client.GetProfileDetails(ctx)` | تفاصيل الملف الشخصي للمشترك: الاسم الكامل، البريد، تاريخ الميلاد، ورابط الصورة |
| **07** | `GET` | `/api/v1/profile-img` | `client.GetProfileDetails(ctx)` | دفق صورة الحساب الرمزية للمستخدم (Avatar) |
| **08** | `GET` | `/api/v2/4gcontent` | `client.GetUnlimited4GBundles(ctx)` | باقات الإنترنت المفتوح 4G (يومية، أسبوعية، شهرية) وأسعارها |
| **09** | `GET` | `/api/v1/products/` | `client.GetSpecialOffers(ctx)` | العروض الخاصة والتخفيضات المخصصة لرقم الهاتف |
| **10** | `GET` | `/api/v1/addon/tags/` | `client.GetAddonTags(ctx)` | أقسام وتصنيفات الباقات الإضافية (إنترنت، تجوال، اتصالات، خدمات) |
| **11** | `GET` | `/api/v1/addon/summary/` | `client.GetAddonSummary(ctx, tagID)` | تفاصيل الباقات والأسعار لقسم محدد من الباقات |
| **12** | `POST` | `/api/v1/addon/subscribe` | `client.SubscribeAddon(ctx, addonID)` | تفعيل واشتراك فوري في الباقة مع خصم قيمتها من الرصيد مباشرة |
| **13** | `GET` | `/api/v1/cdr/detail?type=btransfer` | `client.GetCDRTransferHistory(ctx, page, limit)` | جلب سجل كشف الحساب لتحويلات الرصيد الواردة والصادرة مع رقم المرسل والمبلغ والتاريخ |
| **14** | `POST` | `/api/v1/cdr/send-otp` | `client.SendCDROTP(ctx)` | طلب رمز OTP لتفعيل خدمة كشف الحساب (CDR) على الرقم للجلسة |
| **15** | `POST` | `/api/v1/cdr/confirm` | `client.ConfirmCDROTP(ctx, otp)` | تأكيد رمز OTP لتفعيل صلاحية الوصول لكشف الحساب وسجلات التحويل |
| **16** | `—` | *(خوارزمية داخلية)* | `client.VerifyIncomingTransfer(ctx, phone, minAmt)` | **التحقق التلقائي الفوري من وصول تحويل رصيد من رقم مرسل ومطابقته آلياً** |
| **17** | `POST` | `/api/v1/credit-transfer/start` | `client.StartCreditTransfer(ctx, to, amt)` | بدء تحويل رصيد من الشريحة لرقم آخر وتوليد معرف العملية PID |
| **18** | `POST` | `/api/v1/credit-transfer/do-transfer` | `client.ConfirmCreditTransfer(ctx, pid, code)` | تأكيد تحويل الرصيد بإدخال رمز التحقق المرسل إلى الهاتف |
| **19** | `POST` | `/api/v1/top-up` | `client.RechargeVoucher(ctx, to, code, type)` | شحن وتعبئة كارت آسياسيل عبر كود الكارت (13-14 رقماً) عادي أو إنترنت |
| **20** | `GET` | `/api/v1/transaction/recharge` | `client.GetRechargeHistory(ctx)` | سجل وتاريخ عمليات شحن الكروت السابقة على الخط |
| **21** | `GET` | `/api/v2/services/management` | `client.GetServicesManagement(ctx)` | إدارة الخدمات المشترك بها وحالة حماية الرصيد عند انتهاء الباقة (*223#) |
| **22** | `POST` | `/api/v2/services/action` | `client.ExecuteServiceAction(ctx, svc, act)` | تفعيل أو إيقاف حماية الرصيد وإلغاء الخدمات المزعجة (299/4151/300) |
| **23** | `GET` | `/api/v1/digital-services` | `client.GetDigitalServices(ctx)` | الخدمات الترفيهية والقيمة المضافة الرقمية |
| **24** | `GET` | `/api/v1/partners/cities` | `client.GetCities(ctx)` | قائمة بجميع المحافظات والمدن العراقية ومعرفاتها |
| **25** | `GET` | `/api/v1/shops` | `client.GetCityShops(ctx, cityName)` | مراكز وفروع آسياسيل المعتمدة: العناوين، الإحداثيات، وحالة الفتح/الإغلاق |
| **26** | `GET` | `/api/v1/spin-wheel` | `client.GetSpinWheelStatus(ctx)` | حالة عجلة الحظ اليومية وما إذا كانت متاحة للدوران |
| **27** | `POST` | `/api/v1/spin-wheel/play` | `client.PlaySpinWheel(ctx)` | تدوير عجلة الحظ واستلام الجائزة اليومية (ميغابايت أو رصيد مجاني) |

---

## 2. الشرح التقني المفصل للعمليات

### 2.1 كشف الحساب والتحقق التلقائي من تحويلات الرصيد (CDR & Transfer Verification)

#### `GET /api/v1/cdr/detail?type=btransfer&page={page}&limit={limit}&lang=ar`
- **الغرض**: الاستعلام من خوادم آسياسيل عن السجل الحقيقي لتحويلات الرصيد (Call Detail Records) الواردة والصادرة على الشريحة.
- **الهيدرز**:
  - `Authorization: Bearer <access_token>`
  - `DeviceId: <UUID-v4>`
  - `X-ODP-API-KEY: 1ccbc4c913bc4ce785a0a2de444aa0d6`
- **استجابة الخادم النموذجية `200 OK`**:
  ```json
  {
    "code": 200,
    "message": "success",
    "success": true,
    "data": {
      "total": 1,
      "data": [
        {
          "amount": "1000 IQD",
          "unit": "TRANSFERS",
          "title": "تحويل الرصيد",
          "subTitle": "7744298878",
          "description": "١١/٠٩/٢٠٢٦ ٠٦:٢٣:٤٣"
        }
      ]
    }
  }
  ```
  - `subTitle`: رقم هاتف المرسل (في الحوالات الواردة).
  - `amount`: المبلغ بالدينار العراقي (المبلغ الإيجابي للحوالة الواردة، وبداية سالب `-` للحوالة الصادرة).
  - `description`: تاريخ ووقت العملية.

#### `POST /api/v1/cdr/send-otp`
- **الغرض**: إرسال كود OTP عبر رسالة نصية لتفعيل خدمة كشف الحساب لأول مرة على الجلسة.
- **البايلود**: `{}`

#### `POST /api/v1/cdr/confirm`
- **الغرض**: تأكيد كود الـ OTP وتفعيل الوصول لكشف الحساب.
- **البايلود**:
  ```json
  {
    "code": "676076"
  }
  ```

#### خوارزمية التحقق التلقائي (`VerifyIncomingTransfer`)
- **الآلية**: تقوم الدالة بالاستعلام المباشر من سجل الـ CDR عن آخر 30 عملية، وتطابق رقم الهاتف المرسل (بمقارنة آخر 9 أرقام لتجاوز اختلافات الصيغ الدولية والمحلية `077` أو `96477`)، وتتحقق من أن الحوالة واردة (ليست سالبة) ومن القيمة المالية وتاريخ العملية، لمنع الاحتيال وضمان الإيداع والتفعيل الفوري بدون أي تدخل بشري.

---

### 2.2 تسجيل الدخول وإدارة الجلسات (Authentication)

#### `POST /api/v1/auth/phone`
- **الغرض**: بدء تسجيل الدخول بإرسال رمز OTP مكون من 6 أرقام إلى رقم الهاتف.
- **البايلود**:
  ```json
  {
    "username": "07701234567"
  }
  ```

#### `POST /api/v1/auth/login-passcode`
- **الغرض**: تأكيد رمز الـ OTP واستخراج التوكنات.
- **البايلود**:
  ```json
  {
    "pid": "1af574e3-a8cc-48d8-acec-f9aee2d98e60",
    "passcode": "123456"
  }
  ```
- **الاستجابة**:
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

#### `POST /api/v1/auth/refresh-token`
- **الغرض**: تجديد توكن الوصول تلقائياً عند انتهائه باستخدام الـ `refresh_token`.

---

### 2.3 الاستعلام عن الرصيد والحساب (Profile & Balance)

#### `GET /api/v1/profile?lang=ar`
- **الغرض**: استرجاع الرصيد الحالي بالدينار العراقي، تاريخ انتهاء صلاحية الخط، وباقات الإنترنت والدقائق الفعالة.
- **الحقول المستخرجة**:
  - `Balance`: الرصيد بالدينار العراقي (IQD).
  - `Validity`: تاريخ انتهاء صلاحية الخط.
  - `RemainingData`: حجم الإنترنت المتبقي.
  - `RemainingCalls`: الدقائق المتبقية.
  - `RemainingSMS`: الرسائل المتبقية.

#### `GET /api/v1/profile/view2?lang=ar`
- **الغرض**: تفاصيل المشترك الرسمية (الاسم، الهاتف، البريد، الصورة الشخصية).

---

### 2.4 باقات الإنترنت والعروض (Bundles & Offers)

#### `GET /api/v2/4gcontent?lang=ar`
- **الغرض**: استعراض باقات الـ 4G غير المحدودة المتاحة للاشتراك.
- **الباقات النموذجية**:
  - باقة بلا حدود يومية (2,500 د.ع)
  - باقة بلا حدود أسبوعية (10,000 د.ع)
  - باقة بلا حدود شهرية (40,000 د.ع)

#### `POST /api/v1/addon/subscribe`
- **الغرض**: الاشتراك وتفعيل الباقة فورياً بخصم مباشر من رصيد الخط.

---

### 2.5 كروت الشحن (Top-Up)

#### `POST /api/v1/top-up`
- **الغرض**: شحن الخط عبر كود كارت التعبئة (13 إلى 14 رقماً).
- **الأنواع المدعومة**:
  - `normal` / `1`: كارت رصيد عادي.
  - `data` / `2`: كارت إنترنت مخصص.
- **البايلود**:
  ```json
  {
    "msisdn": "",
    "rechargeType": 1,
    "voucher": "12345678901234"
  }
  ```

---

### 2.6 حماية الرصيد وإلغاء الخدمات المزعجة

#### `GET /api/v2/services/management?lang=ar`
- **الغرض**: فحص حالة خدمة حماية الرصيد (المكافئة للرمز `*223#`).
- **الميزة**: إيقاف استهلاك الرصيد الأساسي بعد نفاد باقة الإنترنت لحماية رصيدك من الضياع.

#### `POST /api/v2/services/action`
- **الغرض**: تفعيل أو إلغاء حماية الرصيد وإلغاء اشتراكات الخدمات المزعجة.

---

### 2.7 الفروع والمراكز المعتمدة وعجلة الحظ

#### `GET /api/v1/partners/cities?lang=ar`
- **الغرض**: قائمة المحافظات العراقية (بغداد، البصرة، أربيل، النجف، كربلاء، السليمانية، كركوك...).

#### `GET /api/v1/shops?cities=<CityName>&lang=ar`
- **الغرض**: مراكز خدمة الزبائن وفروع آسياسيل المعتمدة متضمنة العنوان وساعات العمل والإحداثيات الجغرافية (GPS).

#### `GET /api/v1/spin-wheel` و `POST /api/v1/spin-wheel/play`
- **الغرض**: فحص وتدوير عجلة الحظ اليومية المجانية واستلام المكافآت.

---

## 3. الهيدرز الرسمية المطلوبة (Standard HTTP Headers)

جميع الطلبات المرسلة تتضمن الهيدرز الرسمية المعتمدة من بوابة آسياسيل:

```http
Accept: application/json, text/plain, */*
User-Agent: okhttp/5.0.0-alpha.2
DeviceId: <UUID-v4>
x-device-id: <UUID-v4>
X-ODP-API-KEY: 1ccbc4c913bc4ce785a0a2de444aa0d6
X-OS-Version: 15
X-ODP-APP-VERSION: 4.2.5
X-Device-Type: [Android][INFINIX][Infinix X6871 15][VANILLA_ICE_CREAM][HMS][4.2.5:90000256]
X-FROM-APP: odp
X-ODP-CHANNEL: mobile
X-SCREEN-TYPE: MOBILE
Authorization: Bearer <JWT_ACCESS_TOKEN>
```
