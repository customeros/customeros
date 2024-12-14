# Test: Redirect to pay existing paid invoice
GET {{cos_url}}/invoice/b15a1559-2ffc-4e35-b6bc-301e6e7fbf80/pay
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
address: silviu@customeros.ai

HTTP 303
[Asserts]
header "Content-Type" == "text/html; charset=utf-8"
body == "<a href=\"https://customeros.ai/payments/status/paid/\">See Other</a>.\n\n"
