<div align="center">

<img src="assets/asiacell.svg" alt="Asiacell Logo" width="240" />

# asiacell-go

**Production-grade, zero-dependency Go SDK for Asiacell Iraq APIs**

[![Go Version](https://img.shields.io/badge/Go-1.26+-18181b?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Verified Endpoints](https://img.shields.io/badge/Endpoints-78%20Verified-E7242A?style=flat-square)](ENDPOINTS.md)
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
│   └── asiacell/                # حزمة الـ SDK الأساسية (بدون اعتمادات خارجية)
│       ├── client.go            # إعداد العميل، الترويسات، وإدارة الجلسات
│       ├── auth.go              # المصادقة، OTP، والكابتشا
│       ├── profile.go           # الملف الشخصي، الرصيد، والباقات الفعالة
│       ├── services.go          # كافة خدمات آسياسيل الـ 78 الرسمية
│       └── types.go             # هياكل البيانات ونماذج الاستجابة JSON
├── examples/                    # أمثلة عملية مستقلة لكل خدمة
│   ├── 01_auth_and_session/     # تسجيل الدخول وحفظ الجلسات
│   ├── 02_account_and_profile/  # فحص الرصيد والباقات والملف الشخصي
│   ├── 03_bundles_and_4g/       # تصفح باقات 4G والعروض
│   ├── 04_credit_transfer/      # تحويل الرصيد
│   ├── 05_recharge_voucher/     # شحن كروت التعبئة
│   ├── 06_services_management/  # إدارة الاشتراكات وUSSD السحابي وحماية الرصيد
│   ├── 07_shops_and_governorates/ # الفروع والمدن
│   ├── 08_spin_wheel_and_rewards/ # عجلة الحظ
│   ├── 09_cdr_incoming_transfer_verification/ # التحقق التلقائي من الحوالات
│   └── 10_advanced_features/    # الأرقام المميزة، التعويضات، والشكاوى
├── ENDPOINTS.md                 # التوثيق التقني لجميع نقاط النهاية الـ 78
├── README.md
└── go.mod
```

---

## ميزات وقدرات المكتبة (Features)

المكتبة تغطي 78 واجهة برمجية (Endpoints) رسمية تم فحصها والتحقق من سلامتها بنسبة 100%:

### 1. المصادقة وإدارة الجلسات (Authentication & Sessions)
- **طلب رمز التحقق (SMS OTP)**: `Login(ctx, phone)` لإرسال رمز الدخول مباشرة للهاتف.
- **تأكيد الرمز**: `VerifySMS(ctx, pid, passcode)` واستخراج مفاتيح الوصول (`AccessToken` و `RefreshToken`).
- **التجديد التلقائي للتوكن**: `AutoRefreshToken(ctx)` لفحص صلاحية الجلسة وتجديدها في الخلفية دون انقطاع.
- **تسجيل الخروج السحابي**: `Logout(ctx)` لإبطال التوكن رسمياً على خوادم آسياسيل.
- **تخزين الجلسات واستعادتها**: `SaveSessionToFile(path)` و `LoadSessionFromFile(path)` لحفظ بيانات الدخول كملف JSON واستعادتها لاحقاً دون الحاجة لإعادة طلب رمز SMS.
- **تصدير واستيراد الكائنات**: `ExportSession()` و `ImportSession(session)` للتعامل البرمجي المباشر.
- **التعامل مع الكابتشا**: كشف إجباري للكابتشا `ErrCaptchaRequired`، وتدوير هوية الجهاز `RotateDeviceID()` لتجاوزها تلقائياً، مع دعم محلل ذكي `SolveCaptchaOCR()`.

### 2. الحساب والرصيد والباقات (Account, Balance & Addons)
- **ملخص الرصيد والباقات**: `GetProfile(ctx)` لجلب الرصيد الأساسي، باقات الإنترنت المتبقية، دقائق الاتصال، والرسائل.
- **معلومات الحساب الرسمية**: `GetProfileDetails(ctx)` لجلب اسم صاحب الخط، البريد، الصورة، وتاريخ الميلاد.
- **تفاصيل الباقة النشطة والتحكم**: `GetBundleDetail(ctx, bundleKey)` لجلب تفاصيل استهلاك الباقة بدقة وأزرار الإجراءات (`actionButtons`).
- **باقات 4G غير المحدودة**: `GetUnlimited4GBundles(ctx)` لاستعراض باقات الإنترنت المفتوحة اليومية والأسبوعية والشهرية مع تفاصيل الأسعار وفترات الصلاحية.
- **تصنيفات الباقات الإضافية**: `GetAddonTags(ctx)` و `GetAddonSummary(ctx, tagID)`.
- **الاشتراك المباشر بالباقات**: `SubscribeAddon(ctx, addonID)` لتفعيل أي باقة فورياً من الرصيد.
- **إلغاء الاشتراك الفوري**: `UnsubscribeAddon(ctx, addonID)` لإلغاء أي باقة أو إيقاف تجديدها التلقائي برمجياً (`actionKey=unsubscribe`).
- **العروض الحصرية**: `GetSpecialOffers(ctx)` و `SubscribeSpecialOffer(ctx, index)` و `CancelSpecialOffer(ctx)` لإدارة باقات عروضي سحابياً.

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
- **تتبع نموذج وتفاصيل التذكرة**: `GetTicketDetail(ctx, ticketNumber)` و `GetTicketForm(ctx, category)`.
- **فحص التعويضات**: `CheckCompensation(ctx)` للاستعلام عن التعويضات الرسمية المتاحة للخط واستلامها.
- **إحالات خطوط Yooz**: `GetYoozMGM(ctx)` و `ApplyYoozMGMCode(ctx, code)`.
- **كروت الألعاب والشحن الرقمي**: `GetEVoucherPackages(ctx)` لتصفح بطاقات PUBG و iTunes و PlayStation.

### 6. الاشتراكات، الحسابات المتعددة، والـ USSD السحابي (Subscriptions & Multi-Account)
- **الاشتراكات والخدمات الفعالة**: `GetMySubscriptions(ctx)` للاستعلام المباشر من خوادم آسياسيل عن جميع الخدمات النشطة وتواريخ صلاحيتها وأزرار الإلغاء.
- **خدمات USSD السحابية التفاعلية**: `GetUSSDMenu(ctx, parentId)` و `SubmitUSSDAction(ctx, params)` لتصفح وتنفيذ خدمات وقوائم الـ USSD عبر السحابة مباشرة.
- **إدارة الحسابات المتعددة (Multi-Account)**:
  - `GetLinkedAccounts(ctx)`: جلب الأرقام والخطوط المرتبطة بالحساب.
  - `AddLinkedAccount(ctx, phone)` و `ConfirmLinkedAccount(ctx, phone, pin)`: ربط خط جديد وتأكيده.
  - `SwitchActiveAccount(ctx, phone)`: التبديل بين الخطوط لإدارتها بنفس الجلسة.
  - `RemoveLinkedAccount(ctx, phone)`: حذف ارتباط الخط.
- **مشاركة البيانات وسقوف الاستهلاك (Data Sharing & Caps)**:
  - `GetDataCapLimit(ctx)` و `SetDataCapLimit(ctx, limitMB)`: فحص وتعيين سقف استهلاك يومي للبيانات.
  - `GetBundleShareLimit(ctx)` و `SetBundleShareLimit(ctx, phone, limitMB)`: فحص وتحديد سقف ميغابايت لكل خط مشارك.
  - `GetManageLines(ctx)`: استعلام الخطوط المشاركة في الباقة.
- **سجل وتاريخ الباقات**: `GetSubscriptionHistory(ctx)`.
- **شحن الكارت**: `RechargeVoucher(ctx, phone, voucher, rechargeType)`.
- **التحكم في حماية الرصيد**: تفعيل أو إلغاء ميزة حماية الرصيد لمنع استهلاك الرصيد الأساسي بعد نفاد باقة الإنترنت (`*223#`).
- **إرشادات إلغاء الخدمات الإعلانية**: إرشادات إلغاء الخدمات والمحتوى المزعج (299، 4151، 300).
- **الخدمات الرقمية**: `GetDigitalServices(ctx)`.

### 7. المحافظات والفروع وعجلة الحظ (Shops & Spin Wheel)
- **قائمة المدن والمحافظات**: `GetCities(ctx)`.
- **مراكز ووكلاء آسياسيل**: `GetCityShops(ctx, cityName)` مع العناوين وساعات العمل والإحداثيات (GPS).
- **عجلة الحظ اليومية**: `GetSpinWheelStatus(ctx)` و `PlaySpinWheel(ctx)`.
- **رصيد الطوارئ ونقاط شكراً**: `GetShukranInfo(ctx)` و `RequestShukranCredit(ctx)`.
- **خدمات التجوال الدولي**: `GetRoamingInfo(ctx)`.

### 8. هز واربح ودفع الفواتير (Shake & Win & Billing)
- **هز واربح (Shake & Win)**: `GetShakeAndWinStatus(ctx)` و `PlayShakeAndWin(ctx, txID)`.
- **هدايا ما بعد الشحن**: `GetTopupShakeAndWin(ctx)`.
- **فواتير الخطوط الآجلة الدفع (Postpaid)**: `GetBillAmount(ctx)` و `PayBill(ctx, phone, amount)`.

### 9. خطوط البيانات والراوترات ومحرك البحث (Data Lines, Routers & Search)
- **خطوط البيانات والراوترات 4G**: `GetDataLineInfo(ctx)` و `PairDataLine(ctx, msisdn, iccid)`.
- **إدارة الاتصالات المتعددة**: `GetMultiLineConnections(ctx)` و `GetMultiLineHome(ctx)`.
- **محرك البحث الشامل**: `Search(ctx, query)` و `GetSearchSuggestions(ctx, query)`.
- **تحديث الملف الشخصي**: `UpdateProfileInfo(ctx, name, email)`.

### 10. برنامج مكافآت وفاء والتعرفة الدولية (Loyalty Rewards & International Tariffs)
- **برنامج مكافآت وفاء (Loyalty Rewards)**: `GetLoyaltyRewards(ctx)` و `GetLoyaltyRewardDetail(ctx)` و `RedeemLoyaltyReward(ctx)` و `CheckLoyaltyRewardStatus(ctx, pid)`.
- **التعرفة والخدمات الدولية**: `GetInternationalTariffs(ctx)` و `GetInternationalServices(ctx)`.

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
