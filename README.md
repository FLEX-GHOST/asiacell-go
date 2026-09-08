<div align="center">

<img src="assets/asiacell.svg" alt="Asiacell Logo" width="240" />

# asiacell-go

**Production-grade, zero-dependency Go SDK for Asiacell Iraq APIs**

[![Go Version](https://img.shields.io/badge/Go-1.26+-18181b?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Verified Endpoints](https://img.shields.io/badge/Endpoints-24%20Verified-E7242A?style=flat-square)](ENDPOINTS.md)
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
│   ├── 04_credit_transfer/         # تحويل الرصيد وتأكيده والتحقق من استلام التحويلات
│   ├── 05_recharge_voucher/        # شحن كروت الرصيد والإنترنت وسجل التعبئة
│   ├── 06_services_management/     # حماية رصيد الإنترنت وإلغاء الاشتراكات والخدمات الرقمية
│   ├── 07_shops_and_governorates/  # تصفح محافظات العراق وفروع ومراكز آسياسيل المعتمدة
│   ├── 08_spin_wheel_and_rewards/  # عجلة الحظ اليومية ورصيد الطوارئ (شكراً) وباقات التجوال
│   ├── interactive_cli/            # تطبيق تفاعلي متكامل عبر موجه الأوامر (Terminal CLI)
│   └── README.md                   # دليل تشغيل واستخدام كافة الأمثلة الجاهزة
└── go.mod
```

---

## ميزات وقدرات المكتبة (Features)

المكتبة تغطي 24 واجهة برمجية (Endpoints) رسمية تم فحصها والتحقق من سلامتها بنسبة 100%:

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

### 3. تحويل الرصيد والتحقق المالي (Credit Transfer & Verification)
- **تحويل الرصيد (P2P)**: `StartCreditTransfer(ctx, receiverPhone, amount)` لبدء التحويل وإرسال كود التأكيد.
- **تأكيد التحويل**: `ConfirmCreditTransfer(ctx, pid, passcode)` لتنفيذ التحويل بالرمز السري.
- **تحويل مباشر للمحفظة**: `TransferToWallet(ctx, amount)` لدعم التحويل إلى المحفظة المعتمدة.
- **التحقق التلقائي من استلام التحويلات**: `VerifyIncomingTransfer(ctx, targetPhone, minAmount)` للتحقق البرمجي التلقائي من وصول الرصيد إلى رقمك/محفظتك مع منع التكرار (De-duplication) لتأكيد المعاملات المالية والمبيعات تلقائياً.
- **سجلات العمليات**: `GetTransferHistory(ctx)` و `GetRechargeHistory(ctx)` و `GetSubscriptionHistory(ctx)`.

### 4. شحن كروت التعبئة (Voucher Recharge)
- **شحن الكارت**: `RechargeVoucher(ctx, phone, voucher, rechargeType)` لشحن الرصيد برمز كارت التعبئة المكون من 14 رقماً (كروت الرصيد العادية وكروت الإنترنت).

### 5. حماية الرصيد وإدارة الخدمات (Balance Protection & Services)
- **استعراض الخدمات المفعلة**: `GetServicesManagement(ctx)`.
- **التحكم في حماية الرصيد**: تفعيل أو إلغاء ميزة حماية الرصيد لمنع استهلاك الرصيد الأساسي بعد نفاد باقة الإنترنت (المكافئة لـ `*223#`).
- **إرشادات إلغاء الخدمات الإعلانية**: توفير تعليمات إلغاء الخدمات الإعلانية والمحتوى المزعج عبر الرسائل القصيرة (299، 4151، 300).
- **الخدمات الرقمية**: `GetDigitalServices(ctx)` لاستعراض الخدمات الترفيهية والمحتوى.

### 6. المحافظات والفروع المعتمدة (Governorates & Certified Shops)
- **قائمة المدن والمحافظات**: `GetCities(ctx)` لجلب جميع محافظات العراق (بغداد، أربيل، البصرة، النجف، كربلاء، السليمانية، كركوك...).
- **مراكز ووكلاء آسياسيل**: `GetCityShops(ctx, cityName)` لجلب العناوين الدقيقة، ساعات العمل، حالة الفرع (مفتوح حالياً أو مغلق)، وأرقام الهواتف والإحداثيات الجغرافية (GPS).

### 7. عجلة الحظ والمكافآت (Spin Wheel & Shukran Loyalty)
- **عجلة الحظ اليومية**: `GetSpinWheelStatus(ctx)` لفحص إمكانية التدوير، و `PlaySpinWheel(ctx)` للمطالبة بالجائزة اليومية (إنترنت أو رصيد مجاني).
- **رصيد الطوارئ ونقاط شكراً**: `GetShukranSummary(ctx)` للاستعلام عن نقاط المكافآت وخدمات رصيد الطوارئ.
- **باقات التجوال الدولي**: `GetRoamingBundles(ctx)` لاستعراض باقات السفر والتجوال.

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

### 4. استعراض باقات 4G المفتوحة والاشتراك

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/FLEX-GHOST/asiacell-go/pkg/asiacell"
)

func main() {
	client, _ := asiacell.NewClient()
	client.LoadSessionFromFile("session.json")
	ctx := context.Background()

	bundles, err := client.GetUnlimited4GBundles(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, b := range bundles {
		fmt.Printf("الباقة: %s | السعر: %s د.ع | الكود: %s\n", b.Title, b.Price, b.ID)
	}
}
```

---

## دليل الأمثلة الجاهزة (Examples Suite)

يتضمن المشروع مجلداً غنياً بالأمثلة المستقلة في `examples/`، تم إعداد كل مثال ليعمل بشكل فوري:

| المجلد | الوصف | أمر التشغيل المباشر |
| :--- | :--- | :--- |
| **`01_auth_and_session`** | دورة تسجيل الدخول الكاملة عبر SMS، تدوير هوية الجهاز لتجاوز الكابتشا، وحفظ واستيراد الجلسات. | `go run examples/01_auth_and_session/main.go` |
| **`02_account_and_profile`** | الاستعلام عن الرصيد، الباقات النشطة، الدقائق، الرسائل، وتفاصيل الحساب الشخصي. | `go run examples/02_account_and_profile/main.go` |
| **`03_bundles_and_4g`** | استعراض باقات 4G المفتوحة والعروض الخاصة وتصنيفات الاشتراكات وتفعيلها. | `go run examples/03_bundles_and_4g/main.go` |
| **`04_credit_transfer`** | تحويل الرصيد وتأكيده بالـ OTP، والتحقق التلقائي الذكي من وصول الرصيد للمحفظة. | `go run examples/04_credit_transfer/main.go` |
| **`05_recharge_voucher`** | شحن الرصيد بكروت التعبئة المكونة من 14 رقماً، واستعراض سجل عمليات الشحن. | `go run examples/05_recharge_voucher/main.go` |
| **`06_services_management`** | إدارة حماية رصيد الإنترنت (*223#)، إرشادات إلغاء الخدمات الإعلانية (299/4151/300). | `go run examples/06_services_management/main.go` |
| **`07_shops_and_governorates`** | استعراض المحافظات العراقية، وفروع ومراكز مبيعات آسياسيل المعتمدة ومواعيد عملها. | `go run examples/07_shops_and_governorates/main.go` |
| **`08_spin_wheel_and_rewards`** | التحقق من تدوير عجلة الحظ اليومية، رصيد الطوارئ ونقاط شكراً، وباقات التجوال الدولي. | `go run examples/08_spin_wheel_and_rewards/main.go` |
| **`interactive_cli`** | تطبيق تيرمينال تفاعلي شامل يتيح تجربة جميع ميزات المكتبة عبر قائمة نصية مرئية. | `go run examples/interactive_cli/main.go` |

لتشغيل التطبيق التفاعلي الشامل:
```bash
go run examples/interactive_cli/main.go
```

---

## الاختبارات وضمان الجودة (Testing & Quality)

المكتبة تخضع لاختبارات صارمة تضمن استقرار ونظافة الكود، وعدم وجود أي تسريب للذاكرة أو الـ Goroutines:

```bash
# تشغيل جميع اختبارات الحزمة
go test -v ./...

# التحقق من فحص مفسر Go القياسي
go vet ./...
```

---

## الأمان ومعايير الذاكرة (Security & Memory Policy)

- **صفر اعتمادات خارجية**: الاعتماد الكامل على مكتبة Go القياسية لحماية تامة من ثغرات سلاسل التوريد.
- **تشفير الاتصال**: جميع الاتصالات مشفرة بـ HTTPS بالكامل ومطابقة للمواصفات الرسمية لشبكة آسياسيل.
- **إدارة الذاكرة والاتصالات**: غلق قاطع لجميع اتصالات الشبكة ومقابض الملفات عبر `defer` لمنع أي تسريب للموارد.

---

## الترخيص وإخلاء المسؤولية (License & Legal Disclaimer)

- **الترخيص**: هذا المشروع مرخص ومفتوح المصدر بموجب رخصة [MIT](LICENSE).
- **العلامة التجارية**: اسم "Asiacell" وشعارها علامتان تجاريتان مسجلتان لشركة آسياسيل للاتصالات (Asiacell Telecom PJSC).
- **إخلاء المسؤولية**: هذا المشروع (`asiacell-go`) هو مكتبة برمجية مستقلة غير رسمية تم تطويرها لأغراض تعليمية وتطويرية، وليست تابعة لشركة آسياسيل أو معتمدة منها بشكل رسمي.

