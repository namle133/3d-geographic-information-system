# VNPay payment deployment

The payment integration uses VNPay PAY 2.1.0. The Merchant View migration notice concerns the provider's payment-link management service; it does not require replacing the PAY endpoint.

## Railway variables

Set `VNPAY_TMN_CODE` and `VNPAY_HASH_SECRET` from your sandbox registration. Do not put credentials in source or logs.

```dotenv
VNPAY_PAY_URL=https://sandbox.vnpayment.vn/paymentv2/vpcpay.html
VNPAY_RETURN_URL=https://3d-geographic-information-system.vercel.app/pages/payment-result.html
VNPAY_IPN_URL=https://3d-geographic-information-system-production.up.railway.app/api/v1/payments/vnpay/ipn
```

Register the public IPN URL with VNPay. Setting the variable alone does not register it. Error 71, "Website này chưa được phê duyệt", requires VNPay to check the terminal approval/activation status.

The current application confirms payments only through IPN. The return page verifies the signed return parameters, displays the database state, polls briefly while waiting for IPN, and offers a status retry without starting another payment.

## Merchant View account

The official link shown by the old sandbox Merchant Admin points to https://doitac-test.vnpaytest.vn/. It uses VNPAY ID. In the account check on 19 November 2026, VNPAY ID reported that the supplied merchant email did not exist. Ask VNPay to link/provision the appropriate sandbox VNPAY ID and associate the existing terminal; existing merchant credentials alone do not establish successful migration.

Suggested support request (send through your established VNPay support contact):

> Terminal sandbox KD4N9ECY returns error 71, "Website này chưa được phê duyệt". Please check terminal approval/activation, register the IPN URL above, and advise how to link the existing merchant account to VNPAY ID/Merchant View. Merchant View currently reports that the registered merchant email does not exist. The backend health endpoint and IPN endpoint are publicly reachable.

## Data changes and reconciliation

Startup adds nullable `invoice_payments.payment_transaction_id` with a unique index, and `payment_transactions.expires_at` / `requires_review`. Existing rows retain their values. Locks and the unique index prevent duplicate credit for a single transaction. Active URLs are reused for the initiating user; another payment attempt is blocked until the active order completes or expires. Late captured funds are kept in the ledger and flagged for reconciliation if the invoice was cancelled or no longer owes that amount. Historical duplicate payments are not automatically rewritten.

## Password recovery

Public registration accepts residents only. Privileged users are created through the authenticated admin API. Password reset codes are hashed in the database, expire after 30 minutes, are consumed once, and are sent only to the account's registered email. API responses never return reset codes. Configure `SMTP_HOST`, `SMTP_PORT` (587 STARTTLS or 465 TLS), `SMTP_USER`, `SMTP_PASSWORD`, and optionally `SMTP_FROM` (defaults to SMTP_USER). Without SMTP configuration, recovery returns a clear unavailable message; it does not fall back to exposing tokens.

## Verification

```sh
cd backend
./scripts/test-integration.sh  # temporary PostgreSQL container; runs Go tests with race detector
```

```sh
node --test tests/payment-result.test.mjs
```

Integration tests cover concurrent IPNs, concurrent order creation, status/amount validation, return URLs without ledger writes, rollback on database failures, overpayment review, public registration restrictions, and private single-use reset codes. No provider transaction is created by these tests.
