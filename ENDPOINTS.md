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
| **16** | `GET` | `/api/v1/transaction/transfer` | `client.GetTransferHistory(ctx)` | سجل وتاريخ عمليات تحويل الرصيد السابقة على الخط (المحفظة) |
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
| **28** | `GET` | `/api/v2/vanity/classes` | `client.GetVanityClasses(ctx)` | استعراض فئات وتصنيفات الأرقام المميزة (الماسية، الذهبية، الفضية) وأسعارها |
| **29** | `GET` | `/api/v2/vanity` | `client.SearchVanityNumbers(ctx, pat, class, p, l)` | البحث في الأرقام المميزة المتاحة للبيع بنمط أو فئة معينة |
| **30** | `GET` | `/api/v2/vanity/{msisdn}/detail` | `client.GetVanityDetail(ctx, msisdn)` | تفاصيل وسعر وشروط حجز رقم مميز محدد |
| **31** | `POST` | `/api/v2/vanity` | `client.ReserveVanityNumber(ctx, msisdn, classId)` | حجز الرقم المميز مباشرة باسم المشترك وتوليد معرف العملية |
| **32** | `POST` | `/api/v1/addon/send-as-gift` | `client.SendGiftAddon(ctx, addonId, to)` | شراء باقة إنترنت أو اتصالات وإهداؤها لرقم آخر بخصم من الرصيد |
| **33** | `GET` | `/api/v1/resolution-center/categories` | `client.GetTicketCategories(ctx)` | أقسام وتصنيفات الشكاوى الفنية المعتمدة في آسياسيل |
| **34** | `GET` | `/api/v1/resolution-center` | `client.GetTickets(ctx)` | سجل تذاكر الشكاوى المفتوحة وتحديثات مسار المعالجة |
| **35** | `POST` | `/api/v1/resolution-center` | `client.CreateTicket(ctx, catId, desc)` | فتح وإرسال تذكرة شكوى رسمية جديدة لإدارة العمليات والدعم |
| **36** | `GET` | `/api/v1/compensation` | `client.CheckCompensation(ctx)` | فحص استحقاق الخط للتعويضات الرسمية (جيجابايت أو رصيد مجاني) |
| **37** | `GET` | `/api/v1/yooz-mgm` | `client.GetYoozMGM(ctx)` | استخراج كود الإحالة وإحصائيات دعوة الأصدقاء لخطوط Yooz |
| **38** | `POST` | `/api/v1/yooz-mgm/apply-code` | `client.ApplyYoozMGMCode(ctx, code)` | تفعيل كود دعوة للحصول على البونص والمكافآت المجانية |
| **39** | `GET` | `/api/v2/e-voucher/packages` | `client.GetEVoucherPackages(ctx)` | تصفح كروت الألعاب والشحن الرقمي (PUBG, PlayStation, iTunes) |

---

## 2. الشرح التقني المفصل للعمليات

### 2.1 كشف الحساب والتحقق التلقائي من تحويلات الرصيد (CDR & Transfer Verification)

#### المسار:
<div dir="ltr" align="left">

```http
GET /api/v1/cdr/detail?type=btransfer&page={page}&limit={limit}&lang=ar
```

</div>

- **الغرض**: الاستعلام من خوادم آسياسيل عن السجل الحقيقي لتحويلات الرصيد (Call Detail Records) الواردة والصادرة على الشريحة.

**الهيدرز المطلوبة (Headers):**
<div dir="ltr" align="left">

```http
Authorization: Bearer <access_token>
DeviceId: <UUID-v4>
X-ODP-API-KEY: 1ccbc4c913bc4ce785a0a2de444aa0d6
```

</div>

**استجابة الخادم النموذجية (Server Response):**
<div dir="ltr" align="left">

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

</div>

#### تفعيل كشف الحساب (CDR Activation):
- `POST /api/v1/cdr/send-otp`: إرسال كود OTP برسالة SMS لتفعيل صلاحية الوصول للجلسة.
- `POST /api/v1/cdr/confirm`: تأكيد كود الـ OTP.

#### سجل تحويلات المحفظة السريعة (Transfer History):
- `GET /api/v1/transaction/transfer`: استعلام سجل عمليات تحويل الرصيد السابقة الموثقة في محفظة الحساب (My Pocket) من خوادم آسياسيل مباشرة مستخرجة من كود تطبيق آسياسيل الرسمي (`cw7.java`).

<div dir="ltr" align="left">

