<div align="center">

<img src="assets/asiacell.svg" alt="Asiacell Logo" width="240" />

# asiacell-go

**Production-grade, zero-dependency Go SDK for Asiacell Iraq APIs**

[![Go Version](https://img.shields.io/badge/Go-1.26+-18181b?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Verified Endpoints](https://img.shields.io/badge/Endpoints-39%20Verified-E7242A?style=flat-square)](ENDPOINTS.md)
[![Dependencies](https://img.shields.io/badge/Dependencies-Zero%20(Stdlib)-18181b?style=flat-square)](https://pkg.go.dev/)
[![License](https://img.shields.io/badge/License-MIT-18181b?style=flat-square)](LICENSE)

<br />

مكتبة برمجية متكاملة، احترافية، وعالية الأداء مكتوبة بلغة **Go (Golang)** للتعامل مع واجهات برمجة تطبيقات شركة **آسياسيل (Asiacell)** العراقية.  
المكتبة مبنية بنسبة **100% بالاعتماد على المكتبة القياسية للغة Go** وبدون أي مكاتب أو اعتمادات خارجية نهائياً.

</div>

---

## بنية المشروع (Project Structure)

```
asiacell-go/
├── pkg/
│   └── asiacell/
│       ├── client.go               # إدارة العميل، الجلسة، الكوكيز، تدوير هوية الجهاز وحل الكابتشا
│       ├── auth.go                 # المصادقة، طلب وتأكيد رموز التحقق SMS OTP، وتجديد التوكن
│       ├── profile.go              # الاستعلام عن الرصيد، الباقات النشطة، وبيانات الحساب
│       ├── services.go             # شحن الكروت، باقات 4G، تحويل الرصيد، المحافظات والفروع، والخدمات
│       ├── types.go                # هياكل ونماذج البيانات لكافة الطلبات والاستجابات
│       ├── errors.go               # تعريف الأخطاء وتصنيف استجابات السيرفر
│       └── client_test.go          # اختبارات الوحدة والتحقق من سلامة الواجهات
├── examples/
│   ├── 01_auth_and_session/        # تسجيل الدخول وحفظ واستعادة الجلسة وتجاوز الكابتشا
│   ├── 02_account_and_profile/     # فحص الرصيد والباقات المتبقية وتفاصيل المستخدم
│   ├── 03_bundles_and_4g/          # تصفح باقات 4G المفتوحة والعروض الخاصة والاشتراك
│   ├── 04_credit_transfer/         # تحويل الرصيد وتأكيده بالرمز السري
│   ├── 05_recharge_voucher/        # شحن كروت الرصيد والإنترنت وسجل التعبئة
│   ├── 06_services_management/     # حماية رصيد الإنترنت وإلغاء الاشتراكات والخدمات الرقمية
│   ├── 07_shops_and_governorates/  # تصفح محافظات العراق وفروع ومراكز آسياسيل المعتمدة
│   ├── 08_spin_wheel_and_rewards/  # عجلة الحظ اليومية ورصيد الطوارئ (شكراً) وباقات التجوال
│   ├── 09_cdr_incoming_transfer_verification/ # التحقق التلقائي من استلام تحويلات الرصيد عبر سجل CDR
│   ├── 10_advanced_features/       # سوق الأرقام المميزة، إهداء الباقات، تذاكر الدعم، والتعويضات
│   ├── interactive_cli/            # تطبيق تفاعلي متكامل عبر موجه الأوامر (Terminal CLI)
│   └── README.md                   # دليل تشغيل واستخدام كافة الأمثلة الجاهزة
└── go.mod
```

---

## ميزات وقدرات المكتبة (Features)

المكتبة تغطي 39 واجهة برمجية (Endpoints) رسمية تم فحصها والتحقق من سلامتها بنسبة 100%:

### 1. المصادقة وإدارة الجلسات (Authentication & Sessions)
- **طلب رمز التحقق (SMS OTP)**: `Login(ctx, phone)` لإرسال رمز الدخول مباشرة للهاتف.
- **تأكيد الرمز**: `VerifySMS(ctx, pid, passcode)` واستخراج مفاتيح الوصول (`AccessToken` و `RefreshToken`).
- **التجديد التلقائي للتوكن**: `AutoRefreshToken(ctx)` لفحص صلاحية الجلسة وتجديدها في الخلفية دون انقطاع.
- **تخزين الجلسات واستعادتها**: `SaveSessionToFile(path)` و `LoadSessionFromFile(path)` لحفظ بيانات الدخول كملف JSON واستعادتها لاحقاً دون الحاجة لإعادة طلب رمز SMS.
- **تصدير واستيراد الكائنات**: `ExportSession()` و `ImportSession(session)` للتعامل البرمجي المباشر.
- **التعامل مع الكابتشا**: كشف إجباري للكابتشا `ErrCaptchaRequired`، وتدوير هوية الجهاز `RotateDeviceID()` لتجاوزها تلقائياً، مع دعم محلل ذكي `SolveCaptchaOCR()`.

### 2. الحساب والرصيد والباقات (Account, Balance & Addons)
- **ملخص الرصيد والباقات**: `GetProfile(ctx)` لجلب الرصيد الأساسي، باقات الإنترنت المتبقية، دقائق الاتصال، والرسائل.
- **معلومات الحساب الرسمية**: `GetProfileDetails(ctx)` لجلب اسم صاحب الخط، البريد، الصورة، وتاريخ الميلاد.
- **باقات 4G غير المحدودة**: `GetUnlimited4GBundles(ctx)` لاستعراض باقات الإنترنت المفتوحة اليومية والأسبوعية والشهرية مع تفاصيل الأسعار وفترات الصلاحية.
- **تصنيفات الباقات الإضافية**: `GetAddonTags(ctx)` و `GetAddonSummary(ctx, tagID)`.
- **الاشتراك المباشر بالباقات**: `SubscribeAddon(ctx, addonID)` لتفعيل أي باقة فورياً من الرصيد.
- **العروض الحصرية**: `GetSpecialOffers(ctx)` لجلب العروض المخصصة لخط المستخدم.

### 3. كشف الحساب والتحقق التلقائي من التحويلات (CDR & Transfer Verification)
- **كشف حساب تحويلات الرصيد**: `GetCDRTransferHistory(ctx, page, limit)` لجلب سجل العمليات الواردة والصادرة مع رقم المرسل والمبلغ والتاريخ الدقيق.
- **تفعيل كشف الحساب**: `SendCDROTP(ctx)` و `ConfirmCDROTP(ctx, otp)`.
- **التحقق التلقائي المباشر**: `VerifyIncomingTransfer(ctx, senderPhone, minAmount)` للتحقق البرمجي التلقائي والفوري من استلام حوالة رصيد من زبون بدون أي تدخل يدوي للأدمن.
- **تحويل الرصيد (P2P)**: `StartCreditTransfer(ctx, to, amount)` و `ConfirmCreditTransfer(ctx, pid, code)`.
- **سجلات العمليات السابقة**: `GetTransferHistory(ctx)` و `GetRechargeHistory(ctx)`.

### 4. سوق الأرقام المميزة وإهداء الباقات (Vanity Numbers & Gifting)
- **فئات الأرقام المميزة**: `GetVanityClasses(ctx)` لتصفح تصنيفات VIP (الماسية، الذهبية، الفضية).
- **البحث في الأرقام المعروضة للبيع**: `SearchVanityNumbers(ctx, pattern, classId, page, limit)` للبحث عن أرقام بنمط أو فئة معينة مع أسعارها.
- **تفاصيل الرقم المميز**: `GetVanityDetail(ctx, msisdn)`.
- **حجز الرقم المميز**: `ReserveVanityNumber(ctx, msisdn, classId)`.
- **إهداء الباقات للغير**: `SendGiftAddon(ctx, addonId, receiverPhone)` لشراء باقة إنترنت أو مكالمات وإرسالها كهدية لأي رقم مع الخصم من الرصيد.

### 5. الدعم الفني والتعويضات والخدمات الرقمية (Tickets, Compensation & E-Cards)
- **نظام التذاكر والشكاوى**: `GetTicketCategories(ctx)` و `GetTickets(ctx)` و `CreateTicket(ctx, categoryId, description)` لفتح ومتابعة تذاكر الدعم الفني لدى آسياسيل.
- **فحص التعويضات**: `CheckCompensation(ctx)` للاستعلام عن التعويضات الرسمية المتاحة للخط واستلامها.
- **إحالات خطوط Yooz**: `GetYoozMGM(ctx)` و `ApplyYoozMGMCode(ctx, code)`.
- **كروت الألعاب والشحن الرقمي**: `GetEVoucherPackages(ctx)` لتصفح بطاقات PUBG و iTunes و PlayStation.

### 6. شحن كروت التعبئة وحماية الرصيد (Recharge & Protection)
- **شحن الكارت**: `RechargeVoucher(ctx, phone, voucher, rechargeType)` لشحن الرصيد برمز كارت التعبئة المكون من 14 رقماً.
- **التحكم في حماية الرصيد**: تفعيل أو إلغاء ميزة حماية الرصيد لمنع استهلاك الرصيد الأساسي بعد نفاد باقة الإنترنت (`*223#`).
- **إرشادات إلغاء الخدمات الإعلانية**: إرشادات إلغاء الخدمات والمحتوى المزعج (299، 4151، 300).
- **الخدمات الرقمية**: `GetDigitalServices(ctx)`.

### 7. المحافظات والفروع وعجلة الحظ (Shops & Spin Wheel)
- **قائمة المدن والمحافظات**: `GetCities(ctx)`.
- **مراكز ووكلاء آسياسيل**: `GetCityShops(ctx, cityName)` مع العناوين وساعات العمل والإحداثيات (GPS).
- **عجلة الحظ اليومية**: `GetSpinWheelStatus(ctx)` و `PlaySpinWheel(ctx)`.
- **رصيد الطوارئ ونقاط شكراً**: `GetShukranSummary(ctx)`.

---

## البدء السريع (Quick Start)

### 1. تثبيت الحزمة

```bash
go get github.com/FLEX-GHOST/asiacell-go/pkg/asiacell
```

### 2. مثال تسجيل الدخول وحفظ الجلسة

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func main() {
	client, err := asiacell.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// طلب إرسال رمز التحقق SMS
	pid, err := client.Login(ctx, "07701234567")
	if err != nil {
		log.Fatalf("login error: %v", err)
	}
	fmt.Printf("SMS sent successfully! PID: %s\n", pid)

	// تأكيد الرمز المكون من 6 أرقام
	resp, err := client.VerifySMS(ctx, pid, "123456")
	if err != nil {
		log.Fatalf("verify error: %v", err)
	}
	fmt.Printf("Logged in successfully! User ID: %s\n", resp.UserID)

	// حفظ الجلسة للاستخدام المستقبلي
	if err := client.SaveSessionToFile("session.json"); err != nil {
		log.Fatalf("failed to save session: %v", err)
	}
	fmt.Println("Session saved to session.json!")
}
```

### 3. فحص الرصيد والباقات بجلسة سابقة

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func main() {
	client, err := asiacell.NewClient()
	if err != nil {
		log.Fatalf("client error: %v", err)
	}

	if err := client.LoadSessionFromFile("session.json"); err != nil {
		log.Fatalf("session load error: %v", err)
	}

	ctx := context.Background()
	profile, err := client.GetProfile(ctx)
	if err != nil {
		log.Fatalf("profile error: %v", err)
	}

	fmt.Printf("الرصيد الحالي: %s د.ع\n", profile.Balance)
	fmt.Printf("الإنترنت المتبقي: %s\n", profile.RemainingData)
	fmt.Printf("الدقائق المتبقية: %s\n", profile.RemainingVoice)
}
```

---

## دليل الأمثلة الجاهزة (Examples Suite)

| المجلد | الوصف | أمر التشغيل المباشر |
| :--- | :--- | :--- |
| **`01_auth_and_session`** | دورة تسجيل الدخول الكاملة عبر SMS، تدوير هوية الجهاز لتجاوز الكابتشا، وحفظ واستيراد الجلسات. | `go run examples/01_auth_and_session/main.go` |
| **`02_account_and_profile`** | الاستعلام عن الرصيد، الباقات النشطة، الدقائق، الرسائل، وتفاصيل الحساب الشخصي. | `go run examples/02_account_and_profile/main.go` |
| **`03_bundles_and_4g`** | استعراض باقات 4G المفتوحة والعروض الخاصة وتصنيفات الاشتراكات وتفعيلها. | `go run examples/03_bundles_and_4g/main.go` |
| **`04_credit_transfer`** | تحويل الرصيد وتأكيده بالـ OTP، والتحويل المباشر للمحفظة. | `go run examples/04_credit_transfer/main.go` |
| **`05_recharge_voucher`** | شحن الرصيد بكروت التعبئة المكونة من 14 رقماً، واستعراض سجل عمليات الشحن. | `go run examples/05_recharge_voucher/main.go` |
| **`06_services_management`** | إدارة حماية رصيد الإنترنت (*223#)، إرشادات إلغاء الخدمات الإعلانية (299/4151/300). | `go run examples/06_services_management/main.go` |
| **`07_shops_and_governorates`** | استعراض المحافظات العراقية، وفروع ومراكز مبيعات آسياسيل المعتمدة ومواعيد عملها. | `go run examples/07_shops_and_governorates/main.go` |
| **`08_spin_wheel_and_rewards`** | التحقق من تدوير عجلة الحظ اليومية، رصيد الطوارئ ونقاط شكراً، وباقات التجوال الدولي. | `go run examples/08_spin_wheel_and_rewards/main.go` |
| **`09_cdr_incoming_transfer_verification`** | التحقق التلقائي من تحويلات الرصيد عبر سجل كشف الحساب CDR بدون أدمن. | `go run examples/09_cdr_incoming_transfer_verification/main.go` |
| **`10_advanced_features`** | سوق الأرقام المميزة، إهداء الباقات للغير، تذاكر الدعم الفني، التعويضات، وكروت الألعاب. | `go run examples/10_advanced_features/main.go` |
| **`interactive_cli`** | تطبيق تيرمينال تفاعلي شامل يتيح تجربة جميع ميزات المكتبة عبر قائمة نصية مرئية. | `go run examples/interactive_cli/main.go` |

---

## الاختبارات وضمان الجودة (Testing & Quality)

```bash
# تشغيل جميع اختبارات الحزمة
go test -v ./...

# التحقق من فحص مفسر Go القياسي
go vet ./...
```

---

## الترخيص وإخلاء المسؤولية (License & Legal Disclaimer)

- **الترخيص**: هذا المشروع مرخص ومفتوح المصدر بموجب رخصة [MIT](LICENSE).
- **العلامة التجارية**: اسم "Asiacell" وشعارها علامتان تجاريتان مسجلتان لشركة آسياسيل للاتصالات (Asiacell Telecom PJSC).
- **إخلاء المسؤولية**: هذا المشروع (`asiacell-go`) هو مكتبة برمجية مستقلة غير رسمية تم تطويرها لأغراض تعليمية وتطويرية، وليست تابعة لشركة آسياسيل أو معتمدة منها بشكل رسمي.
