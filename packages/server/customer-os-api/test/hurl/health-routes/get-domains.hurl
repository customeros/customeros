# Test: Get Domains
GET {{cos_url}}/mailstack/v1/domains/
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 200
[Asserts]
body == ""
