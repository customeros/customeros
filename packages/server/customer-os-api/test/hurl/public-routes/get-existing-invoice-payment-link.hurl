# Test: Get existing invoice payment link
GET {{cos_url}}/invoice/a10f97ff-a14a-4675-807d-3079f5cadff2/paymentLink
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
address: silviu@customeros.ai

HTTP 200
[Asserts]
header "Content-Type" == "text/plain; charset=utf-8"
body matches "https://checkout.stripe.com/c/pay/cs_test_.*"
