# Test: Redirect to pay non-existing invoice
GET {{cos_url}}/invoice/a10f97ff-a14a-4675-807d-3079f5cadff1/pay
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}
[QueryStringParams]
address: silviu@customeros.ai

HTTP 404
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
jsonpath "$.error" == "Invoice not found"
