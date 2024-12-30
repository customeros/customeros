# Test: Get Domains [Bug COS-6334]
GET {{cos_url}}/mailstack/v1/domains/
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 *
[Asserts]
status == 200
body == ""
