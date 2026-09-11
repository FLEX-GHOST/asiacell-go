# Asiacell Go API - Examples

This directory contains clean, production-grade, standalone examples demonstrating every capability of the `pkg/asiacell` library.

---

## Examples Overview

| Directory | Description | Run Command |
| :--- | :--- | :--- |
| **`01_auth_and_session`** | Phone login via SMS OTP, Captcha auto-bypass/solving, session import & export. | `go run examples/01_auth_and_session/main.go` |
| **`02_account_and_profile`** | Real-time balance query, remaining calls, SMS, active internet packages, and profile info. | `go run examples/02_account_and_profile/main.go` |
| **`03_bundles_and_4g`** | Browsing 4G unlimited bundles, special offers, addon categories, and API subscription. | `go run examples/03_bundles_and_4g/main.go` |
| **`04_credit_transfer`** | P2P credit transfer, SMS OTP confirmation, wallet transfers, and transfer verification. | `go run examples/04_credit_transfer/main.go` |
| **`05_recharge_voucher`** | Recharging phone balance via 14-digit voucher cards and checking recharge history. | `go run examples/05_recharge_voucher/main.go` |
| **`06_services_management`**| Internet balance protection (*223#), cancellation guides (299/4151/300), digital services. | `go run examples/06_services_management/main.go` |
| **`07_shops_and_governorates`**| Iraqi governorates/cities coverage, certified branches, coordinates, and real-time open status. | `go run examples/07_shops_and_governorates/main.go` |
| **`08_spin_wheel_and_rewards`**| Checking and playing the daily spin wheel, Shukran emergency loans, roaming, promotions. | `go run examples/08_spin_wheel_and_rewards/main.go` |
| **`09_cdr_incoming_transfer_verification`**| Automated incoming credit transfer verification via CDR ledger without admin intervention. | `go run examples/09_cdr_incoming_transfer_verification/main.go` |
| **`interactive_cli`** | An all-in-one terminal CLI application with an interactive text menu for all features. | `go run examples/interactive_cli/main.go` |

---

## Quick Start

### 1. Interactive CLI (All-in-One)
To test and interact with all features using a terminal menu:
```bash
go run examples/interactive_cli/main.go
```

### 2. Login & Session Setup
To log in with your phone and save a persistent `session.json` file:
```bash
go run examples/01_auth_and_session/main.go
```

### 3. Check Account Balance & Bundles
```bash
go run examples/02_account_and_profile/main.go
```

### 4. Automated CDR Transfer Verification (Payment Gateway)
To verify incoming customer payments against the official Asiacell CDR ledger:
```bash
go run examples/09_cdr_incoming_transfer_verification/main.go
```
