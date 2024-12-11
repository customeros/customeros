# Test: Get Invoices For Inexisting Organization
GET {{cos_url}}/billing/v1/organizations/f75add8c-4f00-442d-8bc5-e739a655bafq/invoices
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 404
[Asserts]
jsonpath "$.status" == "error"
jsonpath "$.requestId" exists
jsonpath "$.message" == "Organization does not exist"

