# Test: Get Mailboxes
GET {{cos_url}}/mailstack/v1/domains/q1w2e3r4/mailboxes
Content-Type: application/json
X-CUSTOMER-OS-API-KEY: {{api_key}}

HTTP/2 200
[Asserts]
body == ""
