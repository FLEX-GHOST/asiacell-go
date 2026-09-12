<div align="center">

<img src="assets/asiacell.svg" alt="شعار آسياسيل" width="200" />

# مرجع نقاط نهاية واجهة برمجة تطبيقات آسياسيل (Asiacell API Endpoints)

**المواصفات الكاملة لجميع نقاط نهاية HTTP الرسمية لتطبيق آسياسيل (Asiacell API) باللغة العربية مع نماذج الطلب والاستجابة.**

<br />

[![Specification](https://img.shields.io/badge/Specification-100%25%20Verified-E7242A?style=flat-square)](ENDPOINTS.md)
[![Endpoints](https://img.shields.io/badge/Endpoints-147%20Verified-18181b?style=flat-square)](ENDPOINTS.md)
[![Protocol](https://img.shields.io/badge/Protocol-HTTPS%2FREST-18181b?style=flat-square)](ENDPOINTS.md)
[![Client](https://img.shields.io/badge/Go%20Client-asiacell--go-18181b?style=flat-square)](https://github.com/FLEX-GHOST/asiacell-go)

<br />

جميع نقاط النهاية الموضحة أدناه تم فحصها والتحقق منها مباشرة مقابل خوادم آسياسيل الرسمية (`odpapp.asiacell.com` و `app.asiacell.com`) مع اجتياز اختبارات الوحدة والاختبارات الحية ومعالجة الأخطاء.

</div>

---

## 1. جدول نقاط النهاية المعتمدة (Endpoints Matrix)

| # | طريقة HTTP | مسار النقطة (Endpoint Path) | دالة Go SDK | الوصف التفصيلي |
| :---: | :---: | :--- | :--- | :--- |
| **01** | `POST` | `/api/v1/login` (`/api/v1/auth/phone`) | `client.Login(ctx, phone)` | إرسال رمز التحقق OTP المكون من 6 أرقام برسالة SMS وتوليد معرّف العملية PID لتسجيل الدخول |
| **02** | `POST` | `/api/v1/smsvalidation` (`/api/v1/auth/login-passcode`) | `client.VerifySMS(ctx, pid, code)` | التحقق من كود الـ SMS واستبداله بتوكنات الجلسة وتفعيل التسجيل البيومتري للجهاز تلقائياً |
| **03** | `POST` | `/api/v1/validate` (`/api/v1/auth/refresh-token`) | `client.RefreshToken(ctx)` / `client.RefreshSession(ctx)` | تجديد توكن الوصول المنتهي تلقائياً واسترجاع السر البيومتري المشفر في الخلفية بدون طلب كود مجدداً |
| **04** | `GET` | `/api/v1/captcha` (`/api/v1/auth/captcha`) | `client.SolveCaptchaOCR(data)` / `client.GetCaptcha(ctx)` | جلب صورة اختبار الكابتشا وحلها آلياً في حال فرضها لمنع الطلبات المتكررة |
| **05** | `GET` | `/api/v1/profile` | `client.GetProfile(ctx)` | استعلام رصيد الحساب الحالي، مدة الصلاحية، المتبقي من الإنترنت والمكالمات والرسائل |
| **06** | `GET` | `/api/v1/profile/view2` | `client.GetProfileDetails(ctx)` | تفاصيل الملف الشخصي للمشترك: الاسم الكامل، البريد، تاريخ الميلاد، ورابط الصورة |
| **07** | `GET` | `/api/v1/profile-img` | `client.GetProfileDetails(ctx)` | دفق صورة الحساب الرمزية للمستخدم (Avatar) |
| **08** | `GET` | `/api/v2/4gcontent` | `client.GetUnlimited4GBundles(ctx)` | باقات الإنترنت المفتوح 4G (يومية، أسبوعية، شهرية) وأسعارها |
| **09** | `GET` | `/api/v1/products/` | `client.GetSpecialOffers(ctx)` | العروض الخاصة والتخفيضات المخصصة لرقم الهاتف |
| **10** | `GET` | `/api/v1/addon/tags/` | `client.GetAddonTags(ctx)` | أقسام وتصنيفات الباقات الإضافية (إنترنت، تجوال، اتصالات، خدمات) |
| **11** | `GET` | `/api/v1/addon/summary/` | `client.GetAddonSummary(ctx, tagID)` | تفاصيل الباقات والأسعار لقسم محدد من الباقات |
| **12** | `POST` | `/api/v1/addon/subscribe` | `client.SubscribeAddon(ctx, addonID)` | تفعيل واشتراك فوري في الباقة مع خصم قيمتها من الرصيد مباشرة |
| **13** | `GET` | `/api/v1/cdr/detail?type=btransfer` | `client.GetCDRTransferHistory(ctx, page, limit)` | جلب سجل كشف الحساب لتحويلات الرصيد الواردة والصادرة مع رقم المرسل والمبلغ والتاريخ |
| **14** | `POST` | `/api/v1/cdr/send-otp` | `client.SendCDROTP(ctx)` | طلب رمز OTP واستخراج معرّف العملية PID من رابط nextUrl لتفعيل خدمة كشف الحساب للجلسة |
| **15** | `POST` | `/api/v1/cdr/confirm` | `client.ConfirmCDROTP(ctx, pid, code)` | تأكيد رمز OTP بواسطة GenericSMSConfirmationDTO لتفعيل كشف الحساب (حظر تام لـ smsvalidation) |
| **16** | `GET` | `/api/v1/transaction/transfer` | `client.GetTransferHistory(ctx)` | سجل وتاريخ عمليات تحويل الرصيد السابقة على الخط (المحفظة) |
| **17** | `POST` | `/api/v1/credit-transfer/start` | `client.StartCreditTransfer(ctx, to, amt)` | بدء تحويل رصيد من الشريحة لرقم آخر وتوليد معرف العملية PID |
| **18** | `POST` | `/api/v1/credit-transfer/do-transfer` | `client.ConfirmCreditTransfer(ctx, pid, code)` | تأكيد تحويل الرصيد بإدخال رمز التحقق المرسل إلى الهاتف |
| **19** | `POST` | `/api/v1/top-up` | `client.RechargeVoucher(ctx, to, code, type)` | شحن وتعبئة كارت آسياسيل عبر كود الكارت (13-14 رقماً) عادي أو إنترنت |
| **20** | `GET` | `/api/v1/transaction/recharge` | `client.GetRechargeHistory(ctx)` | سجل وتاريخ عمليات شحن الكروت السابقة على الخط |
| **21** | `GET` | `/api/v1/profile/subscriptions` | `client.GetMySubscriptions(ctx)` | استعلام جميع الخدمات والاشتراكات الفعالة على الخط وصلاحياتها مباشرة من الخادم |
| **22** | `POST` | `/api/v1/ussd` | `client.SubmitUSSDAction(ctx, params)` | تنفيذ وتفعيل أوامر وخدمات الـ USSD السحابية التفاعلية عبر HTTP مباشرة |
| **23** | `GET` | `/api/v1/ussd` | `client.GetUSSDMenu(ctx, parentID)` | تصفح واستعلام قوائم الـ USSD التفاعلية السحابية من الخادم |
| **24** | `GET` | `/api/v1/transaction/bundle` | `client.GetSubscriptionHistory(ctx)` | سجل وتاريخ عمليات شراء واشتراك الباقات السابقة على الخط |
| **25** | `GET` | `/api/v1/digital-services` | `client.GetDigitalServices(ctx)` | الخدمات الترفيهية والقيمة المضافة الرقمية |
| **26** | `GET` | `/api/v1/partners/cities` | `client.GetCities(ctx)` | قائمة بجميع المحافظات والمدن العراقية ومعرفاتها |
| **27** | `GET` | `/api/v1/shops` | `client.GetCityShops(ctx, cityID)` | مراكز وفروع آسياسيل المعتمدة: العناوين، الإحداثيات، وحالة الفتح/الإغلاق |
| **28** | `GET` | `/api/v1/spin-wheel` | `client.GetSpinWheelStatus(ctx)` | حالة عجلة الحظ اليومية وما إذا كانت متاحة للدوران |
| **29** | `POST` | `/api/v1/spin-wheel/play` | `client.PlaySpinWheel(ctx)` | تدوير عجلة الحظ واستلام الجائزة اليومية (ميغابايت أو رصيد مجاني) |
| **30** | `GET` | `/api/v2/vanity/classes` | `client.GetVanityClasses(ctx)` | استعراض فئات وتصنيفات الأرقام المميزة (الماسية، الذهبية، الفضية) وأسعارها |
| **31** | `GET` | `/api/v2/vanity` | `client.SearchVanityNumbers(ctx, pat, class, p, l)` | البحث في الأرقام المميزة المتاحة للبيع بنمط أو فئة معينة |
| **32** | `GET` | `/api/v2/vanity/{msisdn}/detail` | `client.GetVanityDetail(ctx, msisdn)` | تفاصيل وسعر وشروط حجز رقم مميز محدد |
| **33** | `POST` | `/api/v2/vanity` | `client.ReserveVanityNumber(ctx, msisdn, classId)` | حجز الرقم المميز مباشرة باسم المشترك وتوليد معرف العملية |
| **34** | `POST` | `/api/v1/addon/send-as-gift` | `client.SendGiftAddon(ctx, addonId, to)` | شراء باقة إنترنت أو اتصالات وإهداؤها لرقم آخر بخصم من الرصيد |
| **35** | `GET` | `/api/v1/resolution-center/categories` | `client.GetTicketCategories(ctx)` | أقسام وتصنيفات الشكاوى الفنية المعتمدة في آسياسيل |
| **36** | `GET` | `/api/v1/resolution-center` | `client.GetTickets(ctx)` | سجل تذاكر الشكاوى المفتوحة وتحديثات مسار المعالجة |
| **37** | `POST` | `/api/v1/resolution-center` | `client.CreateTicket(ctx, catId, desc)` | فتح وإرسال تذكرة شكوى رسمية جديدة لإدارة العمليات والدعم |
| **38** | `GET` | `/api/v1/compensation` | `client.CheckCompensation(ctx)` | فحص استحقاق الخط للتعويضات الرسمية (جيجابايت أو رصيد مجاني) |
| **39** | `GET` | `/api/v1/yooz-mgm` | `client.GetYoozMGM(ctx)` | استخراج كود الإحالة وإحصائيات دعوة الأصدقاء لخطوط Yooz |
| **40** | `POST` | `/api/v1/yooz-mgm/apply-code` | `client.ApplyYoozMGMCode(ctx, code)` | تفعيل كود دعوة للحصول على البونص والمكافآت المجانية |
| **41** | `GET` | `/api/v2/e-voucher/packages` | `client.GetEVoucherPackages(ctx)` | تصفح كروت الألعاب والشحن الرقمي (PUBG, PlayStation, iTunes) |
| **42** | `GET` | `/api/v1/notifications` | `client.GetNotifications(ctx)` | جلب صندوق الإشعارات والتنبيهات المستلمة من آسياسيل |
| **43** | `GET` | `/api/v1/promotions` | `client.GetPromotions(ctx)` | العروض الترويجية والحملات الإعلانية الحالية |
| **44** | `GET` | `/api/v1/shukran` | `client.GetShukranInfo(ctx)` | فحص خدمة شكراً لسلفة الرصيد أو الإنترنت عند الطوارئ |
| **45** | `GET` | `/api/v1/roaming` | `client.GetRoamingInfo(ctx)` | استعلام خدمات وباقات التجوال الدولي المتاحة للخط |
| **46** | `POST` | `/api/v1/addon` | `client.UnsubscribeAddon(ctx, addonID)` | إلغاء الاشتراك وإيقاف التجديد التلقائي لأي باقة برمجياً (`actionKey=unsubscribe`) |
| **47** | `GET` | `/api/v1/profile/bundle/{key}` | `client.GetBundleDetail(ctx, key)` | تفاصيل الباقة الفعالة وسجل الحصص المستهلكة وأزرار الإجراءات |
| **48** | `POST` | `/api/v1/logout` | `client.Logout(ctx)` | تسجيل الخروج وإبطال التوكن رسمياً على خوادم آسياسيل السحابية |
| **49** | `GET` | `/api/v1/map-account` | `client.GetLinkedAccounts(ctx)` | استعلام كافة الأرقام والخطوط الإضافية المرتبطة بالحساب |
| **50** | `POST` | `/api/v1/map-account` | `client.AddLinkedAccount(ctx, phone)` | إرسال طلب ربط خط إضافي للحساب وتوليد كود OTP للتحقق |
| **51** | `POST` | `/api/v1/map-account/confirm` | `client.ConfirmLinkedAccount(ctx, phone, pin)` | تأكيد ربط الرقم الإضافي وإدراجه ضمن حسابات المستخدم |
| **52** | `POST` | `/api/v1/account-action/switch` | `client.SwitchActiveAccount(ctx, phone)` | التبديل بين الخطوط المرتبطة لإدارتها بشكل فوري |
| **53** | `POST` | `/api/v1/account-action/remove` | `client.RemoveLinkedAccount(ctx, phone)` | فك وحذف ارتباط رقم إضافي من الحساب |
| **54** | `POST` | `/api/v1/addon/datacap/set-limit` | `client.SetDataCapLimit(ctx, limitMB)` | تعيين سقف استهلاك البيانات اليومي للخط بالكامل لمنع استنزاف الرصيد |
| **55** | `POST` | `/api/v1/addon/share/set-limit` | `client.SetBundleShareLimit(ctx, phone, limitMB)` | تحديد سقف ميغابايت لكل رقم مشارك في باقة الإنترنت المشتركة |
| **56** | `GET` | `/api/v1/shake-and-win` | `client.GetShakeAndWinStatus(ctx)` | استعلام حالة لعبة هز واربح اليومية والجوائز المتاحة |
| **57** | `POST` | `/api/v1/shake-and-win` | `client.PlayShakeAndWin(ctx, txID)` | تنفيذ الهزة واستلام الجائزة الفورية (رصيد أو ميغابايت) |
| **58** | `GET` | `/api/v1/top-up/shake-and-win` | `client.GetTopupShakeAndWin(ctx)` | فحص استحقاق جوائز هز واربح بعد تعبئة الرصيد |
| **59** | `GET` | `/api/v1/top-up/bill-amount` | `client.GetBillAmount(ctx)` | استعلام مبلغ الفاتورة المستحقة وتاريخ استحقاقها للخطوط الآجلة الدفع |
| **60** | `POST` | `/api/v1/top-up/pay-bill` | `client.PayBill(ctx, phone, amount)` | تسديد ودفع فاتورة الخط الآجل الدفع |
| **61** | `GET` | `/api/v1/resolution-center/{num}` | `client.GetTicketDetail(ctx, ticketNum)` | متابعة تفاصيل تذكرة شكوى محددة مع مسار المعالجة وردود الدعم |
| **62** | `GET` | `/api/v1/resolution-center/ticket-form` | `client.GetTicketForm(ctx, category)` | استعلام حقول الإدخال والشروط المطلوبة لتقديم شكوى في قسم محدد |
| **63** | `GET` | `/api/v1/data-line` | `client.GetDataLineInfo(ctx)` | استعلام بيانات خط البيانات أو راوتر الإنترنت المرتبط بالحساب |
| **64** | `POST` | `/api/v1/data-line` | `client.PairDataLine(ctx, msisdn, iccid)` | إقران وربط خط راوتر / شريحة بيانات جديدة بالحساب |
| **65** | `GET` | `/api/v1/multi-line` | `client.GetMultiLineConnections(ctx)` | جلب قائمة كافة الخطوط والراوترات المربوطة بالحساب |
| **66** | `GET` | `/api/v1/multi-line/home` | `client.GetMultiLineHome(ctx)` | الشاشة الرئيسية ولوحة التحكم بالخطوط المتعددة والراوترات |
| **67** | `GET` | `/api/v1/search` | `client.Search(ctx, query)` | محرك البحث الشامل في خدمات وباقات وعروض ومراكز آسياسيل |
| **68** | `GET` | `/api/v1/search/suggestions` | `client.GetSearchSuggestions(ctx, query)` | اقتراحات الإكمال التلقائي الفوري أثناء البحث |
| **69** | `GET` | `/api/v1/international-services/tariff` | `client.GetInternationalTariffs(ctx)` | جدول أسعار وتعرفة المكالمات الدولية لكل دولة حول العالم |
| **70** | `GET` | `/api/v1/international-services` | `client.GetInternationalServices(ctx)` | باقات وتخفيضات الاتصال الدولي المتاحة للخط |
| **71** | `GET` | `/api/v1/reward` | `client.GetLoyaltyRewards(ctx)` | قائمة جوائز ومكافآت نقاط برنامج الولاء (وفاء) |
| **72** | `GET` | `/api/v1/reward/detail` | `client.GetLoyaltyRewardDetail(ctx)` | تفاصيل ومواصفات المكافأة المحددة والنقاط المطلوبة |
| **73** | `POST` | `/api/v1/eo` | `client.RedeemLoyaltyReward(ctx)` | إرسال طلب استبدال النقاط بالمكافأة وتوليد رمز الاستلام |
| **74** | `GET` | `/api/v1/eo/check-status` | `client.CheckLoyaltyRewardStatus(ctx, pid)` | فحص حالة كود استلام المكافأة المستبدلة |
| **75** | `POST` | `/api/v3/profile/update` | `client.UpdateProfileInfo(ctx, name, email)` | تعديل وتحديث بيانات المشترك (الاسم الكامل والبريد الإلكتروني) |
| **76** | `GET` | `/api/v1/addon/datacap/limit` | `client.GetDataCapLimit(ctx)` | استعلام سقف البيانات اليومي المفروض حالياً على الخط |
| **77** | `GET` | `/api/v1/addon/share/limit` | `client.GetBundleShareLimit(ctx)` | استعلام سقف الميغابايت المخصص لكل رقم مشارك في الباقة |
| **78** | `GET` | `/api/v1/addon/share` | `client.GetManageLines(ctx)` | جلب قائمة الأرقام والخطوط المشاركة الفعالة في الباقة العائلية |
| **79** | `GET` | `/api/v1/fanzone/home` | `client.GetFanZoneHome(ctx, compId)` | استعلام المنافسات والبطولات النشطة لدوري نجوم العراق |
| **80** | `GET` | `/api/v1/fanzone/kick-and-win/home` | `client.GetFanZoneKickAndWin(ctx, compId)` | استعلام لعبة ركل واربح والمحاولات المتبقية |
| **81** | `POST` | `/api/v1/fanzone/kick-and-win/play-finish` | `client.FinishFanZoneKickAndWin(ctx, compId, score)` | إنهاء لعبة ركل واربح وتسجيل النتيجة المحققة |
| **82** | `GET` | `/api/v1/fanzone/kick-and-win/reward` | `client.GetFanZoneKickAndWinReward(ctx, compId, ticketId)` | استلام وصرف مكافأة وجائزة لعبة ركل واربح |
| **83** | `GET` | `/api/v1/fanzone/leader-board` | `client.GetFanZoneLeaderBoard(ctx, compId)` | استعراض لوحة المتصدرين والنقاط وترتيب المستخدم |
| **84** | `GET` | `/api/v1/fanzone/predict-and-win` | `client.GetFanZonePredictions(ctx, compId)` | جدول المباريات وتوقع النتائج للفوز بجوائز البطولة |
| **85** | `GET` | `/api/v1/fanzone/grand-prizes` | `client.GetFanZoneGrandPrizes(ctx, compId)` | استعراض قائمة الجوائز الكبرى للموسم والمسابقات |
| **86** | `GET` | `/api/v1/fanzone/rewards-history` | `client.GetFanZoneRewardsHistory(ctx, compId)` | سجل وتاريخ الجوائز والمكافآت السابقة التي حصل عليها المستخدم |
| **87** | `POST` | `/api/v1/fanzone/onboarding/gen-nickname` | `client.GenerateFanZoneNickname(ctx, nickname)` | توليد واقتراح اسم مستعار للمشارك في مسابقات FanZone |
| **88** | `GET` | `/api/v1/fanzone/answer-and-win` | `client.GetFanZoneAnswerAndWin(ctx, compId)` | استعلام لعبة أجب واربح الرياضية التفاعلية |
| **89** | `GET` | `/api/v1/fanzone/favourite-team/pick` | `client.PickFanZoneFavoriteTeam(ctx, compId, teamId)` | اختيار وتثبيت الفريق المفضل للمشترك في الدوري |
| **90** | `GET` | `/api/v1/fanzone/champion-team/pick` | `client.PickFanZoneChampionTeam(ctx, compId, teamId)` | اختيار وتوقع الفريق البطل المتوج بالدوري |
| **91** | `GET` | `/api/v1/partners/categories` | `client.GetPartnerCategories(ctx)` | أقسام وتصنيفات الشركاء والمتاجر في برنامج الخصومات (مطاعم، فنادق، مقاهي) |
| **92** | `GET` | `/api/v2/partners/cities` | `client.GetPartnerCities(ctx)` | قائمة المحافظات والمدن المشمولة ببرنامج الخصومات |
| **93** | `GET` | `/api/v2/partners/cities/{cityId}/categories` | `client.GetPartnerCityCategories(ctx, cityId)` | تصنيفات الشركاء والمتاجر المتاحة في محافظة محددة |
| **94** | `GET` | `/api/v2/categories/{catId}/cities/{cityId}/partners` | `client.GetPartnersByCategoryAndCity(ctx, catId, cityId)` | قائمة المتاجر والشركاء مع العنوان والإحداثيات ونسبة الخصم |
| **95** | `POST` | `/api/v2/partners/register` | `client.RegisterPartner(ctx, req)` | تقديم طلب تسجيل متجر أو نشاط تجاري جديد في برنامج شركاء آسياسيل |
| **96** | `GET` | `/api/v5/avocado/home` | `client.GetYoozHome(ctx)` | لوحة تحكم خطوط Yooz الشبابية الرسمية (الرصيد، الصلاحية، الميغابايت) |
| **97** | `GET` | `/api/v3/avocado/bundles/screen` | `client.GetYoozBundlesScreen(ctx, groupId)` | شاشة عرض باقات وعروض خطوط يوز |
| **98** | `GET` | `/api/v3/avocado/bundles/classic-plans` | `client.GetYoozClassicPlans(ctx)` | باقات وخطط يوز كلاسيك الشهرية |
| **99** | `GET` | `/api/v3/avocado/bundles/omega-plans` | `client.GetYoozOmegaPlans(ctx, voucher, msisdn)` | خطط وباقات يوز أوميغا المخصصة |
| **100** | `GET` | `/api/v2/avocado/bundles` | `client.GetYoozBundles(ctx)` | استعراض جميع حزم وباقات يوز المتاحة للشراء |
| **101** | `GET` | `/api/v2/avocado/data-cap` | `client.GetYoozDataCap(ctx)` | استعلام سقف استهلاك البيانات المفروض على شريحة يوز |
| **102** | `POST` | `/api/v1/avocado/data-cap` | `client.SetYoozDataCap(ctx, limitMB)` | تعديل وتعيين سقف استهلاك البيانات لشريحة يوز |
| **103** | `GET` | `/api/v1/avocado/reward` | `client.GetYoozReward(ctx)` | استعلام مكافآت وهدايا ونقاط خطوط يوز |
| **104** | `GET` | `/api/v1/recharge/screen1` | `client.GetRechargeNumbers(ctx, option)` | المرحلة 1 من الشحن الإلكتروني: اختيار الرقم المراد تعبئته |
| **105** | `GET` | `/api/v1/recharge/screen2` | `client.GetRechargeTypes(ctx, option, msisdn)` | المرحلة 2: اختيار نوع الشحن والتعبئة |
| **106** | `GET` | `/api/v1/recharge/screen3` | `client.GetRechargeMethods(ctx, option, msisdn, method)` | المرحلة 3: بوابات ومحافظ الدفع الإلكتروني المتاحة (زين كاش، فاست باي، FIB) |
| **107** | `GET` | `/api/v1/recharge/screen4` | `client.GetOnlinePaymentDetails(ctx, option, msisdn, method, pgName)` | المرحلة 4: مبالغ الشحن المتاحة وتجهيز مسار بوابة الدفع |
| **108** | `GET` | `/api/v1/recharge/confirmation` | `client.GetRechargeConfirmation(ctx, txId)` | استعلام تأكيد وإيصال نجاح عملية الدفع والشحن الإلكتروني |
| **109** | `GET` | `/api/v1/top-up/payment-selection` | `client.GetPaymentSelection(ctx)` | خيارات وطرق الدفع لشحن الرصيد لأرقام أخرى |
| **110** | `GET` | `/api/v3/home` | `client.GetHomeDashboardV3(ctx, lat, lng, roaming)` | الصفحة الرئيسية والداشبورد المحدث v3 مع الأزرار السريعة والإحداثيات |
| **111** | `GET` | `/api/v3/addon/cyo` | `client.GetCYOBundles(ctx, groupId)` | خدمة صمم باقتك بنفسك (Create Your Own) باختيار سعة الإنترنت والدقائق |
| **112** | `GET` | `/api/v1/offer/active-offers` | `client.GetActiveOffers(ctx)` | استعلام قائمة العروض المباشرة الفعالة على الخط |
| **113** | `POST` | `/api/v3/profile/manage-quick-actions` | `client.ManageQuickActions(ctx, ids)` | حفظ وترتيب أزرار الوصول السريع في واجهة التطبيق |
| **114** | `POST` | `/api/v1/app-feedback` | `client.SubmitAppFeedback(ctx, cat, comment, rating)` | إرسال تقييم وملاحظات واقتراحات المشترك حول خدمات آسياسيل |
| **115** | `GET` | `/api/v1/survey` | `client.GetSurveys(ctx)` | استعلام استبيانات رضا المستخدمين وأسئلتها |
| **116** | `POST` | `/api/v1/survey` | `client.SubmitSurvey(ctx, surveyId, answers)` | إرسال إجابات استبيان رضا المشتركين |
| **117** | `GET` | `/api/v1/voc` | `client.GetVoiceOfCustomer(ctx)` | بوابة صوت العميل (Voice of Customer) واستعلام هوية المشترك |
| **118** | `GET` | `/api/v1/promotions/video-tutorials` | `client.GetVideoTutorials(ctx)` | جلب الفيديوهات والشروحات التعليمية الرسمية لكيفية استخدام الخدمات |
| **119** | `GET` | `/api/v1/one-yad` | `client.GetOneYadHome(ctx)` | لوحة تحكم مبادرة "يد واحدة" الاجتماعية والتبرعات |
| **120** | `GET` | `/api/v1/one-yad/teams` | `client.GetOneYadTeams(ctx)` | قائمة الفرق والمجموعات التطوعية التابعة للمبادرة |
| **121** | `POST` | `/api/v1/one-yad/request` | `client.SubmitOneYadRequest(ctx, teamId, amount)` | تقديم طلب تبرع أو مساهمة في مبادرة يد واحدة |
| **122** | `GET` | `/api/v1/epic` | `client.GetEpicLines(ctx)` | لوحة تحكم خطوط الشركات والخطوط المؤسسية (Epic Corporate) |
| **123** | `GET` | `/api/v1/epic/remaining/{msisdn}` | `client.GetEpicLineUsage(ctx, msisdn)` | استعلام الرصيد والمتبقي من البيانات والمكالمات للخط المؤسسي |
| **124** | `POST` | `/api/v1/profile/upload` | `client.UploadProfileImage(ctx, filename, r)` | رفع وتحديث صورة الحساب الرمزية (Avatar) كملف مالتيبارت |
| **125** | `POST` | `/api/v2/logo/upload` | `client.UploadPartnerLogo(ctx, filename, r)` | رفع شعار المتجر للشركاء التجاريين كملف مالتيبارت |
| **126** | `POST` | `/api/v1/notifications/register` | `client.RegisterNotificationToken(ctx, token, os)` | تسجيل توكن جهاز المشترك (FCM) لاستقبال الإشعارات السحابية |
| **127** | `POST` | `/api/v1/notifications/log` | `client.LogNotificationRead(ctx, notifId)` | تسجيل قراءة وتفاعل المشترك مع إشعار محدد |
| **128** | `POST` | `/api/v1/biometrics/register` | `client.RegisterBiometrics(ctx)` | تسجيل وتفعيل المصادقة البيومترية للجهاز واستخراج السر المشفر (Secret) لتأمين الدخول الدائم |
| **129** | `POST` | `/api/v1/biometrics/do-login` | `client.LoginBiometric(ctx, secret)` / `client.BiometricLogin(ctx, key)` | تسجيل الدخول الصامت والسريع عبر المفتاح والسر البيومتري بدون الحاجة لطلب OTP |
| **130** | `GET` | `/api/v1/watch/home` | `client.GetWatchDashboard(ctx)` | لوحة تحكم الساعات الذكية (Apple Watch و Wear OS) |
| **131** | `GET` | `/protected/v1/payments/{txId}/status` | `client.GetProtectedPaymentStatus(ctx, txId)` | فحص حالة عملية الدفع المحمية المشفرة |
| **132** | `POST` | `/protected/v1/payments/{txId}/cancel` | `client.CancelProtectedPayment(ctx, txId)` | إلغاء عملية الدفع المحمية المعلقة |
| **133** | `POST` | `/api/v1/top-up/omega` | `client.TopUpOmega(ctx, phone, voucher)` | شحن وتعبئة كروت وباقات خطوط أوميغا |
| **134** | `POST` | `/api/v1/delete` | `client.DeleteAccount(ctx)` | تقديم طلب رسمي لإغلاق وحذف الحساب نهائياً |
| **135** | `GET` | `/api/v1/asiaverse` | `client.GetAsiaverseHome(ctx)` | لوحة تحكم واستعلام منصة Asiaverse التفاعلية |
| **136** | `GET` | `/api/v1/shazam` | `client.GetShazamScanStatus(ctx)` | فحص حالة ومحاولات مسابقة مسح الباركود والفوز (Scan-to-Win) |
| **137** | `POST` | `/api/v1/shazam` | `client.SubmitShazamScan(ctx, qrCode)` | إرسال كود الباركود الممسوح للمسابقة واستلام الجائزة |
| **138** | `GET` | `/api/v1/profile/interests` | `client.GetUserInterests(ctx)` | استعلام قائمة اهتمامات وتفضيلات المشترك المحفوظة |
| **139** | `POST` | `/api/v1/avocado/profile/save-interests` | `client.SaveUserInterests(ctx, interests)` | حفظ وتحديث قائمة اهتمامات وتفضيلات المشترك |
| **140** | `GET` | `/api/v1/cdr/summary` | `client.GetCDRSummary(ctx)` | ملخص وإحصائيات تحويلات الرصيد الواردة والصادرة (CDR Summary) |
| **141** | `POST` | `/api/v1/map-account/resend` | `client.ResendLinkedAccountSMS(ctx, phone)` | إعادة إرسال رمز SMS لتأكيد ربط الخط الإضافي بالحساب |
| **142** | `POST` | `/api/v1/avocado/bundles/migrate-line` | `client.MigrateLineToYooz(ctx, req)` | تحويل الخط العادي إلى باقة وخط يوز (Yooz) بإرسال بيانات المستخدم وتاريخ الميلاد |
| **143** | `GET` | `/api/v1/avocado/migrate-out` | `client.GetYoozMigrateOutHome(ctx)` | استعلام شاشة وتعليمات طلب الخروج والرجوع من خط يوز إلى الخط العادي |
| **144** | `GET` | `/api/v1/avocado/migrate-out/select` | `client.GetYoozMigrateOutLocations(ctx)` | جلب قائمة مراكز وفروع الخدمة المعتمدة لإتمام تحويل الخط |
| **145** | `POST` | `/api/v1/avocado/migrate-out` | `client.SubmitYoozMigrateOut(ctx)` | تأكيد وإرسال طلب الخروج والتحويل النهائي من يوز للخط العادي |
| **146** | `GET` | `/api/v1/mosaic/migrate-out` | `client.GetMosaicMigrateOutHome(ctx)` | استعلام شاشة وتعليمات الرجوع من خطوط موزايك (Mosaic) إلى الخط العادي |
| **147** | `GET` | `/api/v1/mosaic/migrate-out/select` | `client.GetMosaicMigrateOutLocations(ctx)` | استعلام فروع ومراكز الخدمة المعتمدة لتحويل خطوط موزايك |

<br />

### دوال ومحركات إدارة الجلسة في Go SDK (Session Persistence & Keep-Alive Engine)

| # | دالة / واجهة Go SDK | نمط الاستخدام | الوصف الفني والمعماري |
| :---: | :--- | :--- | :--- |
| **S1** | `client.RefreshSession(ctx)` | تجديد ذكي ثلاثي الطبقات | تجديد متسلسل يبدأ بالتوكن (`/api/v1/validate`) ثم ينتقل تلقائياً لتسجيل الدخول البيومتري الصامت (`/api/v1/biometrics/do-login`) عند انتهاء الجلسة |
| **S2** | `client.StartKeepAlive(ctx, interval)` | نبض دوري خلفي (Heartbeat) | تشغيل مؤقت خلفي (Goroutine Ticker) كل 15 دقيقة لإبقاء الجلسة نشطة ومراقبة كشف الحساب وإطلاق `OnCDRExpired` |
| **S3** | `asiacell.NewFileSessionStorage(path)` | تخزين ذري دائم للجلسة | حفظ ذري (Atomic Write عبر ملف مؤقت) لبيانات الجلسة (Tokens, DeviceID, Phone, BiometricSecret) لاسترجاعها عند الإقلاع |
| **S4** | `asiacell.NewMemorySessionStorage()` | تخزين مؤقت بالذاكرة | تخزين آمن لبيانات الجلسة في الذاكرة (Thread-safe) للأنظمة والحاويات عديمة الحالة (Stateless) |
| **S5** | `client.ExportSession()` / `ImportSession()` | تصدير واستيراد الجلسة | تصدير واستيراد كائن الجلسة المشفر كـ JSON وتمريره بين السيرفرات أو البوتات المختلفة بسلاسة |

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

#### تفعيل كشف الحساب والتحقق الثنائي (CDR 2FA Activation):
وفقاً للهندسة العكسية الدقيقة لتطبيق آسياسيل الرسمي (v5.1.0)، تتطلب عملية تفعيل كشف الحساب خطوتين متتاليتين بحزم:

1. **إرسال رمز التحقق واستخراج معرّف العملية (Send OTP & Extract PID)**:
   <div dir="ltr" align="left">

   ```http
   POST /api/v1/cdr/send-otp?lang=ar
   ```

   ```json
   // استجابة الخادم الرسمية:
   {
     "code": 200,
     "message": "success",
     "success": true,
     "nextUrl": "/api/v1/cdr/confirm?PID=82d7c541-16ef-46be-91c6-8fb50b5557ef"
   }
   ```

   </div>

   تقوم دالة `client.SendCDROTP(ctx)` بإرسال الطلب، وتحليل حقل `nextUrl` لاستخراج معرّف العملية `PID` وتخزينه تلقائياً في كائن العميل لتأمين إتمام العملية.

2. **تأكيد الرمز عبر بنية GenericSMSConfirmationDTO الرسمية**:
   <div dir="ltr" align="left">

   ```http
   POST /api/v1/cdr/confirm?lang=ar
   ```

   ```json
   {
     "PID": "82d7c541-16ef-46be-91c6-8fb50b5557ef",
     "passcode": "123456"
   }
   ```

   </div>

   تقوم دالة `client.ConfirmCDROTP(ctx, pid, passcode)` (أو `client.ConfirmCDROTP(ctx, passcode)` التي تستخدم الـ PID المحفوظ مسبقاً) بتمرير بنية `GenericSMSConfirmationDTO` المعتمدة رسمياً من البوابة السحابية.

   > [!CAUTION]
   > **قاعدة هندسية حرجة**: يُحظر تماماً استدعاء مسار `/api/v1/smsvalidation` لتأكيد كود كشف الحساب (CDR). مسار `smsvalidation` مخصص فقط لعملية تسجيل الدخول الأولى، واستدعاؤه مع كود كشف الحساب يؤدي فوراً إلى احتراق الرمز (OTP Burn) وإبطال صلاحية جلسة المستخدم بالكامل.

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

### 2.2 الاشتراكات المفعلة وخدمات الـ USSD السحابية (Subscriptions & USSD API)

- `GET /api/v1/profile/subscriptions`: استعلام الخادم الرسمي عن جميع الخدمات النشطة على الخط مع تواريخ الصلاحية وأزرار الإلغاء (`MySubscriptionsResponse`).
- `GET /api/v1/ussd?parent_id={id}`: استعلام قائمة خيارات أوامر الـ USSD السحابية التفاعلية من الخادم.
- `POST /api/v1/ussd`: إرسال وتأكيد خيار أو أمر USSD عبر السحابة مباشرة دون الحاجة للاتصال من الهاتف.

<div dir="ltr" align="left">

```http
POST /api/v1/ussd
```

```json
{
  "parent_id": "0",
  "choice": "1"
}
```

</div>

---

### 2.3 سوق الأرقام المميزة (Vanity VIP Numbers)

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

### 2.4 إهداء الباقات للغير (Send Addon as a Gift)

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

### 2.5 نظام الشكاوى والتذاكر الفنية (Resolution Center)

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

### 2.6 نظام التعويضات التلقائي (Compensation System)

<div dir="ltr" align="left">

```http
GET /api/v1/compensation
```

</div>

- **الغرض**: فحص استحقاق الخط للتعويضات الرسمية المعتمدة من آسياسيل عند وجود أعطال شبكة عامة، واستلام باقات إنترنت أو رصيد مجاني.

---

### 2.7 خطوط الشباب Yooz (MGM Referral Program)

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

### 2.8 كروت الألعاب والشحن الرقمي (E-Vouchers)

<div dir="ltr" align="left">

```http
GET /api/v2/e-voucher/packages?recharge-type=1
```

</div>

- **الغرض**: استعراض كروت الألعاب والتطبيقات المتاحة للشراء برصيد الهاتف (PUBG, PlayStation, iTunes, etc.).

---

### 2.9 معمارية تسجيل الدخول والجلسة الخالدة (Authentication & Immortal Session Architecture)

تعتمد حزمة Go SDK نظام مصادقة ثلاثي الطبقات متقدم (3-Layer Immortal Authentication) مستخرج ومطابق بنسبة 100% لسلوك تطبيق آسياسيل الرسمي (v5.1.0)، مما يضمن استمرارية عمل البوتات والخدمات السحابية لأشهر متواصلة دون انقطاع ودون الحاجة لإعادة طلب رمز SMS من المشترك.

#### 1. مسارات دورة حياة المصادقة الرسمية:
<div dir="ltr" align="left">

```http
POST /api/v1/login
POST /api/v1/smsvalidation
POST /api/v1/validate
POST /api/v1/biometrics/register
POST /api/v1/biometrics/do-login
```

</div>

---

#### 2. تدفق تسجيل الدخول الأولي وتأسيس البصمة (Initial Login & Bootstrapping):
1. **طلب رمز SMS**: استدعاء `client.Login(ctx, phone)` الذي يخاطب `POST /api/v1/login` ويستخرج معرّف `PID` من حقل `nextUrl`.
2. **التحقق من الكود وتوليد الجلسة**: استدعاء `client.VerifySMS(ctx, pid, code)` الذي يخاطب `POST /api/v1/smsvalidation` لإرجاع توكنات الجلسة (`access_token`, `refresh_token`, `handshake_token`, `secret`).
3. **تفعيل البصمة التلقائي**: تقوم دالة `VerifySMS` فور نجاحها باستدعاء `POST /api/v1/biometrics/register` لتسجيل الجهاز على خوادم آسياسيل واستخراج المفتاح السري المشفر `BiometricSecret` وتخزينه في كائن الجلسة.

---

#### 3. معمارية الجلسة الخالدة ثلاثية الطبقات (3-Layer Fallback Architecture):
تعمل دالة الطلبات المركزية `doRequest` في الـ SDK بمحرك إعادة محاولة ذاتي وشفاف عند استقبال أي خطأ مصادقة (`401 Unauthorized` أو `403 Forbidden` أو `493 Session Expired`):

```
+-------------------------------------------------------------+
|               الطبقة 1: توكن الوصول الفعال                  |
|       Authorization: Bearer <access_token>                  |
+-------------------------------------------------------------+
                              |
                     فشل (401 / 403 / 493)
                              v
+-------------------------------------------------------------+
|               الطبقة 2: تجديد التوكن الصامت                 |
|               POST /api/v1/validate                         |
|        Payload: {"refreshToken": "Bearer <refresh_token>"}  |
|   يعيد توكنات وصول وتجديد جديدة + مفتاح Secret محدث         |
+-------------------------------------------------------------+
                              |
                    فشل (انتهاء صلاحية التجديد)
                              v
+-------------------------------------------------------------+
|        الطبقة 3: تسجيل الدخول البيومتري الصامت التام        |
|             POST /api/v1/biometrics/do-login                |
|      Payload: {"msisdn": phone, "secret": biometricSecret}  |
|  يعيد توليد جلسة كاملة فوراً بدون SMS وبدون تدخل بشري       |
+-------------------------------------------------------------+
```

---

#### 4. معرّف الجهاز الثابت الحصين (Persistent Immutable DeviceID):
- يتم توليد معرّف جهاز عشوائي بصيغة **UUID v4** مرة واحدة فقط ويُحفظ دائماً مع بيانات الجلسة (`SessionData.DeviceID`).
- يُحظر تماماً تدوير أو إعادة توليد المعرّف عشوائياً عند كل طلب أو إعادة تشغيل، لتجنب قيام خوادم آسياسيل بحظر الخط أو إبطال الجلسة.
- يُرسل المعرّف إلزامياً في جميع الطلبات عبر ترويستين متطابقتين:
  - `DeviceId: <UUID-v4>`
  - `x-device-id: <UUID-v4>`

---

#### 5. محرك النبض الدوري الخلفي (Keep-Alive Heartbeat Daemon):
تشغيل حارس الجلسة في الخلفية عبر استدعاء:
```go
client.StartKeepAlive(ctx, 15*time.Minute)
```
- يرسل نبضات دورية منتظمة (Heartbeat Pulses) كل 15 دقيقة لفحص حالة الملف الشخصي (`GET /api/v1/profile`) وسجل كشف الحساب (`GET /api/v1/cdr/detail?type=btransfer&page=1&limit=1`).
- يضمن منع خمول الجلسة على السيرفر (Keep Session Warm).
- يراقب صلاحية كشف الحساب، وعند انتهائها يقوم بإطلاق هوك `OnCDRExpired()` فوراً لإشعار النظام أو إرسال تنبيه لإعادة تفعيل التحقق الثنائي (2FA).

---

#### 6. التخزين الذري الدائم للجلسة (Atomic Session Persistence):
توفر المكتبة واجهة `SessionStorage` لحفظ واسترجاع حالة الجلسة تلقائياً عند أي عملية تجديد أو تغيير:
- `FileSessionStorage`: حفظ مشفر وذري (Atomic Write عبر كتابة ملف مؤقت ثم استبداله `os.Rename`) لضمان عدم تلف ملف الجلسة عند انقطاع التيار أو إيقاف السيرفر المفاجئ.
- `MemorySessionStorage`: تخزين خفيف وآمن تزامنياً (Thread-safe) للأنظمة والحاويات عديمة الحالة (Stateless).

---

### 2.10 تحويل الخط إلى يوز والرجوع إلى الخط العادي (Line Migration Flow)

<div dir="ltr" align="left">

```http
POST /api/v1/avocado/bundles/migrate-line
GET  /api/v1/avocado/migrate-out
GET  /api/v1/avocado/migrate-out/select
POST /api/v1/avocado/migrate-out
GET  /api/v1/mosaic/migrate-out
GET  /api/v1/mosaic/migrate-out/select
```

```json
// طلب التحويل إلى خط Yooz
{
  "dob": "1998-05-20",
  "name": "Ali Ahmed",
  "avatar": "avatar-5"
}
```

</div>

- **الغرض**:
  1. **التحويل إلى Yooz**: ترقية ونقل الخط العادي إلى نظام وباقات Yooz الرقمية مع تسجيل الاسم وتاريخ الميلاد والأفاتار.
  2. **الرجوع إلى الخط العادي (Migrate Out)**:
     - `GET /api/v1/avocado/migrate-out`: استعلام شاشة التأكيد والتعليمات والشروط الخاصة بإلغاء باقة Yooz والرجوع لنظام الخطوط العادي.
     - `GET /api/v1/avocado/migrate-out/select`: جلب فروع ومراكز الخدمة المعتمدة لإتمام التحويل وتثبيت الهوية.
     - `POST /api/v1/avocado/migrate-out`: إرسال وتثبيت طلب الخروج والرجوع المباشر للخط العادي.
  3. **خطوط موزايك (Mosaic)**: استعلام شاشات ومراكز الرجوع للخط العادي عبر `/api/v1/mosaic/migrate-out`.

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

> [!IMPORTANT]
> **تطابق ترويسات معرّف الجهاز (DeviceID Dual-Header Invariant)**:
> يجب تمرير معرّف الجهاز الثابت (UUID v4) في الترويستين معاً: `DeviceId` بالحرف الكبير و `x-device-id` بالأحرف الصغيرة مع تطابق تام في القيمة وعدم تغييرها طوال دورة حياة الحساب لضمان عدم سقوط الجلسة.