```json
{
  "code": 200,
  "message": "success",
  "success": true,
  "data": [
    {
      "type": "TRANSFER",
      "msisdn": "07701234567",
      "receiverMsisdn": "07744298878",
      "createdAt": 1726068223000,
      "amount": 1000.0
    }
  ]
}
```

</div>

#### خوارزمية التحقق التلقائي المدمجة (`VerifyIncomingTransfer`):
- تقوم دالة الـ SDK بالاستعلام المزدوج الذكي: فحص سجل الـ CDR السحابي (`/api/v1/cdr/detail?type=btransfer`) مع الرجوع لسجل المحفظة (`/api/v1/transaction/transfer`)، وتطابق رقم هاتف المرسل (بمقارنة آخر 9 أرقام لتجاوز اختلافات البادئات `077` أو `96477`)، وتتأكد أن الحوالة واردة (إيجابية وليست سالبة) وأن القيمة مساوية أو أكبر من المطلوب، لمنع الاحتيال وضمان الإيداع التلقائي فورياً بدون تدخل يدوي للأدمن.

---

### 2.2 سوق الأرقام المميزة (Vanity VIP Numbers)

- `GET /api/v2/vanity/classes`: استعراض فئات الأرقام المميزة (الماسية، الذهبية، الفضية، البرونزية).
- `GET /api/v2/vanity?msisdn={pattern}&classId={id}&page={p}&limit={l}`: البحث عن أرقام مميزة بنمط محدد وأسعارها.

#### حجز الرقم المميز:
<div dir="ltr" align="left">

```http
POST /api/v2/vanity
```

```json
{
  "msisdn": "07700001111",
  "classId": "1"
}
```

</div>

---

### 2.3 إهداء الباقات للغير (Send Addon as a Gift)

- **الغرض**: شراء باقة إنترنت أو مكالمات وإرسالها كهدية لأي رقم آسياسيل آخر، مع الخصم المباشر من رصيد الشريحة المرسلة.

<div dir="ltr" align="left">

```http
POST /api/v1/addon/send-as-gift
```

```json
{
  "addOnId": 142,
  "receiverMsisdn": "07701234567"
}
```

</div>

---

### 2.4 نظام الشكاوى والتذاكر الفنية (Resolution Center)

- `GET /api/v1/resolution-center/categories`: جلب تصنيفات المشاكل والشكاوى المعتمدة.
- `GET /api/v1/resolution-center`: استعراض التذاكر السابقة المفتوحة ومسار معالجتها.

#### فتح تذكرة دعم فني جديدة:
<div dir="ltr" align="left">

```http
POST /api/v1/resolution-center
```

```json
{
  "category": "cat_network",
  "description": "انقطاع مفاجئ في إشارة الـ 4G في منطقة المنصور"
}
```

</div>

---

### 2.5 نظام التعويضات التلقائي (Compensation System)

<div dir="ltr" align="left">

```http
GET /api/v1/compensation
```

</div>

- **الغرض**: فحص استحقاق الخط للتعويضات الرسمية المعتمدة من آسياسيل عند وجود أعطال شبكة عامة، واستلام باقات إنترنت أو رصيد مجاني.

---

### 2.6 خطوط الشباب Yooz (MGM Referral Program)

- `GET /api/v1/yooz-mgm`: استخراج كود الإحالة ورابط المشاركة وعدد الإحالات الناجحة.

#### تفعيل كود دعوة:
<div dir="ltr" align="left">

```http
POST /api/v1/yooz-mgm/apply-code
```

```json
{
  "promoCode": "YOOZ2026"
}
```

</div>

---

### 2.7 كروت الألعاب والشحن الرقمي (E-Vouchers)

<div dir="ltr" align="left">

```http
GET /api/v2/e-voucher/packages?recharge-type=1
```

</div>

- **الغرض**: استعراض كروت الألعاب والتطبيقات المتاحة للشراء برصيد الهاتف (PUBG, PlayStation, iTunes, etc.).

---

### 2.8 تسجيل الدخول وإدارة الجلسات (Authentication)

<div dir="ltr" align="left">

```http
POST /api/v1/auth/phone
POST /api/v1/auth/login-passcode
POST /api/v1/auth/refresh-token
```

</div>

- **الغرض**: تسجيل الدخول برقم الهاتف عبر رمز التحقق SMS واستخراج وتجديد التوكنات تلقائياً.

---

## 3. الهيدرز الرسمية المطلوبة (Standard HTTP Headers)

جميع الطلبات المرسلة تتضمن الهيدرز الرسمية المعتمدة من بوابة آسياسيل:

<div dir="ltr" align="left">

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

</div>
