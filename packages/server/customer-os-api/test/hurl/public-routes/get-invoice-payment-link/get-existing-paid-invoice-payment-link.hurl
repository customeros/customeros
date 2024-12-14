# Test: Get existing paid invoice payment link
GET {{cos_url}}/invoice/b15a1559-2ffc-4e35-b6bc-301e6e7fbf80/paymentLink
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
address: silviu@customeros.ai

HTTP 409
[Asserts]
header "Content-Type" == "text/plain; charset=utf-8"
body == ""
